package api

import (
	"omnievent-backend/internal/service"
	"omnievent-backend/pkg/api"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	tokenService *service.TokenService
}

func NewTokenHandler(tokenService *service.TokenService) *TokenHandler {
	return &TokenHandler{tokenService: tokenService}
}

func (h *TokenHandler) ListTokens(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrInvalidCredentials
	}

	tokens, err := h.tokenService.GetTokensByUID(uid)
	if err != nil {
		return nil, errs.ErrTokenNotFound
	}

	return gin.H{
		"tokens": tokens,
	}, nil
}

type RefreshTokenRequest struct {
	UserTokenID int64 `json:"user_token_id" binding:"required"`
}

func (h *TokenHandler) RefreshToken(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrInvalidCredentials
	}

	var req RefreshTokenRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	token, _, err := h.tokenService.CreateToken(uid, 1)
	if err != nil {
		return nil, errs.ErrTokenInvalid
	}

	return gin.H{
		"token": token,
	}, nil
}

type RevokeTokenRequest struct {
	UserTokenID int64 `json:"user_token_id" binding:"required"`
}

func (h *TokenHandler) RevokeToken(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrInvalidCredentials
	}

	var req RevokeTokenRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	err := h.tokenService.RevokeToken(uid, req.UserTokenID)
	if err != nil {
		return nil, errs.ErrTokenNotFound
	}

	return gin.H{
		"success": true,
	}, nil
}

func (h *TokenHandler) RevokeAllTokens(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrInvalidCredentials
	}

	err := h.tokenService.RevokeAllTokens(uid)
	if err != nil {
		return nil, errs.ErrTokenNotFound
	}

	return gin.H{
		"success": true,
	}, nil
}
