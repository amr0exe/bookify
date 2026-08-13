package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Unifed way of handling response

type Body struct {
	Success bool         `json:"success"`
	Data    any          `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
	Meta    any          `json:"meta,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Responder interface {
	Success(c *gin.Context, statusCode int, data any)
	SuccessWithMeta(c *gin.Context, statusCode int, data any, meta any)
	Error(c *gin.Context, statusCode int, code string, message string, details any)
	ValidationError(c *gin.Context, err error)
}

type jsonResponder struct{}

func NewResponder() Responder {
	return &jsonResponder{}
}

func (r *jsonResponder) Success(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, Body{
		Success: true,
		Data:    data,
	})
}

func (r *jsonResponder) SuccessWithMeta(c *gin.Context, statusCode int, data any, meta any) {
	c.JSON(statusCode, Body{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func (r *jsonResponder) Error(c *gin.Context, statusCode int, code string, message string, details any) {
	c.JSON(statusCode, Body{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func (r *jsonResponder) ValidationError(c *gin.Context, err error) {
	details := formatValidationError(err)

	r.Error(
		c,
		http.StatusUnprocessableEntity,
		"VALIDATION_ERROR",
		"Invalid request payload",
		details,
	)
}
