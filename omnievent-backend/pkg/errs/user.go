package errs

import (
	"net/http"
)

// Error codes related to users (2000-2999)
var (
	ErrLoginNameInvalid         = NewNormalError(NormalSubcategoryUser, 0, http.StatusUnauthorized, "login name is invalid")
	ErrLoginNameOrPasswordInvalid = NewNormalError(NormalSubcategoryUser, 1, http.StatusUnauthorized, "login name or password is invalid")
	ErrLoginNameOrPasswordWrong   = NewNormalError(NormalSubcategoryUser, 2, http.StatusUnauthorized, "login name or password is wrong")
	ErrUserIdInvalid              = NewNormalError(NormalSubcategoryUser, 3, http.StatusBadRequest, "user id is invalid")
	ErrUsernameIsEmpty            = NewNormalError(NormalSubcategoryUser, 4, http.StatusBadRequest, "username is empty")
	ErrEmailIsEmpty               = NewNormalError(NormalSubcategoryUser, 5, http.StatusBadRequest, "email is empty")
	ErrNicknameIsEmpty            = NewNormalError(NormalSubcategoryUser, 6, http.StatusBadRequest, "nickname is empty")
	ErrPasswordIsEmpty            = NewNormalError(NormalSubcategoryUser, 7, http.StatusBadRequest, "password is empty")
	ErrUserNotFound               = NewNormalError(NormalSubcategoryUser, 8, http.StatusBadRequest, "user not found")
	ErrUserPasswordWrong          = NewNormalError(NormalSubcategoryUser, 9, http.StatusBadRequest, "password is wrong")
	ErrUsernameAlreadyExists      = NewNormalError(NormalSubcategoryUser, 10, http.StatusBadRequest, "username already exists")
	ErrUserEmailAlreadyExists     = NewNormalError(NormalSubcategoryUser, 11, http.StatusBadRequest, "email already exists")
	ErrUserRegistrationNotAllowed = NewNormalError(NormalSubcategoryUser, 12, http.StatusBadRequest, "user registration not allowed")
	ErrUserIsDisabled             = NewNormalError(NormalSubcategoryUser, 13, http.StatusBadRequest, "user is disabled")
	ErrEmailIsInvalid             = NewNormalError(NormalSubcategoryUser, 14, http.StatusBadRequest, "email is invalid")
	ErrEmailIsEmptyOrInvalid      = NewNormalError(NormalSubcategoryUser, 15, http.StatusBadRequest, "email is empty or invalid")
	ErrNewPasswordEqualsOldInvalid = NewNormalError(NormalSubcategoryUser, 16, http.StatusBadRequest, "new password equals old password")
	ErrNotPermittedToPerformThisAction = NewNormalError(NormalSubcategoryUser, 17, http.StatusBadRequest, "not permitted to perform this action")
	ErrCannotLoginByPassword      = NewNormalError(NormalSubcategoryUser, 18, http.StatusBadRequest, "cannot login by password")
	ErrUserNameIsInvalid          = NewNormalError(NormalSubcategoryUser, 19, http.StatusBadRequest, "user name is invalid")
	ErrNickNameIsInvalid          = NewNormalError(NormalSubcategoryUser, 20, http.StatusBadRequest, "nick name is invalid")
	ErrMinUsernameLength          = NewNormalError(NormalSubcategoryUser, 21, http.StatusBadRequest, "username is too short")
	ErrMaxUsernameLength          = NewNormalError(NormalSubcategoryUser, 22, http.StatusBadRequest, "username is too long")
	ErrMinPasswordLength          = NewNormalError(NormalSubcategoryUser, 23, http.StatusBadRequest, "password is too short")
)