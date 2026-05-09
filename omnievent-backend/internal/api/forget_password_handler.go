package api

import (
	"strings"

	"omnievent-backend/internal/service"
	"omnievent-backend/pkg/api"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"

	"github.com/gin-gonic/gin"
)

// ForgetPasswordHandler represents forget password handler
type ForgetPasswordHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

// NewForgetPasswordHandler creates a new forget password handler
func NewForgetPasswordHandler(userService *service.UserService, tokenService *service.TokenService) *ForgetPasswordHandler {
	return &ForgetPasswordHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

// UserForgetPasswordRequest handles password reset request (sends email)
func (h *ForgetPasswordHandler) UserForgetPasswordRequest(c *context.WebContext) (interface{}, *errs.Error) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	req.Email = strings.TrimSpace(req.Email)

	user, err := h.userService.GetUserByEmail(c, req.Email)
	if err != nil {
		// Return success even if user doesn't exist to prevent email enumeration
		return gin.H{"success": true}, nil
	}

	if user.Disabled {
		return nil, errs.ErrUserIsDisabled
	}

	token, _, err := h.tokenService.CreatePasswordResetToken(c, user)
	if err != nil {
		return nil, errs.ErrTokenGenerating
	}

	// In production, send email with token here
	_ = token

	return gin.H{"success": true}, nil
}

// UserResetPassword handles password reset with token
func (h *ForgetPasswordHandler) UserResetPassword(c *context.WebContext) (interface{}, *errs.Error) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Token    string `json:"token" binding:"required"`
	}
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

	if user.Disabled {
		return nil, errs.ErrUserIsDisabled
	}

	if user.Email != req.Email {
		return nil, errs.ErrEmailIsInvalid
	}

	// Check if new password is same as old
	if h.userService.IsPasswordEqualsUserPassword(req.Password, user) {
		return nil, errs.ErrNewPasswordEqualsOldInvalid
	}

	// Update password
	user.Password = req.Password
	err = h.userService.UpdateUserPassword(c, user)
	if err != nil {
		return nil, errs.ErrDatabaseOperationFailed
	}

	return gin.H{"success": true}, nil
}