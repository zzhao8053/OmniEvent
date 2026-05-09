package errs

// ErrorCategory represents error category
type ErrorCategory = int32

// Error categories
const (
	CATEGORY_SYSTEM ErrorCategory = 1
	CATEGORY_NORMAL ErrorCategory = 2
)

// Sub categories of system error
const (
	SystemSubcategoryDefault  ErrorCategory = 0
	SystemSubcategorySetting  ErrorCategory = 1
	SystemSubcategoryDatabase ErrorCategory = 2
	SystemSubcategoryMail     ErrorCategory = 3
	SystemSubcategoryLogging  ErrorCategory = 4
	SystemSubcategoryCron     ErrorCategory = 5
)

// Sub categories of normal error
const (
	NormalSubcategoryGlobal      ErrorCategory = 0
	NormalSubcategoryUser       ErrorCategory = 1
	NormalSubcategoryToken      ErrorCategory = 2
	NormalSubcategoryValidation ErrorCategory = 3
	NormalSubcategoryActivity   ErrorCategory = 4
)