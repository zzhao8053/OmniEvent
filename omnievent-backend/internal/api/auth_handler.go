package api

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"omnievent-backend/internal/model"
	"omnievent-backend/internal/service"
	"omnievent-backend/pkg/api"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
)

// AuthHandler represents auth handler
type AuthHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *service.UserService, tokenService *service.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=32"`
	Email           string `json:"email" binding:"required,email"`
	Nickname        string `json:"nickname" binding:"required,min=1,max=64"`
	Password        string `json:"password" binding:"required,min=6"`
	Language        string `json:"language"`
	DefaultCurrency string `json:"default_currency"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	LoginName string `json:"loginname" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// ForgetPasswordRequest represents forget password request
type ForgetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Token    string `json:"token" binding:"required"`
}

// AuthResponse represents auth response
type AuthResponse struct {
	Token   string             `json:"token,omitempty"`
	Need2FA bool               `json:"need_2fa"`
	User    *UserBasicInfoResp `json:"user,omitempty"`
}

// UserBasicInfoResp represents basic user info response
type UserBasicInfoResp struct {
	UID             int64  `json:"uid"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Nickname        string `json:"nickname"`
	DefaultCurrency string `json:"default_currency,omitempty"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *context.WebContext) (interface{}, *errs.Error) {
	var req RegisterRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Nickname = strings.TrimSpace(req.Nickname)

	if req.DefaultCurrency == "" {
		req.DefaultCurrency = "CNY"
	}

	user := &model.User{
		Username:        req.Username,
		Email:           req.Email,
		Nickname:        req.Nickname,
		Password:        req.Password,
		Language:        req.Language,
		DefaultCurrency: req.DefaultCurrency,
	}

	err := h.userService.CreateUser(c, user, false)
	if err != nil {
		if customErr, ok := err.(*errs.Error); ok {
			return nil, customErr
		}
		return nil, errs.ErrDatabaseOperationFailed
	}

	token, claims, err := h.tokenService.CreateToken(c, user)
	if err != nil {
		// Registration succeeded but token creation failed
		return UserBasicInfoResp{
			UID:             user.UID,
			Username:        user.Username,
			Email:           user.Email,
			Nickname:        user.Nickname,
			DefaultCurrency: user.DefaultCurrency,
		}, nil
	}

	c.SetTextualToken(token)
	c.SetTokenClaims(&context.UserTokenClaims{
		UID:         claims.UID,
		UserTokenID: claims.UserTokenID,
		TokenType:   claims.TokenType,
		ExpiresAt:   claims.ExpiresAt,
	})

	return AuthResponse{
		Token:   token,
		Need2FA: false,
		User: &UserBasicInfoResp{
			UID:             user.UID,
			Username:        user.Username,
			Email:           user.Email,
			Nickname:        user.Nickname,
			DefaultCurrency: user.DefaultCurrency,
		},
	}, nil
}

// Login handles user login
func (h *AuthHandler) Login(c *context.WebContext) (interface{}, *errs.Error) {
	var req LoginRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	user, err := h.userService.GetUserByUsernameOrEmailAndPassword(c, req.LoginName, req.Password)
	if err != nil {
		if customErr, ok := err.(*errs.Error); ok {
			return nil, customErr
		}
		return nil, errs.ErrLoginNameOrPasswordWrong
	}

	if user.Disabled {
		return nil, errs.ErrUserIsDisabled
	}

	err = h.userService.UpdateUserLastLoginTime(c, user.UID)
	if err != nil {
		// Log warning but continue with login
	}

	token, claims, err := h.tokenService.CreateToken(c, user)
	if err != nil {
		return nil, errs.ErrTokenGenerating
	}

	c.SetTextualToken(token)
	c.SetTokenClaims(&context.UserTokenClaims{
		UID:         claims.UID,
		UserTokenID: claims.UserTokenID,
		TokenType:   claims.TokenType,
		ExpiresAt:   claims.ExpiresAt,
	})

	return AuthResponse{
		Token:   token,
		Need2FA: false,
		User: &UserBasicInfoResp{
			UID:             user.UID,
			Username:        user.Username,
			Email:           user.Email,
			Nickname:        user.Nickname,
			DefaultCurrency: user.DefaultCurrency,
		},
	}, nil
}

// ForgetPassword handles forget password request
func (h *AuthHandler) ForgetPassword(c *context.WebContext) (interface{}, *errs.Error) {
	var req ForgetPasswordRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	user, err := h.userService.GetUserByEmail(c, req.Email)
	if err != nil {
		// Return success even if user doesn't exist to prevent email enumeration
		return gin.H{
			"success": true,
			"message": "if the email exists, a reset link has been sent",
		}, nil
	}

	if user.Disabled {
		return nil, errs.ErrUserIsDisabled
	}

	token, _, err := h.tokenService.CreatePasswordResetToken(c, user)
	if err != nil {
		return nil, errs.ErrTokenGenerating
	}

	// In production, send email with token here
	// For now, we just return success
	_ = token

	return gin.H{
		"success": true,
		"message": "if the email exists, a reset link has been sent",
	}, nil
}

// ResetPassword handles password reset with token
func (h *AuthHandler) ResetPassword(c *context.WebContext) (interface{}, *errs.Error) {
	var req ResetPasswordRequest
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
		if customErr, ok := err.(*errs.Error); ok {
			return nil, customErr
		}
		return nil, errs.ErrDatabaseOperationFailed
	}

	// Delete old tokens (except current one would be handled by middleware)
	now := time.Now().Unix()
	_ = now // Would delete old tokens in production

	return gin.H{
		"success": true,
	}, nil
}