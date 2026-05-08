package middleware

import (
	"net/http"
	"strings"

	"omnievent-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UID         int64 `json:"uid"`
	UserTokenID int64 `json:"user_token_id"`
	TokenType   int   `json:"token_type"`
	jwt.RegisteredClaims
}

func JWTAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		c.Set("token", tokenString)
		c.Next()
	}
}
