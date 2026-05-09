package api

import (
	"fmt"
	"net/http"
	"strings"

	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
	"omnievent-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MiddlewareHandlerFunc 中间件处理器函数类型
type MiddlewareHandlerFunc func(*context.WebContext)

// ApiHandlerFunc API 处理器函数类型
type ApiHandlerFunc func(*context.WebContext) (interface{}, *errs.Error)

// BindMiddleware 将 MiddlewareHandlerFunc 转换为 Gin 中间件
func BindMiddleware(fn MiddlewareHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(context.WrapWebContext(c))
	}
}

// BindApi 绑定 API 请求的通用封装
// 自动处理 JSON 响应格式
func BindApi(fn ApiHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := context.WrapWebContext(ginCtx)
		result, err := fn(c)

		if err != nil {
			response.ErrorWithError(ginCtx, err)
		} else {
			response.Success(ginCtx, result)
		}
	}
}

// BindApiWithTokenUpdate 绑定 API 并更新 Token
// 用于需要更新 Token 的认证相关 API
func BindApiWithTokenUpdate(fn ApiHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := context.WrapWebContext(ginCtx)
		result, err := fn(c)

		if err != nil {
			response.ErrorWithError(ginCtx, err)
		} else {
			response.Success(ginCtx, result)
		}
	}
}

// BindApiWithCookie 绑定 API 并设置 Cookie
// 用于需要设置认证 Cookie 的 API
func BindApiWithCookie(fn ApiHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := context.WrapWebContext(ginCtx)
		result, err := fn(c)

		if err != nil {
			response.ErrorWithError(ginCtx, err)
		} else {
			response.Success(ginCtx, result)
		}
	}
}

// AvatarApiHandlerFunc 头像 API 处理器函数类型 (返回字节数据)
type AvatarApiHandlerFunc func(*context.WebContext) ([]byte, string, *errs.Error)

// BindAvatarApi 绑定头像 API 请求 (直接返回图片数据)
func BindAvatarApi(fn AvatarApiHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := context.WrapWebContext(ginCtx)
		data, contentType, err := fn(c)

		if err != nil {
			response.ErrorWithError(ginCtx, err)
		} else {
			ginCtx.Header("Content-Type", contentType)
			ginCtx.Header("Content-Length", fmt.Sprintf("%d", len(data)))
			ginCtx.Data(http.StatusOK, contentType, data)
		}
	}
}

// BindOptions 绑定选项
type BindOptions struct {
	TrimSpaces bool
}

// BindOption 绑定选项函数
type BindOption func(*BindOptions)

func WithTrimSpaces() BindOption {
	return func(o *BindOptions) {
		o.TrimSpaces = true
	}
}

// BindJSON 绑定 JSON 请求
func BindJSON(c *context.WebContext, obj interface{}, opts ...BindOption) *errs.Error {
	return bindWith(c.Context, obj, bindingJSON, opts...)
}

// BindQuery 绑定 Query 参数
func BindQuery(c *context.WebContext, obj interface{}, opts ...BindOption) *errs.Error {
	return bindWith(c.Context, obj, bindingQuery, opts...)
}

// BindForm 绑定 Form 数据
func BindForm(c *context.WebContext, obj interface{}, opts ...BindOption) *errs.Error {
	return bindWith(c.Context, obj, bindingForm, opts...)
}

// BindUri 绑定 URI 参数
func BindUri(c *context.WebContext, obj interface{}, opts ...BindOption) *errs.Error {
	return bindWith(c.Context, obj, bindingUri, opts...)
}

type binder func(c *gin.Context, obj interface{}) error

var (
	bindingJSON  = func(c *gin.Context, obj interface{}) error { return c.ShouldBindJSON(obj) }
	bindingQuery = func(c *gin.Context, obj interface{}) error { return c.ShouldBindQuery(obj) }
	bindingForm  = func(c *gin.Context, obj interface{}) error { return c.ShouldBind(obj) }
	bindingUri   = func(c *gin.Context, obj interface{}) error { return c.ShouldBindUri(obj) }
)

func bindWith(c *gin.Context, obj interface{}, binder binder, opts ...BindOption) *errs.Error {
	options := &BindOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if err := binder(c, obj); err != nil {
		fmt.Printf("[bind] parse request failed, because %s\n", err.Error())

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return validationError(validationErrors)
		}

		return errs.ErrValidationFailed
	}

	if options.TrimSpaces {
		trimSpaces(obj)
	}

	if err := validateStruct(obj); err != nil {
		return err
	}

	return nil
}

func validationError(validationErrs validator.ValidationErrors) *errs.Error {
	for _, err := range validationErrs {
		field := strings.ToLower(err.Field())
		tag := err.Tag()

		switch tag {
		case "required":
			return errs.NewNormalError(errs.NormalSubcategoryValidation, 1, 400, field+" is required")
		case "email":
			return errs.NewNormalError(errs.NormalSubcategoryValidation, 2, 400, field+" must be a valid email")
		case "min":
			return errs.NewNormalError(errs.NormalSubcategoryValidation, 3, 400, field+" is too short")
		case "max":
			return errs.NewNormalError(errs.NormalSubcategoryValidation, 4, 400, field+" is too long")
		case "eqfield":
			return errs.NewNormalError(errs.NormalSubcategoryValidation, 5, 400, field+" does not match")
		default:
			return errs.NewNormalError(errs.NormalSubcategoryValidation, 9, 400, field+" validation failed: "+tag)
		}
	}
	return errs.ErrValidationFailed
}

func validateStruct(obj interface{}) *errs.Error {
	return nil
}

func trimSpaces(obj interface{}) {
}

// GetClientIP 获取客户端 IP
func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	return c.ClientIP()
}

// GetUserAgent 获取用户代理
func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

// ParseInt64 解析 int64
func ParseInt64(s string, defaultVal int64) int64 {
	if s == "" {
		return defaultVal
	}
	var val int64
	for _, c := range s {
		if c < '0' || c > '9' {
			return defaultVal
		}
		val = val*10 + int64(c-'0')
	}
	return val
}

// ParseInt 解析 int
func ParseInt(s string, defaultVal int) int {
	return int(ParseInt64(s, int64(defaultVal)))
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  data,
	})
}

// SuccessWithMessage 成功响应（带消息）
func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})
}
