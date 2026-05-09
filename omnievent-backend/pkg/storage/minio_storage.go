package storage

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	webcontext "omnievent-backend/pkg/context"
)

// MinIOObjectStorage represents MinIO object storage
type MinIOObjectStorage struct {
	minIOClient *minio.Client
	bucket      string
	rootPath    string
}

// NewMinIOObjectStorage returns a MinIO object storage
func NewMinIOObjectStorage(pathPrefix string) (*MinIOObjectStorage, error) {
	config := GetStorageConfig()

	minIOClient, err := minio.New(config.MinIOEndpoint, &minio.Options{
		Region:    config.MinIOLocation,
		Creds:     credentials.NewStaticV4(config.MinIOAccessKeyID, config.MinIOSecretAccessKey, ""),
		Secure:    config.MinIOUseSSL,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: config.MinIOSkipTLSVerify}},
	})

	if err != nil {
		return nil, err
	}

	storage := &MinIOObjectStorage{
		minIOClient: minIOClient,
		bucket:      config.MinIOBucket,
		rootPath:    config.MinIORootPath,
	}

	storage.rootPath = storage.getFinalPath(pathPrefix)
	storage.rootPath = strings.ReplaceAll(storage.rootPath, "\\", "/")

	ctx := context.Background()
	exists, err := minIOClient.BucketExists(ctx, config.MinIOBucket)

	if err != nil {
		return nil, err
	}

	if !exists {
		err := minIOClient.MakeBucket(ctx, config.MinIOBucket, minio.MakeBucketOptions{
			Region: config.MinIOLocation,
		})

		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}

// Exists returns whether the file exists
func (s *MinIOObjectStorage) Exists(ctx *webcontext.WebContext, path string) (bool, error) {
	objectInfo, err := s.minIOClient.StatObject(ctx.Request.Context(), s.bucket, s.getFinalPath(path), minio.StatObjectOptions{})

	if err == nil && !objectInfo.IsDeleteMarker {
		return true, nil
	}

	return false, err
}

// Read returns the object instance according to specified the file path
func (s *MinIOObjectStorage) Read(ctx *webcontext.WebContext, path string) (ObjectInStorage, error) {
	return s.minIOClient.GetObject(ctx.Request.Context(), s.bucket, s.getFinalPath(path), minio.GetObjectOptions{})
}

// Save returns whether save the object instance successfully
func (s *MinIOObjectStorage) Save(ctx *webcontext.WebContext, path string, object ObjectInStorage) error {
	_, err := s.minIOClient.PutObject(ctx.Request.Context(), s.bucket, s.getFinalPath(path), object, -1, minio.PutObjectOptions{})

	return err
}

// Delete returns whether delete the object according to specified the file path successfully
func (s *MinIOObjectStorage) Delete(ctx *webcontext.WebContext, path string) error {
	return s.minIOClient.RemoveObject(ctx.Request.Context(), s.bucket, s.getFinalPath(path), minio.RemoveObjectOptions{})
}

func (s *MinIOObjectStorage) getFinalPath(path string) string {
	rootPath := s.rootPath

	if len(rootPath) > 0 && rootPath[len(rootPath)-1] != '/' {
		rootPath = rootPath + "/"
	}

	if len(rootPath) > 0 && rootPath[0] == '/' {
		rootPath = rootPath[1:]
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	path = strings.ReplaceAll(path, "\\", "/")

	return rootPath + path
}