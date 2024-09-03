package httplib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultResponse struct {
	Status    string      `json:"status"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	DataError interface{} `json:"dataError"`
}

type DefaultPaginationResponse struct {
	Status     string      `json:"status"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	TotalCount uint64      `json:"totalCount"`
	TotalPages uint64      `json:"totalPages"`
	Data       interface{} `json:"data"`
}

func SetSuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, DefaultResponse{
		Status:  http.StatusText(code),
		Code:    code,
		Data:    data,
		Message: message,
	})
	return
}

func SetPaginationResponse(c *gin.Context, code int, message string, data interface{}, totalCount uint64, pg *Query) {
	c.JSON(code, DefaultPaginationResponse{
		Status:     http.StatusText(code),
		Code:       code,
		Message:    message,
		Page:       pg.GetPage(),
		Size:       pg.GetSize(),
		TotalCount: totalCount,
		TotalPages: uint64(GetTotalPages(int(totalCount), pg.GetSize())),
		Data:       data,
	})
	return
}

func SetErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, DefaultResponse{
		Status:  http.StatusText(code),
		Code:    code,
		Data:    nil,
		Message: message,
	})
	return
}

func SetCustomResponse(c *gin.Context, code int, message string, data interface{}, dataErr interface{}) {
	c.JSON(code, DefaultResponse{
		Status:    http.StatusText(code),
		Code:      code,
		Data:      data,
		Message:   message,
		DataError: dataErr,
	})
	return
}
