package response

import (
	"strconv"

	"omnievent-backend/pkg/errs"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
}

type ErrorResponse struct {
	Success     bool   `json:"success"`
	ErrorCode   int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Path        string `json:"path,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, SuccessResponse{
		Success: true,
		Result:  data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Success:     false,
		ErrorMessage: message,
		Path:        c.Request.URL.Path,
	})
}

func ErrorWithCode(c *gin.Context, status int, code int, message string) {
	c.JSON(status, ErrorResponse{
		Success:     false,
		ErrorCode:   code,
		ErrorMessage: message,
		Path:        c.Request.URL.Path,
	})
}

func ErrorWithError(c *gin.Context, err error) {
	if e, ok := err.(*errs.Error); ok {
		c.JSON(e.HttpStatusCode, ErrorResponse{
			Success:     false,
			ErrorCode:   e.Code(),
			ErrorMessage: e.Message,
			Path:        c.Request.URL.Path,
		})
		return
	}
	c.JSON(500, ErrorResponse{
		Success:     false,
		ErrorCode:   100000,
		ErrorMessage: "internal server error",
		Path:        c.Request.URL.Path,
	})
}

func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"success": true,
		"message": message,
	})
}

func ValidationError(c *gin.Context, message string) {
	ErrorWithCode(c, 400, 2001, message)
}

func ServerError(c *gin.Context, err error) {
	if e, ok := err.(*errs.Error); ok {
		ErrorWithError(c, e)
		return
	}
	c.JSON(500, ErrorResponse{
		Success:     false,
		ErrorCode:   100000,
		ErrorMessage: "internal server error: " + err.Error(),
		Path:        c.Request.URL.Path,
	})
}

func BadRequest(c *gin.Context, message string) {
	ErrorWithCode(c, 400, 2002, message)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "unauthorized"
	}
	ErrorWithCode(c, 401, 1003, message)
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "forbidden"
	}
	ErrorWithCode(c, 403, 1000, message)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "resource not found"
	}
	ErrorWithCode(c, 404, 1001, message)
}

func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = "resource conflict"
	}
	ErrorWithCode(c, 409, 1002, message)
}

func StringToInt(s string, defaultVal int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return i
}
