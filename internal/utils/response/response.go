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

func DefaultResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, response{
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
	})
}

// TODO: Implement better error handling for bad request responses when fields fail validation.
func BadRequest(c *gin.Context, message ...string) {
	msg := "Internal Server Error"
	if len(message) > 0 {
		msg = message[0]
	}
	DefaultResponse(c, http.StatusBadRequest, msg, nil)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	DefaultResponse(c, http.StatusUnauthorized, message, nil)
}

func InternalServerError(c *gin.Context, message ...string) {
	msg := "Internal Server Error"
	if len(message) > 0 {
		msg = message[0]
	}
	DefaultResponse(c, http.StatusInternalServerError, msg, nil)
}

func Success(c *gin.Context, message string, data interface{}) {
	DefaultResponse(c, http.StatusOK, message, data)
}
