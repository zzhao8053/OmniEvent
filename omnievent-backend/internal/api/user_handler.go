package api

import (
	"strings"

	"omnievent-backend/internal/service"
	"omnievent-backend/pkg/api"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"

	"github.com/gin-gonic/gin"
)

// UserHandler represents user handler
type UserHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, tokenService *service.TokenService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

// ProfileResponse represents user profile response
type ProfileResponse struct {
	UID             int64  `json:"uid"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Nickname        string `json:"nickname"`
	DefaultCurrency string `json:"default_currency"`
	AvatarURL       string `json:"avatar_url,omitempty"`
}

// GetProfile returns user profile of current user
func (h *UserHandler) GetProfile(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	user, err := h.userService.GetUserById(c, uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	return ProfileResponse{
		UID:             user.UID,
		Username:        user.Username,
		Email:           user.Email,
		Nickname:        user.Nickname,
		DefaultCurrency: user.DefaultCurrency,
		AvatarURL:       user.AvatarURL,
	}, nil
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Nickname        string `json:"nickname"`
	Email           string `json:"email"`
	DefaultCurrency string `json:"default_currency"`
}

// UpdateProfile saves user profile by request parameters for current user
func (h *UserHandler) UpdateProfile(c *context.WebContext) (interface{}, *errs.Error) {
	var req UpdateProfileRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	user, err := h.userService.GetUserById(c, uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Nickname = strings.TrimSpace(req.Nickname)

	if req.Email != "" && req.Email != user.Email {
		user.Email = req.Email
	}

	if req.Nickname != "" && req.Nickname != user.Nickname {
		user.Nickname = req.Nickname
	}

	if req.DefaultCurrency != "" && req.DefaultCurrency != user.DefaultCurrency {
		user.DefaultCurrency = req.DefaultCurrency
	}

	_, err = h.userService.UpdateUser(c, user)
	if err != nil {
		return nil, errs.ErrDatabaseOperationFailed
	}

	return ProfileResponse{
		UID:             user.UID,
		Username:        user.Username,
		Email:           user.Email,
		Nickname:        user.Nickname,
		DefaultCurrency: user.DefaultCurrency,
		AvatarURL:       user.AvatarURL,
	}, nil
}

// UpdateAvatarRequest represents avatar update request
type UpdateAvatarRequest struct {
	AvatarType string `json:"avatar_type" binding:"required"`
}

// UpdateAvatar saves user avatar by request parameters for current user
func (h *UserHandler) UpdateAvatar(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	var req UpdateAvatarRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	user, err := h.userService.GetUserById(c, uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	user.CustomAvatarType = req.AvatarType

	_, err = h.userService.UpdateUser(c, user)
	if err != nil {
		return nil, errs.ErrDatabaseOperationFailed
	}

	return gin.H{
		"avatar_url": user.AvatarURL,
	}, nil
}

// RemoveAvatar removes user avatar by request parameters for current user
func (h *UserHandler) RemoveAvatar(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	user, err := h.userService.GetUserById(c, uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	if user.CustomAvatarType == "" {
		return gin.H{"success": false, "message": "nothing to remove"}, nil
	}

	user.CustomAvatarType = ""

	_, err = h.userService.UpdateUser(c, user)
	if err != nil {
		return nil, errs.ErrDatabaseOperationFailed
	}

	return gin.H{
		"success": true,
	}, nil
}

// ResendVerifyEmail sends verification email for unlogin user
func (h *UserHandler) ResendVerifyEmail(c *context.WebContext) (interface{}, *errs.Error) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	user, err := h.userService.GetUserByEmail(c, req.Email)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	if !h.userService.IsPasswordEqualsUserPassword(req.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	if user.Disabled {
		return nil, errs.ErrUserIsDisabled
	}

	token, _, err := h.tokenService.CreateEmailVerifyToken(c, user)
	if err != nil {
		return nil, errs.ErrTokenGenerating
	}

	_ = token // In production, send email with token here

	return gin.H{"success": true}, nil
}

// SendVerifyEmail sends verification email for logged in user
func (h *UserHandler) SendVerifyEmail(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	user, err := h.userService.GetUserById(c, uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	token, _, err := h.tokenService.CreateEmailVerifyToken(c, user)
	if err != nil {
		return nil, errs.ErrTokenGenerating
	}

	_ = token // In production, send email with token here

	return gin.H{"success": true}, nil
}