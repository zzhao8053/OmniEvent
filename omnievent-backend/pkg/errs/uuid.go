package errs

// Uuid errors (subCategory 40)
var (
	ErrInvalidUuidMode = NewNormalError(40, 1, 400, "invalid uuid generator mode")
)