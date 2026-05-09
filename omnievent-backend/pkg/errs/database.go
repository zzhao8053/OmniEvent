package errs

import (
	"net/http"
)

// Error codes related to database (100000+ system errors)
var (
	ErrDatabaseTypeInvalid      = NewSystemError(SystemSubcategoryDatabase, 0, http.StatusInternalServerError, "database type is invalid")
	ErrDatabaseHostInvalid      = NewSystemError(SystemSubcategoryDatabase, 1, http.StatusInternalServerError, "database host is invalid")
	ErrDatabaseIsNull           = NewSystemError(SystemSubcategoryDatabase, 2, http.StatusInternalServerError, "database cannot be null")
	ErrDatabaseOperationFailed  = NewSystemError(SystemSubcategoryDatabase, 3, http.StatusInternalServerError, "database operation failed")
	ErrDatabaseConnectionFailed = NewSystemError(SystemSubcategoryDatabase, 4, http.StatusInternalServerError, "database connection failed")
	ErrDatabaseQueryFailed      = NewSystemError(SystemSubcategoryDatabase, 5, http.StatusInternalServerError, "database query failed")
	ErrDatabaseTransactionFailed = NewSystemError(SystemSubcategoryDatabase, 6, http.StatusInternalServerError, "database transaction failed")
	ErrLoggingFileOpenFailed     = NewSystemError(SystemSubcategoryLogging, 0, http.StatusInternalServerError, "log file open failed")
	ErrLoggingFileCloseFailed    = NewSystemError(SystemSubcategoryLogging, 1, http.StatusInternalServerError, "log file close failed")
	ErrLoggingFileRenameFailed   = NewSystemError(SystemSubcategoryLogging, 2, http.StatusInternalServerError, "log file rename failed")
	ErrSystemError               = NewSystemError(SystemSubcategoryDefault, 0, http.StatusInternalServerError, "system error")
	ErrInvalidStorageType        = NewSystemError(SystemSubcategorySetting, 0, http.StatusInternalServerError, "invalid storage type")
)