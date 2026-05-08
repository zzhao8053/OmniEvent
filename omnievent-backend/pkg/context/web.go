package context

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyRequestID    = "RequestID"
	ContextKeyToken        = "Token"
	ContextKeyTokenClaims  = "TokenClaims"
	ContextKeyUID          = "UID"
	ContextKeyResponseError = "ResponseError"
)

type UserTokenClaims struct {
	UID         int64
	UserTokenID int64
	TokenType   int
	ExpiresAt   int64
}

type WebContext struct {
	*gin.Context
}

func WrapWebContext(c *gin.Context) *WebContext {
	return &WebContext{Context: c}
}

func (c *WebContext) SetContextID(id string) {
	c.Set(ContextKeyRequestID, id)
}

func (c *WebContext) GetContextID() string {
	if v, exists := c.Get(ContextKeyRequestID); exists {
		return v.(string)
	}
	return ""
}

func (c *WebContext) SetTextualToken(token string) {
	c.Set(ContextKeyToken, token)
}

func (c *WebContext) GetTextualToken() string {
	if v, exists := c.Get(ContextKeyToken); exists {
		return v.(string)
	}
	return ""
}

func (c *WebContext) SetTokenClaims(claims *UserTokenClaims) {
	c.Set(ContextKeyTokenClaims, claims)
}

func (c *WebContext) GetTokenClaims() *UserTokenClaims {
	if v, exists := c.Get(ContextKeyTokenClaims); exists {
		return v.(*UserTokenClaims)
	}
	return nil
}

func (c *WebContext) GetCurrentUID() int64 {
	if claims := c.GetTokenClaims(); claims != nil {
		return claims.UID
	}
	return 0
}

func (c *WebContext) SetResponseError(err interface{}) {
	c.Set(ContextKeyResponseError, err)
}

func (c *WebContext) GetResponseError() interface{} {
	if v, exists := c.Get(ContextKeyResponseError); exists {
		return v
	}
	return nil
}

func (c *WebContext) GetTokenStringFromHeader() string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

func (c *WebContext) GetTokenStringFromQueryString() string {
	return c.Query("token")
}

func (c *WebContext) GetTokenStringFromCookie() string {
	token, _ := c.Cookie("token")
	return token
}

func (c *WebContext) SetTokenStringToCookie(token string, expiredAt int64, path string) {
	maxAge := 0
	if expiredAt > time.Now().Unix() {
		maxAge = int(expiredAt - time.Now().Unix())
	}
	c.SetCookie("token", token, maxAge, path, "", false, true)
}

func (c *WebContext) DeleteTokenCookie() {
	c.SetCookie("token", "", -1, "/", "", false, true)
}

func (c *WebContext) GetClientLocale() string {
	return c.GetHeader("Accept-Language")
}

func (c *WebContext) GetClientTimezoneOffset() int {
	offsetStr := c.GetHeader("X-Timezone-Offset")
	if offsetStr == "" {
		return 0
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0
	}
	return offset
}

func (c *WebContext) GetClientIP() string {
	return c.ClientIP()
}

func (c *WebContext) GetUserAgent() string {
	return c.GetHeader("User-Agent")
}
