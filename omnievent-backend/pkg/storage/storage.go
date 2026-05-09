package storage

import "omnievent-backend/pkg/context"

// ObjectStorage represents an object storage to store file object
type ObjectStorage interface {
	Exists(ctx *context.WebContext, path string) (bool, error)
	Read(ctx *context.WebContext, path string) (ObjectInStorage, error)
	Save(ctx *context.WebContext, path string, object ObjectInStorage) error
	Delete(ctx *context.WebContext, path string) error
}