package r

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wan_go/pkg/utils"
)

const (
	SuccessStatus = 0
	ErrorStatus   = -1
)

type Response struct {
	Code    int    `json:"code,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type response struct {
	Response
	Data any `json:"data"`
}

func Success(data any, c *gin.Context) {
	res := &response{}
	res.Message = http.StatusText(http.StatusOK)
	res.Status = SuccessStatus
	res.Code = http.StatusOK
	res.Data = data

	c.JSON(http.StatusOK, res)
}

func Error(httpStatus int, msg string, c *gin.Context) {
	res := &response{}
	res.Message = utils.IfThen(len(msg) > 0, msg, http.StatusText(httpStatus)).(string)
	res.Status = ErrorStatus
	res.Code = httpStatus

	c.JSON(httpStatus, res)
}

func ErrorInternal(msg string, c *gin.Context) {
	Error(http.StatusInternalServerError, msg, c)
}
