package storage

import (
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
)

const avatarPathPrefix = "avatar"
const transactionPicturePathPrefix = "transaction"

// StorageContainer contains the current object storage
type StorageContainer struct {
	avatarCurrentStorage             ObjectStorage
	transactionPictureCurrentStorage ObjectStorage
}

// Initialize a object storage container singleton instance
var (
	Container     = &StorageContainer{}
	storageConfig *StorageConfig
)

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type                  string
	LocalFileSystemPath  string
	MinIOEndpoint        string
	MinIOLocation        string
	MinIOAccessKeyID     string
	MinIOSecretAccessKey string
	MinIOUseSSL          bool
	MinIOSkipTLSVerify   bool
	MinIOBucket          string
	MinIORootPath        string
	WebDAVURL            string
	WebDAVUsername       string
	WebDAVPassword       string
	WebDAVRootPath       string
	WebDAVRequestTimeout int
	WebDAVProxy          string
	WebDAVSkipTLSVerify  bool
}

// SetStorageConfig sets the storage configuration
func SetStorageConfig(config *StorageConfig) {
	storageConfig = config
}

// GetStorageConfig returns the storage configuration
func GetStorageConfig() *StorageConfig {
	return storageConfig
}

// InitializeStorageContainer initializes the current object storage according to the config
func InitializeStorageContainer() error {
	if storageConfig == nil {
		return errs.ErrSystemError
	}

	if storageConfig.Type == "local_filesystem" {
		avatarStorage, err := newObjectStorage(avatarPathPrefix)

		if err != nil {
			return err
		}

		Container.avatarCurrentStorage = avatarStorage
		Container.transactionPictureCurrentStorage = avatarStorage
	} else if storageConfig.Type == "minio" {
		avatarStorage, err := newObjectStorage(avatarPathPrefix)

		if err != nil {
			return err
		}

		Container.avatarCurrentStorage = avatarStorage
		Container.transactionPictureCurrentStorage = avatarStorage
	} else if storageConfig.Type == "webdav" {
		avatarStorage, err := newObjectStorage(avatarPathPrefix)

		if err != nil {
			return err
		}

		Container.avatarCurrentStorage = avatarStorage
		Container.transactionPictureCurrentStorage = avatarStorage
	}

	return nil
}

// ExistsAvatar returns whether the avatar file exists from the current avatar object storage
func (s *StorageContainer) ExistsAvatar(ctx *context.WebContext, path string) (bool, error) {
	if s.avatarCurrentStorage == nil {
		return false, errs.ErrSystemError
	}

	return s.avatarCurrentStorage.Exists(ctx, path)
}

// ReadAvatar returns the avatar file from the current avatar object storage
func (s *StorageContainer) ReadAvatar(ctx *context.WebContext, path string) (ObjectInStorage, error) {
	if s.avatarCurrentStorage == nil {
		return nil, errs.ErrSystemError
	}

	return s.avatarCurrentStorage.Read(ctx, path)
}

// SaveAvatar returns whether save the avatar file into the current avatar object storage successfully
func (s *StorageContainer) SaveAvatar(ctx *context.WebContext, path string, object ObjectInStorage) error {
	if s.avatarCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.avatarCurrentStorage.Save(ctx, path, object)
}

// DeleteAvatar returns whether delete the avatar file from the current avatar object storage successfully
func (s *StorageContainer) DeleteAvatar(ctx *context.WebContext, path string) error {
	if s.avatarCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.avatarCurrentStorage.Delete(ctx, path)
}

// ExistsTransactionPicture returns whether the transaction picture file exists from the current transaction picture object storage
func (s *StorageContainer) ExistsTransactionPicture(ctx *context.WebContext, path string) (bool, error) {
	if s.transactionPictureCurrentStorage == nil {
		return false, errs.ErrSystemError
	}

	return s.transactionPictureCurrentStorage.Exists(ctx, path)
}

// ReadTransactionPicture returns the transaction picture file from the current transaction picture object storage
func (s *StorageContainer) ReadTransactionPicture(ctx *context.WebContext, path string) (ObjectInStorage, error) {
	if s.transactionPictureCurrentStorage == nil {
		return nil, errs.ErrSystemError
	}

	return s.transactionPictureCurrentStorage.Read(ctx, path)
}

// SaveTransactionPicture returns whether save the transaction picture file into the current transaction picture object storage successfully
func (s *StorageContainer) SaveTransactionPicture(ctx *context.WebContext, path string, object ObjectInStorage) error {
	if s.transactionPictureCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.transactionPictureCurrentStorage.Save(ctx, path, object)
}

// DeleteTransactionPicture returns whether delete the transaction picture file from the current transaction picture object storage successfully
func (s *StorageContainer) DeleteTransactionPicture(ctx *context.WebContext, path string) error {
	if s.transactionPictureCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.transactionPictureCurrentStorage.Delete(ctx, path)
}

func newObjectStorage(pathPrefix string) (ObjectStorage, error) {
	if storageConfig.Type == "local_filesystem" {
		return NewLocalFileSystemObjectStorage(pathPrefix)
	} else if storageConfig.Type == "minio" {
		return NewMinIOObjectStorage(pathPrefix)
	} else if storageConfig.Type == "webdav" {
		return NewWebDAVObjectStorage(pathPrefix)
	}

	return nil, errs.ErrSystemError
}