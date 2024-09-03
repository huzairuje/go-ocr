package health

import (
	"errors"
	"net/http"

	"go-ocr/infrastructure/httplib"
	logger "go-ocr/infrastructure/log"
	"go-ocr/modules/primitive"
	"go-ocr/utils"

	"github.com/gin-gonic/gin"
)

type Http struct {
	serviceHealth InterfaceService
}

func NewHttp(serviceHealth InterfaceService) InterfaceHttp {
	return &Http{
		serviceHealth: serviceHealth,
	}
}

type InterfaceHttp interface {
	GroupHealth(group *gin.RouterGroup)
}

func (h *Http) GroupHealth(g *gin.RouterGroup) {
	g.GET("/ping", h.Ping)
	g.GET("/check", h.HealthCheckApi)
}

func (h *Http) Ping(c *gin.Context) {
	httplib.SetSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), "pong")
}

func (h *Http) HealthCheckApi(c *gin.Context) {
	logCtx := "handler.HealthCheckApi"

	if h.serviceHealth == nil {
		err := errors.New("dependency service health to handler health is nil")
		logger.Error(c, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceHealth")
		httplib.SetErrorResponse(c, http.StatusInternalServerError, primitive.SomethingWentWrong)
		return
	}

	resp, err := h.serviceHealth.CheckUpTime(c)
	if err != nil {
		logger.Error(c, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceHealth.CheckUpTime")
		httplib.SetErrorResponse(c, http.StatusInternalServerError, primitive.SomethingWentWrong)
		return
	}
	httplib.SetSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), resp)
	return
}
