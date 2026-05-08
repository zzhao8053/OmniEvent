package errs

import (
	"fmt"
	"net/http"
)

type ErrorCategory int

const (
	CategorySystem ErrorCategory = 1
	CategoryNormal ErrorCategory = 2
)

type ErrorSubCategory int

const (
	// System subcategories
	SubCategoryDatabase ErrorSubCategory = 1
	SubCategoryMail     ErrorSubCategory = 2
	SubCategoryLogging  ErrorSubCategory = 3
	SubCategoryCron     ErrorSubCategory = 4
	SubCategorySecurity ErrorSubCategory = 5
	SubCategoryStorage  ErrorSubCategory = 6
	SubCategoryConfig   ErrorSubCategory = 7

	// Normal subcategories
	SubCategoryUser        ErrorSubCategory = 1
	SubCategoryToken       ErrorSubCategory = 2
	SubCategoryTwoFactor   ErrorSubCategory = 3
	SubCategoryAccount     ErrorSubCategory = 4
	SubCategoryTransaction ErrorSubCategory = 5
	SubCategoryCategory    ErrorSubCategory = 6
	SubCategoryTag         ErrorSubCategory = 7
	SubCategoryTemplate    ErrorSubCategory = 8
	SubCategoryPicture     ErrorSubCategory = 9
	SubCategoryData        ErrorSubCategory = 10
	SubCategoryValidation  ErrorSubCategory = 11
)

type Error struct {
	Category      ErrorCategory
	SubCategory   ErrorSubCategory
	Index         int
	HttpStatusCode int
	Message       string
	BaseError     error
	Context       map[string]interface{}
}

func (e *Error) Error() string {
	if e.BaseError != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.BaseError)
	}
	return e.Message
}

func (e *Error) Code() int {
	return int(e.Category)*100000 + int(e.SubCategory)*1000 + e.Index
}

func (e *Error) WithContext(key string, value interface{}) *Error {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

func (e *Error) WithBaseError(base error) *Error {
	e.BaseError = base
	return e
}

type MultiErrors struct {
	Errors []*Error
}

func (m *MultiErrors) Add(err *Error) {
	m.Errors = append(m.Errors, err)
}

func (m *MultiErrors) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}
	return fmt.Sprintf("%d errors occurred", len(m.Errors))
}

func NewSystemError(subCategory ErrorSubCategory, index int, httpStatus int, message string) *Error {
	return &Error{
		Category:       CategorySystem,
		SubCategory:    subCategory,
		Index:          index,
		HttpStatusCode: httpStatus,
		Message:        message,
	}
}

func NewNormalError(subCategory ErrorSubCategory, index int, httpStatus int, message string) *Error {
	return &Error{
		Category:       CategoryNormal,
		SubCategory:    subCategory,
		Index:          index,
		HttpStatusCode: httpStatus,
		Message:        message,
	}
}

func Or(err error, defaultErr error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*Error); ok {
		return err
	}
	return defaultErr
}

func IsCustomError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

func GetHttpStatusCode(err error) int {
	if e, ok := err.(*Error); ok {
		return e.HttpStatusCode
	}
	return http.StatusInternalServerError
}
