package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go-ocr/infrastructure/config"
	logger "go-ocr/infrastructure/log"
	redisLocal "go-ocr/infrastructure/redis"
	"go-ocr/modules/primitive"
	"go-ocr/utils"

	"github.com/otiai10/gosseract/v2"
)

const (
	redisFinaleKeyOcr     = "ocr:%d"
	redisListFinaleKeyOcr = "ocr_list"
)

type ServiceInterface interface {
	ProcessOcr(ctx context.Context, payload primitive.OcrRequest, file multipart.File, fileHeader *multipart.FileHeader) (primitive.OCrResponse, error)
	ListOcr(ctx context.Context, isDisablePagination bool, param primitive.ParameterFindOcr) (res []primitive.OCrResponse, count int64, err error)
	GetRecordOcrById(ctx context.Context, id int64) (primitive.OCrResponse, error)
}

type Service struct {
	repository       RepositoryInterface
	redisInterface   redisLocal.LibInterface
	tesseractsClient *gosseract.Client
}

func NewService(repository RepositoryInterface, redisInterface redisLocal.LibInterface, tesseractsClient *gosseract.Client) ServiceInterface {
	return &Service{
		repository:       repository,
		redisInterface:   redisInterface,
		tesseractsClient: tesseractsClient,
	}
}

func (s *Service) ProcessOcr(ctx context.Context, payload primitive.OcrRequest, file multipart.File, fileHeader *multipart.FileHeader) (primitive.OCrResponse, error) {
	logCtx := fmt.Sprintf("service.RecordOcr")

	// Define the file path where the image will be saved
	uploadDir := "./uploads"
	filePath := filepath.Join(uploadDir, fileHeader.Filename)
	payload.Image = filePath

	// Ensure the directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "os.MkdirAll")
		return primitive.OCrResponse{}, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Create a file at the specified location
	fileCreated, err := os.Create(filePath)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "os.Create")
		return primitive.OCrResponse{}, fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		fileCreated.Close()
	}()

	// Write the uploaded file content to the created file
	if _, err := io.Copy(fileCreated, file); err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "io.Copy")
		return primitive.OCrResponse{}, fmt.Errorf("failed to save uploaded file: %w", err)
	}

	err = s.tesseractsClient.SetImage(fileCreated.Name())
	if err != nil {
		return primitive.OCrResponse{}, err
	}

	var textResult string
	var isEnabledHOCR bool
	if payload.HOCREnabled != "" {
		isEnabledHOCR, err = strconv.ParseBool(payload.HOCREnabled)
		if err != nil {
			return primitive.OCrResponse{}, err
		}
	} else {
		isEnabledHOCR = false
	}

	if isEnabledHOCR {
		textResult, err = s.tesseractsClient.HOCRText()
		if err != nil {
			return primitive.OCrResponse{}, err
		}
	} else {
		textResult, err = s.tesseractsClient.Text()
		if err != nil {
			return primitive.OCrResponse{}, err
		}
	}
	textResult = strings.Trim(textResult, "\n")

	payloadDb := primitive.Ocr{
		ImageUrl: payload.Image,
		Text:     textResult,
		Status:   "SUCCESSFUL",
	}

	data, err := s.repository.CreateOcr(ctx, payloadDb)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "u.repository.CountArticle")
		return primitive.OCrResponse{}, err
	}

	//set data to redis on goroutine
	if config.Conf.Redis.EnableRedis && s.redisInterface != nil {
		go func() {
			dataBytes, errMarshall := json.Marshal(data)
			if errMarshall != nil {
				logger.Error(ctx, utils.ErrorLogFormat, errMarshall.Error(), logCtx, "json.Marshal")
			}
			redisFinaleKey := fmt.Sprintf(redisFinaleKeyOcr, data.ID)
			errSetToRedis := s.redisInterface.Set(redisFinaleKey, dataBytes, time.Minute)
			if errSetToRedis != nil {
				logger.Error(ctx, utils.ErrorLogFormat, errSetToRedis.Error(), logCtx, "s.redis.Set")
			}
			fmt.Printf("success SET on redis by key: %s\n", redisFinaleKey)
		}()
	}

	payloadResp := primitive.OCrResponse{
		ID:        data.ID,
		ImageUrl:  data.ImageUrl,
		Text:      data.Text,
		Status:    data.Status,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return payloadResp, nil

}

func (s *Service) ListOcr(ctx context.Context, isDisablePagination bool, param primitive.ParameterFindOcr) (res []primitive.OCrResponse, count int64, err error) {
	logCtx := fmt.Sprintf("service.ListPaymentAll")

	emptySliceDataOcr := make([]primitive.OCrResponse, 0)
	// Data not found in cache, query the database
	count, err = s.repository.CountAllListOcr(ctx, param)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "u.repository.CountAllList")
		return
	}

	var listData []primitive.Ocr
	if isDisablePagination {
		listData, err = s.repository.FindAllListOcrNonPagination(ctx, param)
		if err != nil {
			logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "u.repository.FindListArticle")
			return
		}
	} else {
		if param.SortBy == "" {
			param.SortBy = "id"
		}
		if param.SortOrder == "" {
			param.SortOrder = "desc"
		}
		listData, err = s.repository.FindAllListOcrPagination(ctx, param)
		if err != nil {
			logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "u.repository.FindListArticle")
			return
		}
	}

	if count == 0 || len(listData) == 0 {
		return emptySliceDataOcr, 0, nil
	}

	var list []primitive.OCrResponse
	if len(listData) > 0 {
		for _, val := range listData {

			list = append(list, primitive.OCrResponse{
				ID:        val.ID,
				ImageUrl:  val.ImageUrl,
				Text:      val.Text,
				Status:    val.Status,
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.CreatedAt,
			})
		}
		res = list
	}

	return res, count, nil
}

func (s *Service) GetRecordOcrById(ctx context.Context, id int64) (primitive.OCrResponse, error) {
	logCtx := fmt.Sprintf("service.GetRecordPaymentById")

	data, err := s.repository.FindOcrByID(ctx, id)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "u.repository.FindPaymentByIdAndByCustomerId")
		return primitive.OCrResponse{}, err
	}

	payloadResp := primitive.OCrResponse{
		ID:        data.ID,
		ImageUrl:  data.ImageUrl,
		Text:      data.Text,
		Status:    data.Status,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.CreatedAt,
	}

	return payloadResp, nil

}
