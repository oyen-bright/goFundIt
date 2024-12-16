package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
