package errs

var (
	ErrUserNotFound = NewNormalError(SubCategoryUser, 1, 404, "user not found")
	ErrUserExists   = NewNormalError(SubCategoryUser, 2, 409, "user already exists")
	ErrInvalidCredentials = NewNormalError(SubCategoryUser, 3, 401, "invalid credentials")
	ErrEmailExists  = NewNormalError(SubCategoryUser, 4, 409, "email already exists")
	ErrUsernameExists = NewNormalError(SubCategoryUser, 5, 409, "username already exists")
	ErrUserDisabled = NewNormalError(SubCategoryUser, 6, 403, "user account is disabled")
	ErrEmailNotVerified = NewNormalError(SubCategoryUser, 7, 403, "email not verified")
)
