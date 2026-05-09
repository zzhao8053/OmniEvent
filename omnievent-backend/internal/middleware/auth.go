package middleware

import (
	"strings"

	"omnievent-backend/internal/model"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
	"omnievent-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TokenSourceType represents token source
type TokenSourceType byte

// Token source types
const (
	TOKEN_SOURCE_TYPE_HEADER   TokenSourceType = 1
	TOKEN_SOURCE_TYPE_ARGUMENT TokenSourceType = 2
	TOKEN_SOURCE_TYPE_COOKIE   TokenSourceType = 3
)

// Claims represents JWT claims for OmniEvent
type Claims struct {
	UID         int64 `json:"uid"`
	UserTokenID int64 `json:"user_token_id"`
	TokenType   int   `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTAuthorization verifies whether current request is valid by jwt token in header
func JWTAuthorization() gin.HandlerFunc {
	return jwtAuthorization(TOKEN_SOURCE_TYPE_HEADER)
}

// JWTAuthorizationByQueryString verifies whether current request is valid by jwt token in query string
func JWTAuthorizationByQueryString() gin.HandlerFunc {
	return jwtAuthorization(TOKEN_SOURCE_TYPE_ARGUMENT)
}

// JWTAuthorizationByCookie verifies whether current request is valid by jwt token in cookie
func JWTAuthorizationByCookie() gin.HandlerFunc {
	return jwtAuthorization(TOKEN_SOURCE_TYPE_COOKIE)
}

// JWTTwoFactorAuthorization verifies whether current request is valid by 2fa passcode
func JWTTwoFactorAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		webCtx := context.WrapWebContext(c)
		claims, err := getTokenClaims(webCtx, TOKEN_SOURCE_TYPE_HEADER)

		if err != nil {
			response.ErrorWithError(c, err)
			c.Abort()
			return
		}

		if claims.TokenType != model.TokenTypeRequire2FA {
			response.Unauthorized(c, "user token is not require two-factor authorization")
			c.Abort()
			return
		}

		webCtx.SetTokenClaims(&context.UserTokenClaims{
			UID:         claims.UID,
			UserTokenID: claims.UserTokenID,
			TokenType:   claims.TokenType,
		})
		c.Next()
	}
}

// JWTEmailVerifyAuthorization verifies whether current request is for email verification
func JWTEmailVerifyAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		webCtx := context.WrapWebContext(c)
		claims, err := getTokenClaims(webCtx, TOKEN_SOURCE_TYPE_ARGUMENT)

		if err != nil {
			response.ErrorWithError(c, err)
			c.Abort()
			return
		}

		if claims.TokenType != model.TokenTypeEmailVerify {
			response.Unauthorized(c, "token is not for email verification")
			c.Abort()
			return
		}

		webCtx.SetTokenClaims(&context.UserTokenClaims{
			UID:         claims.UID,
			UserTokenID: claims.UserTokenID,
			TokenType:   claims.TokenType,
		})
		c.Next()
	}
}

// JWTResetPasswordAuthorization verifies whether current request is for password reset
func JWTResetPasswordAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		webCtx := context.WrapWebContext(c)
		claims, err := getTokenClaims(webCtx, TOKEN_SOURCE_TYPE_ARGUMENT)

		if err != nil {
			response.ErrorWithError(c, err)
			c.Abort()
			return
		}

		if claims.TokenType != model.TokenTypePasswordReset {
			response.Unauthorized(c, "token is not for password reset")
			c.Abort()
			return
		}

		webCtx.SetTokenClaims(&context.UserTokenClaims{
			UID:         claims.UID,
			UserTokenID: claims.UserTokenID,
			TokenType:   claims.TokenType,
		})
		c.Next()
	}
}

// JWTMCPAuthorization verifies whether current request is valid by jwt mcp token in header
func JWTMCPAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		webCtx := context.WrapWebContext(c)
		claims, err := getTokenClaims(webCtx, TOKEN_SOURCE_TYPE_HEADER)

		if err != nil {
			response.ErrorWithError(c, err)
			c.Abort()
			return
		}

		if claims.TokenType != model.TokenTypeMCP {
			response.Unauthorized(c, "token type is not mcp token")
			c.Abort()
			return
		}

		webCtx.SetTokenClaims(&context.UserTokenClaims{
			UID:         claims.UID,
			UserTokenID: claims.UserTokenID,
			TokenType:   claims.TokenType,
		})
		c.Next()
	}
}

func jwtAuthorization(source TokenSourceType) gin.HandlerFunc {
	return func(c *gin.Context) {
		webCtx := context.WrapWebContext(c)
		claims, err := getTokenClaims(webCtx, source)

		if err != nil {
			response.ErrorWithError(c, err)
			c.Abort()
			return
		}

		if claims.TokenType == model.TokenTypeRequire2FA {
			response.Unauthorized(c, "token requires two-factor authorization")
			c.Abort()
			return
		}

		if claims.TokenType != model.TokenTypeNormal {
			response.Unauthorized(c, "token type is invalid")
			c.Abort()
			return
		}

		webCtx.SetTokenClaims(&context.UserTokenClaims{
			UID:         claims.UID,
			UserTokenID: claims.UserTokenID,
			TokenType:   claims.TokenType,
		})
		c.Next()
	}
}

func getTokenClaims(webCtx *context.WebContext, source TokenSourceType) (*Claims, *errs.Error) {
	tokenString := ""

	if source == TOKEN_SOURCE_TYPE_ARGUMENT {
		tokenString = webCtx.GetTokenStringFromQueryString()
	} else if source == TOKEN_SOURCE_TYPE_COOKIE {
		tokenString = webCtx.GetTokenStringFromCookie()
	} else {
		tokenString = webCtx.GetTokenStringFromHeader()
	}

	if tokenString == "" {
		return nil, errs.ErrTokenIsEmpty
	}

	token, err := parseToken(tokenString)
	if err != nil {
		return nil, errs.ErrTokenMalformed
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errs.ErrCurrentInvalidToken
	}

	if claims.UID <= 0 {
		return nil, errs.ErrCurrentInvalidToken
	}

	return claims, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrTokenSignatureInvalid
		}
		return []byte("omnievent-secret-key"), nil
	})
	return token, err
}

// getTokenStringFromHeader extracts token from Authorization header
func getTokenStringFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && strings.ToLower(authHeader[:7]) == "bearer " {
		return authHeader[7:]
	}
	return ""
}