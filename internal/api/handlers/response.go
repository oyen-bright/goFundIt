package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/pkg/errs"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func DefaultResponse(c *gin.Context, statusCode int, message string, data interface{}, err interface{}) {
	c.JSON(statusCode, response{
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
		Errors:  err,
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	DefaultResponse(c, http.StatusOK, message, data, nil)
}

func Created(c *gin.Context, message string, data interface{}) {
	DefaultResponse(c, http.StatusCreated, message, data, nil)
}

func BadRequest(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, response{
		Status:  http.StatusText(http.StatusBadRequest),
		Message: message,
		Errors:  errors,
	})
}

func Unauthorized(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnauthorized, response{
		Status:  http.StatusText(http.StatusUnauthorized),
		Message: message,
		Errors:  errors,
	})
}

func Forbidden(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusForbidden, response{
		Status:  http.StatusText(http.StatusForbidden),
		Message: message,
		Errors:  errors,
	})
}

func NotFound(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusNotFound, response{
		Status:  http.StatusText(http.StatusNotFound),
		Message: message,
		Errors:  errors,
	})
}

func InternalServerError(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusInternalServerError, response{
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: message,
		Errors:  errors,
	})
}

func FromError(c *gin.Context, err error) {
	if e, ok := err.(errs.Error); ok {
		DefaultResponse(c, e.Code(), e.Message(), e.Data(), e.Errors())
		return
	}
	DefaultResponse(c, http.StatusInternalServerError, err.Error(), nil, nil)
}
