package errs

var (
	ErrValidationFailed = NewNormalError(SubCategoryValidation, 1, 400, "validation failed")
	ErrInvalidRequest   = NewNormalError(SubCategoryValidation, 2, 400, "invalid request")
)
