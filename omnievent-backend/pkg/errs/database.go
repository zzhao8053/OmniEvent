package errs

var (
	ErrDatabaseOperation = NewSystemError(SubCategoryDatabase, 1, 500, "database operation failed")
	ErrDatabaseNotFound  = NewSystemError(SubCategoryDatabase, 2, 404, "record not found in database")
	ErrDatabaseDuplicate = NewSystemError(SubCategoryDatabase, 3, 409, "duplicate record in database")
)
