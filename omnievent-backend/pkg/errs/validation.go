package errs

import (
	"net/http"
)

// Error codes related to validation (4000-4999)
var (
	ErrValidationFailed            = NewNormalError(NormalSubcategoryValidation, 0, http.StatusBadRequest, "validation failed")
	ErrParameterInvalid            = NewNormalError(NormalSubcategoryValidation, 1, http.StatusBadRequest, "parameter is invalid")
	ErrParameterRequired           = NewNormalError(NormalSubcategoryValidation, 2, http.StatusBadRequest, "parameter is required")
	ErrParameterTooLong            = NewNormalError(NormalSubcategoryValidation, 3, http.StatusBadRequest, "parameter is too long")
	ErrParameterTooShort           = NewNormalError(NormalSubcategoryValidation, 4, http.StatusBadRequest, "parameter is too short")
	ErrInvalidFormat               = NewNormalError(NormalSubcategoryValidation, 5, http.StatusBadRequest, "invalid format")
	ErrInvalidEmailFormat          = NewNormalError(NormalSubcategoryValidation, 6, http.StatusBadRequest, "invalid email format")
	ErrInvalidUsernameFormat       = NewNormalError(NormalSubcategoryValidation, 7, http.StatusBadRequest, "invalid username format")
	ErrInvalidPasswordFormat        = NewNormalError(NormalSubcategoryValidation, 8, http.StatusBadRequest, "invalid password format")
	ErrPasswordsDoNotMatch         = NewNormalError(NormalSubcategoryValidation, 9, http.StatusBadRequest, "passwords do not match")
	ErrInvalidJSONFormat           = NewNormalError(NormalSubcategoryValidation, 10, http.StatusBadRequest, "invalid JSON format")
	ErrFieldRequired              = NewNormalError(NormalSubcategoryValidation, 11, http.StatusBadRequest, "field is required")
	ErrFieldTooLong               = NewNormalError(NormalSubcategoryValidation, 12, http.StatusBadRequest, "field is too long")
	ErrFieldTooShort              = NewNormalError(NormalSubcategoryValidation, 13, http.StatusBadRequest, "field is too short")
	ErrInvalidRange               = NewNormalError(NormalSubcategoryValidation, 14, http.StatusBadRequest, "value is out of valid range")
)