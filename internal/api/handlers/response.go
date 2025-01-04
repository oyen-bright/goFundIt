package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/pkg/errs"
)

// SuccessResponse represents a successful API response
// @Description Successful API response structure
type SuccessResponse struct {
	// Status will always be "OK"
	// @example OK
	Status string `json:"status" example:"OK" enums:"OK"`

	// Response message describing the success
	// @example Operation completed successfully
	Message string `json:"message"`

	// Response data (optional)
	// @example null
	Data interface{} `json:"data,omitempty"`
}

// BadRequestResponse represents a 400 error response
// @Description Bad request error response structure
type BadRequestResponse struct {
	// Status will always be "Bad Request"
	// @example Bad Request
	Status string `json:"status" example:"Bad Request" enums:"Bad Request"`

	// Error message
	// @example Invalid input provided
	Message string `json:"message"`

	// Validation errors (optional)
	// @example [{"field":"field","error":"must be a valid data"}]
	Errors interface{} `json:"errors,omitempty"`
}

// UnauthorizedResponse represents a 401 error response
// @Description Unauthorized error response structure
type UnauthorizedResponse struct {
	// Status will always be false for unauthorized
	// @example false
	Status string `json:"status" example:"Unauthorized" enums:"Unauthorized"`

	// Will always be "Unauthorized"
	// @example Unauthorized
	Message string `json:"message" example:"Unauthorized" enums:"Unauthorized"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field string `json:"field" example:"field"`
	Error string `json:"error" example:"must be a valid field"`
}

// Response represents the standard API response structure
type response struct {
	// Response status
	// @example OK
	Status string `json:"status"`

	// Response message
	// @example Operation completed successfully
	Message string `json:"message"`

	// Response data (optional)
	Data interface{} `json:"data,omitempty"`

	// Error details (optional)
	Errors interface{} `json:"errors,omitempty"`
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
	c.JSON(http.StatusBadRequest, BadRequestResponse{
		Status:  http.StatusText(http.StatusBadRequest),
		Message: message,
		Errors:  errors,
	})
}

func Unauthorized(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnauthorized, UnauthorizedResponse{
		Status:  http.StatusText(http.StatusUnauthorized),
		Message: message,
		// Errors:  errors,
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
