package api

import (
	"omnievent-backend/internal/service"
	"omnievent-backend/pkg/api"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type ProfileResponse struct {
	UID       int64  `json:"uid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

func (h *UserHandler) GetProfile(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrInvalidCredentials
	}

	user, err := h.userService.GetUserByUID(uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	return ProfileResponse{
		UID:       user.UID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}, nil
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (h *UserHandler) UpdateProfile(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrInvalidCredentials
	}

	var req UpdateProfileRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	user, err := h.userService.UpdateProfile(uid, req.Nickname, req.Email)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	return ProfileResponse{
		UID:       user.UID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}, nil
}

type UpdateAvatarRequest struct {
	AvatarType string `json:"avatar_type" binding:"required"`
}

func (h *UserHandler) UpdateAvatar(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrInvalidCredentials
	}

	var req UpdateAvatarRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	user, err := h.userService.UpdateAvatar(uid, req.AvatarType)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	return gin.H{
		"avatar_url": user.AvatarURL,
	}, nil
}

func (h *UserHandler) RemoveAvatar(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrInvalidCredentials
	}

	_, err := h.userService.RemoveAvatar(uid)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	return gin.H{
		"success": true,
	}, nil
}
