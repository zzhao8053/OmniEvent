package api

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"omnievent-backend/internal/service"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
	"omnievent-backend/pkg/log"
	"omnievent-backend/pkg/storage"
	"omnievent-backend/pkg/utils"
)

// AvatarHandler represents avatar handler
type AvatarHandler struct {
	userService *service.UserService
}

// NewAvatarHandler creates a new avatar handler
func NewAvatarHandler(userService *service.UserService) *AvatarHandler {
	return &AvatarHandler{
		userService: userService,
	}
}

// Allowed avatar MIME types
var allowedAvatarTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// MaxAvatarFileSize 5MB
const MaxAvatarFileSize = 5 * 1024 * 1024

// UploadAvatar handles avatar file upload via multipart form
func (h *AvatarHandler) UploadAvatar(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	user, err := h.userService.GetUserById(c, uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	// Get multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max memory
		log.Errorf(c, "[avatar.UploadAvatar] failed to parse multipart form for user %d: %s", uid, err.Error())
		return nil, errs.ErrParameterInvalid
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		log.Warnf(c, "[avatar.UploadAvatar] no avatar file in request for user %d", uid)
		return nil, errs.ErrParameterInvalid
	}
	defer file.Close()

	// Validate file size
	if header.Size < 1 {
		log.Warnf(c, "[avatar.UploadAvatar] avatar file is empty for user %d", uid)
		return nil, errs.ErrParameterInvalid
	}

	if header.Size > MaxAvatarFileSize {
		log.Warnf(c, "[avatar.UploadAvatar] avatar file size %d exceeds max %d for user %d", header.Size, MaxAvatarFileSize, uid)
		return nil, errs.ErrParameterInvalid
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isValidAvatarExtension(ext) {
		log.Warnf(c, "[avatar.UploadAvatar] invalid avatar extension '%s' for user %d", ext, uid)
		return nil, errs.ErrImageTypeNotSupported
	}

	// Validate MIME type
	contentType := header.Header.Get("Content-Type")
	if !allowedAvatarTypes[contentType] {
		// Check by extension as fallback
		if !isValidAvatarExtension(ext) {
			log.Warnf(c, "[avatar.UploadAvatar] unsupported avatar type '%s' for user %d", contentType, uid)
			return nil, errs.ErrImageTypeNotSupported
		}
	}

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		log.Errorf(c, "[avatar.UploadAvatar] failed to read avatar file for user %d: %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	// Delete old avatar if exists
	if user.CustomAvatarType != "" {
		oldPath := h.getAvatarPath(user.UID, user.CustomAvatarType)
		if err := storage.Container.DeleteAvatar(c, oldPath); err != nil {
			log.Warnf(c, "[avatar.UploadAvatar] failed to delete old avatar at %s: %s", oldPath, err.Error())
			// Continue anyway, non-critical error
		}
	}

	// Save new avatar
	newPath := h.getAvatarPath(uid, ext)
	reader := &bytesReader{data: data}
	if err := storage.Container.SaveAvatar(c, newPath, reader); err != nil {
		log.Errorf(c, "[avatar.UploadAvatar] failed to save avatar for user %d: %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	// Update user avatar type
	user.CustomAvatarType = ext
	if _, err := h.userService.UpdateUser(c, user); err != nil {
		log.Errorf(c, "[avatar.UploadAvatar] failed to update user avatar type for user %d: %s", uid, err.Error())
		return nil, errs.ErrDatabaseOperationFailed
	}

	avatarURL := h.getAvatarURL(uid, ext)
	log.Infof(c, "[avatar.UploadAvatar] avatar uploaded successfully for user %d", uid)

	return map[string]interface{}{
		"avatar_url": avatarURL,
	}, nil
}

// GetAvatar returns avatar file data
func (h *AvatarHandler) GetAvatar(c *context.WebContext) ([]byte, string, *errs.Error) {
	fileName := c.Param("fileName")
	if fileName == "" {
		return nil, "", errs.ErrParameterInvalid
	}

	ext := filepath.Ext(fileName)
	if !isValidAvatarExtension(ext) {
		return nil, "", errs.ErrImageTypeNotSupported
	}

	// Validate user can only access their own avatar
	uid := c.GetCurrentUID()
	fileBaseName := strings.TrimSuffix(fileName, ext)
	if !isNumericString(fileBaseName) || utils.Int64ToString(uid) != fileBaseName {
		log.Warnf(c, "[avatar.GetAvatar] user %d attempted to access avatar of user %s", uid, fileBaseName)
		return nil, "", errs.ErrUserIdInvalid
	}

	path := h.getAvatarPath(uid, ext)

	obj, err := storage.Container.ReadAvatar(c, path)
	if err != nil {
		log.Errorf(c, "[avatar.GetAvatar] failed to read avatar at %s: %s", path, err.Error())
		return nil, "", errs.ErrOperationFailed
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		log.Errorf(c, "[avatar.GetAvatar] failed to read avatar data at %s: %s", path, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	contentType := getAvatarContentType(ext)
	return data, contentType, nil
}

// RemoveAvatar removes user's avatar
func (h *AvatarHandler) RemoveAvatar(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	user, err := h.userService.GetUserById(c, uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	if user.CustomAvatarType == "" {
		return nil, errs.ErrNothingWillBeUpdated
	}

	// Delete avatar file
	path := h.getAvatarPath(uid, user.CustomAvatarType)
	if err := storage.Container.DeleteAvatar(c, path); err != nil {
		log.Errorf(c, "[avatar.RemoveAvatar] failed to delete avatar at %s: %s", path, err.Error())
		return nil, errs.ErrOperationFailed
	}

	// Update user avatar type
	user.CustomAvatarType = ""
	if _, err := h.userService.UpdateUser(c, user); err != nil {
		log.Errorf(c, "[avatar.RemoveAvatar] failed to update user avatar type for user %d: %s", uid, err.Error())
		return nil, errs.ErrDatabaseOperationFailed
	}

	log.Infof(c, "[avatar.RemoveAvatar] avatar removed successfully for user %d", uid)

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (h *AvatarHandler) getAvatarPath(uid int64, ext string) string {
	return fmt.Sprintf("avatars/%d%s", uid, ext)
}

func (h *AvatarHandler) getAvatarURL(uid int64, ext string) string {
	return fmt.Sprintf("/api/v1/avatars/%d%s", uid, ext)
}

func isValidAvatarExtension(ext string) bool {
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validExts[ext]
}

func getAvatarContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
	}
	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

func isNumericString(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// bytesReader wraps a byte slice as an ObjectInStorage
type bytesReader struct {
	data []byte
	pos  int
}

func (b *bytesReader) Read(p []byte) (n int, err error) {
	if b.pos >= len(b.data) {
		return 0, io.EOF
	}
	n = copy(p, b.data[b.pos:])
	b.pos += n
	return n, nil
}

func (b *bytesReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0: // SeekStart
		b.pos = int(offset)
	case 1: // SeekCurrent
		b.pos += int(offset)
	case 2: // SeekEnd
		b.pos = len(b.data) + int(offset)
	}
	if b.pos < 0 {
		b.pos = 0
	}
	if b.pos > len(b.data) {
		b.pos = len(b.data)
	}
	return int64(b.pos), nil
}

func (b *bytesReader) Close() error {
	return nil
}