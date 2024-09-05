package ocr

import (
	"fmt"
	"net/http"
	
	"go-ocr/infrastructure/httplib"
	logger "go-ocr/infrastructure/log"
	"go-ocr/infrastructure/validator"
	"go-ocr/modules/primitive"
	"go-ocr/utils"

	"github.com/gin-gonic/gin"
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
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	response, err := h.serviceOcr.ProcessOcr(ctx, requestBody, file)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceOcr.ProcessOcr")
		httplib.SetErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	httplib.SetSuccessResponse(ctx, http.StatusOK, primitive.ProcessOcrSuccess, response)
	return
}
