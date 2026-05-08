package api

import (
	"github.com/gin-gonic/gin"

	"omnievent-backend/internal/model"
	"omnievent-backend/internal/service"
	"omnievent-backend/pkg/api"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
)

type AuthHandler struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

func NewAuthHandler(userService *service.UserService, tokenService *service.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required,min=1,max=64"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *context.WebContext) (interface{}, *errs.Error) {
	var req RegisterRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	user, err := h.userService.CreateUser(req.Username, req.Email, req.Nickname, req.Password)
	if err != nil {
		return nil, errs.NewNormalError(errs.SubCategoryUser, 2, 409, "user already exists")
	}

	return gin.H{
		"uid":      user.UID,
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.Nickname,
	}, nil
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string           `json:"token"`
	Need2FA  bool            `json:"need_2fa"`
	User     *ProfileResponse `json:"user,omitempty"`
}

func (h *AuthHandler) Login(c *context.WebContext) (interface{}, *errs.Error) {
	var req LoginRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	user, err := h.userService.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	token, _, err := h.tokenService.CreateToken(user.UID, model.TokenTypeNormal)
	if err != nil {
		return nil, errs.NewNormalError(errs.SubCategoryToken, 1, 500, "failed to create token")
	}

	h.userService.UpdateLastLogin(user.UID)

	return LoginResponse{
		Token:   token,
		Need2FA: false,
		User: &ProfileResponse{
			UID:      user.UID,
			Username: user.Username,
			Email:    user.Email,
			Nickname: user.Nickname,
		},
	}, nil
}

type ForgetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgetPassword(c *context.WebContext) (interface{}, *errs.Error) {
	var req ForgetPasswordRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	_, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		// 即使邮箱不存在也返回成功，防止泄露用户信息
		return gin.H{
			"success": true,
			"message": "if the email exists, a reset link has been sent",
		}, nil
	}

	return gin.H{
		"success": true,
		"message": "if the email exists, a reset link has been sent",
	}, nil
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ResetPassword(c *context.WebContext) (interface{}, *errs.Error) {
	var req ResetPasswordRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	return gin.H{
		"success": true,
	}, nil
}
