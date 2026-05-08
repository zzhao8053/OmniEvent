package errs

var (
	ErrTokenNotFound       = NewNormalError(SubCategoryToken, 1, 404, "token not found")
	ErrTokenExpired        = NewNormalError(SubCategoryToken, 2, 401, "token expired")
	ErrTokenInvalid        = NewNormalError(SubCategoryToken, 3, 401, "invalid token")
	ErrTokenTypeMismatch   = NewNormalError(SubCategoryToken, 4, 401, "token type mismatch")
	ErrRefreshTokenExpired  = NewNormalError(SubCategoryToken, 5, 401, "refresh token expired")
)
