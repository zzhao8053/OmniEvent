package errs

import (
	"net/http"
)

// Error codes related to tokens (3000-3999)
var (
	ErrTokenGenerating           = NewNormalError(NormalSubcategoryToken, 0, http.StatusInternalServerError, "failed to generate token")
	ErrUnauthorizedAccess        = NewNormalError(NormalSubcategoryToken, 1, http.StatusUnauthorized, "unauthorized access")
	ErrCurrentInvalidToken       = NewNormalError(NormalSubcategoryToken, 2, http.StatusUnauthorized, "current token is invalid")
	ErrCurrentTokenExpired       = NewNormalError(NormalSubcategoryToken, 3, http.StatusUnauthorized, "current token is expired")
	ErrCurrentInvalidTokenType   = NewNormalError(NormalSubcategoryToken, 4, http.StatusUnauthorized, "current token type is invalid")
	ErrInvalidToken              = NewNormalError(NormalSubcategoryToken, 5, http.StatusBadRequest, "token is invalid")
	ErrInvalidTokenId            = NewNormalError(NormalSubcategoryToken, 6, http.StatusBadRequest, "token id is invalid")
	ErrInvalidUserTokenId        = NewNormalError(NormalSubcategoryToken, 7, http.StatusBadRequest, "user token id is invalid")
	ErrTokenRecordNotFound       = NewNormalError(NormalSubcategoryToken, 8, http.StatusBadRequest, "token is not found")
	ErrTokenExpired              = NewNormalError(NormalSubcategoryToken, 9, http.StatusBadRequest, "token is expired")
	ErrTokenIsEmpty              = NewNormalError(NormalSubcategoryToken, 10, http.StatusBadRequest, "token is empty")
	ErrPasswordResetTokenInvalid = NewNormalError(NormalSubcategoryToken, 11, http.StatusBadRequest, "password reset token is invalid or expired")
	ErrEmailVerifyTokenInvalid   = NewNormalError(NormalSubcategoryToken, 12, http.StatusBadRequest, "email verify token is invalid or expired")
	ErrTokenNotProvided          = NewNormalError(NormalSubcategoryToken, 13, http.StatusUnauthorized, "token not provided")
	ErrTokenMalformed           = NewNormalError(NormalSubcategoryToken, 14, http.StatusUnauthorized, "token is malformed")
	ErrTokenSignatureInvalid     = NewNormalError(NormalSubcategoryToken, 15, http.StatusUnauthorized, "token signature is invalid")
)