package api

import (
	"sort"

	"omnievent-backend/internal/model"
	"omnievent-backend/internal/service"
	"omnievent-backend/pkg/api"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
	"omnievent-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// TokenHandler represents token handler
type TokenHandler struct {
	tokenService *service.TokenService
	userService  *service.UserService
}

// NewTokenHandler creates a new token handler
func NewTokenHandler(tokenService *service.TokenService) *TokenHandler {
	return &TokenHandler{tokenService: tokenService}
}

// SetUserService sets the user service (for handler composition)
func (h *TokenHandler) SetUserService(userService *service.UserService) {
	h.userService = userService
}

// TokenInfoResponse represents token info response
type TokenInfoResponse struct {
	TokenID   string `json:"tokenId"`
	TokenType int    `json:"tokenType"`
	UserAgent string `json:"userAgent"`
	LastSeen  int64  `json:"lastSeen"`
	IsCurrent bool   `json:"isCurrent"`
}

// TokenGenerateAPIRequest represents API token generation request
type TokenGenerateAPIRequest struct {
	Password         string `json:"password" binding:"required"`
	ExpiredInSeconds int64  `json:"expired_in_seconds"`
}

// TokenGenerateMCPRequest represents MCP token generation request
type TokenGenerateMCPRequest struct {
	Password         string `json:"password" binding:"required"`
	ExpiredInSeconds int64  `json:"expired_in_seconds"`
}

// TokenGenerateResponse represents token generation response
type TokenGenerateResponse struct {
	Token string `json:"token,omitempty"`
}

// TokenRevokeRequest represents token revoke request
type TokenRevokeRequest struct {
	TokenID string `json:"token_id" binding:"required"`
}

// ListTokens returns available token list of current user
func (h *TokenHandler) ListTokens(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	tokens, err := h.tokenService.GetAllUnexpiredNormalAndMCPTokensByUid(c, uid)
	if err != nil {
		return nil, errs.ErrTokenRecordNotFound
	}

	tokenResps := make([]TokenInfoResponse, len(tokens))
	claims := c.GetTokenClaims()

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tokenResp := TokenInfoResponse{
			TokenID:   h.tokenService.GenerateTokenId(token),
			TokenType: token.TokenType,
			UserAgent: token.UserAgent,
			LastSeen:  token.LastSeenAt,
		}

		if token.UID == claims.UID && utils.Int64ToString(token.UserTokenID) == utils.Int64ToString(claims.UserTokenID) && token.CreatedAt == claims.ExpiresAt-2592000 {
			tokenResp.IsCurrent = true
		}

		tokenResps[i] = tokenResp
	}

	sort.Slice(tokenResps, func(i, j int) bool {
		return tokenResps[i].LastSeen > tokenResps[j].LastSeen
	})

	return gin.H{"tokens": tokenResps}, nil
}

// RefreshToken refresh current token of current user
func (h *TokenHandler) RefreshToken(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	return gin.H{"success": true}, nil
}

// RevokeToken revokes specific token of current user
func (h *TokenHandler) RevokeToken(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	var req TokenRevokeRequest
	if err := api.BindJSON(c, &req); err != nil {
		return nil, err
	}

	tokenRecord, err := h.tokenService.ParseFromTokenId(req.TokenID)
	if err != nil {
		return nil, errs.ErrInvalidTokenId
	}

	if tokenRecord.UID != uid {
		return nil, errs.ErrInvalidTokenId
	}

	err = h.tokenService.DeleteToken(c, tokenRecord)
	if err != nil {
		return nil, errs.ErrTokenRecordNotFound
	}

	return gin.H{"success": true}, nil
}

// RevokeAllTokens revokes all tokens of current user except current token
func (h *TokenHandler) RevokeAllTokens(c *context.WebContext) (interface{}, *errs.Error) {
	uid := c.GetCurrentUID()
	if uid == 0 {
		return nil, errs.ErrUnauthorizedAccess
	}

	err := h.tokenService.DeleteTokensByType(c, uid, 1) // TokenTypeNormal
	if err != nil {
		return nil, errs.ErrTokenRecordNotFound
	}

	return gin.H{"success": true}, nil
}

// RevokeCurrentToken revokes current token of current user
func (h *TokenHandler) RevokeCurrentToken(c *context.WebContext) (interface{}, *errs.Error) {
	tokenString := c.GetTokenStringFromHeader()
	if tokenString == "" {
		return nil, errs.ErrTokenIsEmpty
	}

	_, claims, _, err := h.tokenService.ParseToken(c, tokenString)
	if err != nil {
		return nil, errs.ErrCurrentInvalidToken
	}

	tokenRecord := &model.TokenRecord{
		UID:         claims.UID,
		UserTokenID: claims.UserTokenID,
		CreatedAt:   claims.IssuedAt,
	}

	tokenId := h.tokenService.GenerateTokenId(tokenRecord)
	err = h.tokenService.DeleteToken(c, tokenRecord)
	if err != nil {
		return nil, errs.ErrTokenRecordNotFound
	}

	return gin.H{"success": true, "token_id": tokenId}, nil
}

// GenerateAPIToken generates a new API token for current user
func (h *TokenHandler) GenerateAPIToken(c *context.WebContext) (interface{}, *errs.Error) {
	var req TokenGenerateAPIRequest
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

	if !h.userService.IsPasswordEqualsUserPassword(req.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	token, _, err := h.tokenService.CreateAPIToken(c, user, req.ExpiredInSeconds)
	if err != nil {
		return nil, errs.ErrTokenGenerating
	}

	return gin.H{"token": token}, nil
}

// GenerateMCPToken generates a new MCP token for current user
func (h *TokenHandler) GenerateMCPToken(c *context.WebContext) (interface{}, *errs.Error) {
	var req TokenGenerateMCPRequest
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

	if !h.userService.IsPasswordEqualsUserPassword(req.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	token, _, err := h.tokenService.CreateMCPToken(c, user, req.ExpiredInSeconds)
	if err != nil {
		return nil, errs.ErrTokenGenerating
	}

	return gin.H{"token": token}, nil
}