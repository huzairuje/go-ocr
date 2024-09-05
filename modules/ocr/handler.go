package ocr

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go-ocr/infrastructure/httplib"
	logger "go-ocr/infrastructure/log"
	"go-ocr/infrastructure/validator"
	"go-ocr/modules/primitive"
	"go-ocr/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Http struct {
	serviceOcr ServiceInterface
}

func NewHttp(serviceOcr ServiceInterface) InterfaceHttp {
	return &Http{
		serviceOcr: serviceOcr,
	}
}

type InterfaceHttp interface {
	GroupOcr(group *gin.RouterGroup)
}

func (h *Http) GroupOcr(g *gin.RouterGroup) {
	g.POST("", h.ProcessOCR)
	g.GET("", h.GetListOcr)
	g.GET("/:id", h.DetailOCR)
}

func (h *Http) ProcessOCR(ctx *gin.Context) {
	logCtx := fmt.Sprintf("handler.ProcessOCR")

	var requestBody primitive.OcrRequest
	if err := ctx.ShouldBind(&requestBody); err != nil {
		logger.Error(ctx, logCtx, "ctx.ShouldBind got err : %v", err)
		httplib.SetErrorResponse(ctx, http.StatusBadRequest, primitive.SomethingWrongWithTheBodyRequest)
		return
	}

	errValidateStruct := validator.ValidateStructResponseSliceString(requestBody)
	if errValidateStruct != nil {
		logger.Error(ctx, logCtx, "validator.ValidateStructResponseSliceString got err : %v", errValidateStruct)
		httplib.SetCustomResponse(ctx, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil, errValidateStruct)
		return
	}

	// Get uploaded file
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	response, err := h.serviceOcr.ProcessOcr(ctx, requestBody, file, fileHeader)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceOcr.ProcessOcr")
		httplib.SetErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httplib.SetSuccessResponse(ctx, http.StatusOK, primitive.ProcessOcrSuccess, response)
	return
}

func (h *Http) GetListOcr(ctx *gin.Context) {
	logCtx := fmt.Sprintf("handler.GetListOcr")

	paginationQuery, err := httplib.GetPaginationFromCtx(ctx)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "httplib.GetPaginationFromCtx")
		httplib.SetErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	status := ctx.Request.URL.Query().Get("status")
	if status != "" {
		if !utils.IsValidSanitizeSQL(status) {
			err = errors.New(primitive.QueryIsSuspicious)
			logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "utils.IsValidSanitizeSQL")
			httplib.SetErrorResponse(ctx, http.StatusBadRequest, primitive.QueryIsSuspicious)
			return
		}
	}

	text := ctx.Request.URL.Query().Get("text")
	if text != "" {
		if !utils.IsValidSanitizeSQL(text) {
			err = errors.New(primitive.QueryIsSuspicious)
			logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "utils.IsValidSanitizeSQL")
			httplib.SetErrorResponse(ctx, http.StatusBadRequest, primitive.QueryIsSuspicious)
			return
		}
	}

	disablePaginationHeader := utils.IsDisablePagination(ctx)
	var param primitive.ParameterFindOcr
	if disablePaginationHeader {
		param = primitive.ParameterFindOcr{
			Text:   text,
			Status: status,
		}
	} else {
		param = primitive.ParameterFindOcr{
			Text:      text,
			Status:    status,
			PageSize:  paginationQuery.GetSize(),
			Offset:    paginationQuery.GetOffset(),
			SortBy:    paginationQuery.GetOrderBy(),
			SortOrder: paginationQuery.GetSortOrder(),
		}
	}

	data, _, err := h.serviceOcr.ListOcr(ctx, disablePaginationHeader, param)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceArticle.GetListArticle")
		httplib.SetErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if disablePaginationHeader {
		httplib.SetSuccessResponse(ctx, http.StatusOK, http.StatusText(http.StatusOK), data)
		return
	} else {
		httplib.SetPaginationResponse(ctx, http.StatusOK, http.StatusText(http.StatusOK), data, uint64(len(data)), paginationQuery)
		return

	}
}

func (h *Http) DetailOCR(ctx *gin.Context) {
	logCtx := fmt.Sprintf("handler.DetailOCR")

	idParam := ctx.Param("id")
	if idParam == "" {
		err := errors.New(primitive.ParamIdIsZeroOrNullString)
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "c.Param")
		httplib.SetErrorResponse(ctx, http.StatusBadRequest, primitive.ParamIdIsZeroOrNullString)
		return
	}

	idInt, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || idInt == 0 {
		err := errors.New(primitive.ParamIdIsZeroOrNullString)
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "strconv.Atoi")
		httplib.SetErrorResponse(ctx, http.StatusBadRequest, primitive.ParamIdIsZeroOrNullString)
		return
	}

	data, err := h.serviceOcr.GetRecordOcrById(ctx, idInt)
	if err != nil {
		errNotFound := []error{gorm.ErrRecordNotFound}
		if utils.ContainsError(err, errNotFound) {
			logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceServices.GetRecordServicesById")
			httplib.SetErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceServices.GetRecordServicesById")
		httplib.SetErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httplib.SetSuccessResponse(ctx, http.StatusOK, http.StatusText(http.StatusOK), data)
	return

}
