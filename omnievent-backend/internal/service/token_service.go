package service

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"

	"omnievent-backend/internal/model"
	"omnievent-backend/internal/repository"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
	"omnievent-backend/pkg/utils"
)

const tokenMaxExpiredAtUnixTime = int64(253402300799) // 9999-12-31 23:59:59 UTC

// OmniEventUserTokenClaims implements jwt.Claims for OmniEvent
type OmniEventUserTokenClaims struct {
	UserTokenID int64  `json:"userTokenId"`
	UID         int64  `json:"jti,string"`
	Username    string `json:"username,omitempty"`
	TokenType   int    `json:"type"`
	IssuedAt    int64  `json:"iat"`
	ExpiresAt   int64  `json:"exp"`
}

func (c *OmniEventUserTokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(c.ExpiresAt, 0)}, nil
}

func (c *OmniEventUserTokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(c.IssuedAt, 0)}, nil
}

func (c *OmniEventUserTokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{}, nil
}

func (c *OmniEventUserTokenClaims) GetIssuer() (string, error) {
	return "", nil
}

func (c *OmniEventUserTokenClaims) GetSubject() (string, error) {
	return "", nil
}

func (c *OmniEventUserTokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}

// TokenService represents user token service
type TokenService struct {
	repo *repository.UserRepository
}

// Initialize a user token service singleton instance
var tokenService *TokenService

func NewTokenService(repo *repository.UserRepository) *TokenService {
	if tokenService == nil {
		tokenService = &TokenService{
			repo: repo,
		}
	}
	return tokenService
}

// GetAllTokensByUid returns all token models of given user
func (s *TokenService) GetAllTokensByUid(c *context.WebContext, uid int64) ([]*model.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tokenRecords []*model.TokenRecord
	err := s.repo.GetDB().NewSession().Cols("uid", "user_token_id", "token_type", "user_agent", "created_at", "expires_at").Where("uid=?", uid).Find(&tokenRecords)

	return tokenRecords, err
}

// GetAllUnexpiredNormalAndMCPTokensByUid returns all available token models of given user
func (s *TokenService) GetAllUnexpiredNormalAndMCPTokensByUid(c *context.WebContext, uid int64) ([]*model.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	var tokenRecords []*model.TokenRecord
	err := s.repo.GetDB().NewSession().Cols("uid", "user_token_id", "token_type", "user_agent", "created_at", "expires_at", "last_seen_at").Where("uid=? AND (token_type=? OR token_type=? OR token_type=?) AND expires_at>?", uid, model.TokenTypeNormal, model.TokenTypeMCP, model.TokenTypeAPI, now).Find(&tokenRecords)

	return tokenRecords, err
}

// ParseToken returns the token model according to token content
func (s *TokenService) ParseToken(c *context.WebContext, tokenString string) (*jwt.Token, *OmniEventUserTokenClaims, string, error) {
	return s.parseToken(c, tokenString)
}

// CreateToken generates a new normal token and saves to database
func (s *TokenService) CreateToken(c *context.WebContext, user *model.User) (string, *OmniEventUserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, model.TokenTypeNormal, s.getUserAgent(c), "", getTokenExpiredTimeDuration())
	return token, claims, err
}

// CreateRequire2FAToken generates a new token requiring user to verify 2fa passcode and saves to database
func (s *TokenService) CreateRequire2FAToken(c *context.WebContext, user *model.User) (string, *OmniEventUserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, model.TokenTypeRequire2FA, s.getUserAgent(c), "", getTemporaryTokenExpiredTimeDuration())
	return token, claims, err
}

// CreateEmailVerifyToken generates a new email verify token and saves to database
func (s *TokenService) CreateEmailVerifyToken(c *context.WebContext, user *model.User) (string, *OmniEventUserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, model.TokenTypeEmailVerify, s.getUserAgent(c), "", getEmailVerifyTokenExpiredTimeDuration())
	return token, claims, err
}

// CreateEmailVerifyTokenWithoutUserAgent generates a new email verify token without user agent
func (s *TokenService) CreateEmailVerifyTokenWithoutUserAgent(c *context.WebContext, user *model.User) (string, *OmniEventUserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, model.TokenTypeEmailVerify, "", "", getEmailVerifyTokenExpiredTimeDuration())
	return token, claims, err
}

// CreatePasswordResetToken generates a new password reset token and saves to database
func (s *TokenService) CreatePasswordResetToken(c *context.WebContext, user *model.User) (string, *OmniEventUserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, model.TokenTypePasswordReset, s.getUserAgent(c), "", getPasswordResetTokenExpiredTimeDuration())
	return token, claims, err
}

// CreatePasswordResetTokenWithoutUserAgent generates a new password reset token without user agent
func (s *TokenService) CreatePasswordResetTokenWithoutUserAgent(c *context.WebContext, user *model.User) (string, *OmniEventUserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, model.TokenTypePasswordReset, "", "", getPasswordResetTokenExpiredTimeDuration())
	return token, claims, err
}

// CreateAPIToken generates a new API token and saves to database
func (s *TokenService) CreateAPIToken(c *context.WebContext, user *model.User, expiresInSeconds int64) (string, *OmniEventUserTokenClaims, error) {
	var tokenExpiredTimeDuration time.Duration

	if expiresInSeconds > 0 {
		tokenExpiredTimeDuration = time.Duration(expiresInSeconds) * time.Second
	} else {
		tokenExpiredTimeDuration = time.Unix(tokenMaxExpiredAtUnixTime, 0).Sub(time.Now())
	}

	token, claims, _, err := s.createToken(c, user, model.TokenTypeAPI, s.getUserAgent(c), "", tokenExpiredTimeDuration)
	return token, claims, err
}

// CreateMCPToken generates a new MCP token and saves to database
func (s *TokenService) CreateMCPToken(c *context.WebContext, user *model.User, expiresInSeconds int64) (string, *OmniEventUserTokenClaims, error) {
	var tokenExpiredTimeDuration time.Duration

	if expiresInSeconds > 0 {
		tokenExpiredTimeDuration = time.Duration(expiresInSeconds) * time.Second
	} else {
		tokenExpiredTimeDuration = time.Unix(tokenMaxExpiredAtUnixTime, 0).Sub(time.Now())
	}

	token, claims, _, err := s.createToken(c, user, model.TokenTypeMCP, s.getUserAgent(c), "", tokenExpiredTimeDuration)
	return token, claims, err
}

// CreateOAuth2CallbackRequireVerifyToken generates a new OAuth 2.0 callback token requiring user to verify and saves to database
func (s *TokenService) CreateOAuth2CallbackRequireVerifyToken(c *context.WebContext, user *model.User, ctx string) (string, *OmniEventUserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, model.TokenTypeOAuth2CallbackRequireVerify, s.getUserAgent(c), ctx, getTemporaryTokenExpiredTimeDuration())
	return token, claims, err
}

// CreateOAuth2CallbackToken generates a new OAuth 2.0 callback token and saves to database
func (s *TokenService) CreateOAuth2CallbackToken(c *context.WebContext, user *model.User, ctx string) (string, *OmniEventUserTokenClaims, error) {
	token, claims, _, err := s.createToken(c, user, model.TokenTypeOAuth2Callback, s.getUserAgent(c), ctx, getTemporaryTokenExpiredTimeDuration())
	return token, claims, err
}

// UpdateTokenLastSeen updates the last seen time of specified token
func (s *TokenService) UpdateTokenLastSeen(c *context.WebContext, tokenRecord *model.TokenRecord) error {
	if tokenRecord.UID <= 0 {
		return errs.ErrUserIdInvalid
	}

	if tokenRecord.UserTokenID <= 0 {
		return errs.ErrInvalidUserTokenId
	}

	tokenRecord.LastSeenAt = time.Now().Unix()

	sess := s.repo.GetDB().NewSession()
	defer sess.Close()

	updatedRows, err := sess.Cols("last_seen_at").Where("uid=? AND user_token_id=? AND created_at=?", tokenRecord.UID, tokenRecord.UserTokenID, tokenRecord.CreatedAt).Update(tokenRecord)
	if err != nil {
		return err
	}
	if updatedRows < 1 {
		return errs.ErrTokenRecordNotFound
	}

	return nil
}

// DeleteToken deletes given token from database
func (s *TokenService) DeleteToken(c *context.WebContext, tokenRecord *model.TokenRecord) error {
	if tokenRecord.UID <= 0 {
		return errs.ErrUserIdInvalid
	}

	if tokenRecord.UserTokenID <= 0 {
		return errs.ErrInvalidUserTokenId
	}

	sess := s.repo.GetDB().NewSession()
	defer sess.Close()

	deletedRows, err := sess.Where("uid=? AND user_token_id=? AND created_at=?", tokenRecord.UID, tokenRecord.UserTokenID, tokenRecord.CreatedAt).Delete(&model.TokenRecord{})
	if err != nil {
		return err
	}
	if deletedRows < 1 {
		return errs.ErrTokenRecordNotFound
	}

	return nil
}

// DeleteTokens deletes given tokens from database
func (s *TokenService) DeleteTokens(c *context.WebContext, uid int64, tokenRecords []*model.TokenRecord) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	sess := s.repo.GetDB().NewSession()
	defer sess.Close()

	for i := 0; i < len(tokenRecords); i++ {
		tokenRecord := tokenRecords[i]
		deletedRows, err := sess.Where("uid=? AND user_token_id=? AND created_at=?", uid, tokenRecord.UserTokenID, tokenRecord.CreatedAt).Delete(&model.TokenRecord{})
		if err != nil {
			return err
		}
		if deletedRows < 1 {
			return errs.ErrTokenRecordNotFound
		}
	}

	return nil
}

// DeleteTokenByClaims deletes given token from database
func (s *TokenService) DeleteTokenByClaims(c *context.WebContext, claims *OmniEventUserTokenClaims) error {
	return s.DeleteToken(c, &model.TokenRecord{
		UID:         claims.UID,
		UserTokenID: claims.UserTokenID,
		CreatedAt:   claims.IssuedAt,
	})
}

// DeleteTokensBeforeTime deletes tokens that is created before specific time
func (s *TokenService) DeleteTokensBeforeTime(c *context.WebContext, uid int64, createTime int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	_, err := s.repo.GetDB().Where("uid=? AND created_at<?", uid, createTime).Delete(&model.TokenRecord{})
	return err
}

// DeleteTokensByType deletes specified type tokens
func (s *TokenService) DeleteTokensByType(c *context.WebContext, uid int64, tokenType int) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	_, err := s.repo.GetDB().Where("uid=? AND token_type=?", uid, tokenType).Delete(&model.TokenRecord{})
	return err
}

// DeleteAllExpiredTokens deletes all expired tokens
func (s *TokenService) DeleteAllExpiredTokens(c *context.WebContext) error {
	count, err := s.repo.GetDB().Where("expires_at<=?", time.Now().Unix()).Delete(&model.TokenRecord{})
	if count > 0 {
		// log not available in this context
	}
	return err
}

// ExistsValidTokenByType returns whether the given token type exists
func (s *TokenService) ExistsValidTokenByType(c *context.WebContext, uid int64, tokenType int) (bool, error) {
	if uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	return s.repo.GetDB().NewSession().Cols("uid", "user_token_id", "expires_at").Where("uid=? AND token_type=? AND expires_at>?", uid, tokenType, now).Exist(&model.TokenRecord{})
}

// ParseFromTokenId returns token model according to token id
func (s *TokenService) ParseFromTokenId(tokenId string) (*model.TokenRecord, error) {
	pairs := strings.Split(tokenId, ":")

	if len(pairs) != 3 {
		return nil, errs.ErrInvalidTokenId
	}

	uid := utils.StringToInt64(pairs[0])
	createdAt := utils.StringToInt64(pairs[1])
	userTokenId := utils.StringToInt64(pairs[2])

	tokenRecord := &model.TokenRecord{
		UID:         uid,
		UserTokenID: userTokenId,
		CreatedAt:   createdAt,
	}

	return tokenRecord, nil
}

// GenerateTokenId generates token id according to token model
func (s *TokenService) GenerateTokenId(tokenRecord *model.TokenRecord) string {
	return fmt.Sprintf("%d:%d:%d", tokenRecord.UID, tokenRecord.CreatedAt, tokenRecord.UserTokenID)
}

func (s *TokenService) parseToken(c *context.WebContext, tokenString string) (*jwt.Token, *OmniEventUserTokenClaims, string, error) {
	claims := &OmniEventUserTokenClaims{}
	tokenContext := ""

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (any, error) {
			now := time.Now().Unix()
			userTokenId := claims.UserTokenID

			if userTokenId <= 0 {
				return nil, errs.ErrInvalidUserTokenId
			}

			tokenRecord, err := s.getTokenRecord(c, claims.UID, userTokenId, claims.IssuedAt)

			if err != nil {
				return nil, errs.ErrTokenRecordNotFound
			}

			if tokenRecord.ExpiresAt < now {
				return nil, errs.ErrTokenExpired
			}

			tokenContext = string(tokenRecord.Context)
			return []byte(tokenRecord.Secret), nil
		},
		jwt.WithIssuedAt(),
	)

	if err != nil {
		if errors.Is(err, request.ErrNoTokenInRequest) {
			return nil, nil, "", errs.ErrTokenIsEmpty
		}

		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenUnverifiable) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, nil, "", errs.ErrCurrentInvalidToken
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, nil, "", errs.ErrCurrentTokenExpired
		}

		if errors.Is(err, jwt.ErrTokenUsedBeforeIssued) {
			return nil, nil, "", errs.ErrCurrentInvalidToken
		}

		return nil, nil, "", err
	}

	return token, claims, tokenContext, err
}

func (s *TokenService) createToken(c *context.WebContext, user *model.User, tokenType int, userAgent string, tokenContext string, expiryTime time.Duration) (string, *OmniEventUserTokenClaims, *model.TokenRecord, error) {
	now := time.Now()

	tokenRecord := &model.TokenRecord{
		UID:         user.UID,
		UserTokenID: s.getUserTokenId(),
		TokenType:   tokenType,
		UserAgent:   userAgent,
		Context:     []byte(tokenContext),
		CreatedAt:   now.Unix(),
		ExpiresAt:   now.Add(expiryTime).Unix(),
		LastSeenAt:  now.Unix(),
	}

	tokenRecord.Secret = utils.GenerateRandomString(10)

	claims := &OmniEventUserTokenClaims{
		UserTokenID: tokenRecord.UserTokenID,
		UID:         tokenRecord.UID,
		Username:    user.Username,
		TokenType:   tokenRecord.TokenType,
		IssuedAt:    tokenRecord.CreatedAt,
		ExpiresAt:   tokenRecord.ExpiresAt,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(tokenRecord.Secret))
	if err != nil {
		return "", nil, nil, err
	}

	err = s.createTokenRecord(c, tokenRecord)
	if err != nil {
		return "", nil, nil, err
	}

	return tokenString, claims, tokenRecord, err
}

func (s *TokenService) getTokenRecord(c *context.WebContext, uid int64, userTokenId int64, createAt int64) (*model.TokenRecord, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if userTokenId <= 0 {
		return nil, errs.ErrInvalidUserTokenId
	}

	tokenRecord := &model.TokenRecord{}
	has, err := s.repo.GetDB().NewSession().Where("uid=? AND user_token_id=? AND created_at=?", uid, userTokenId, createAt).Limit(1).Get(tokenRecord)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errs.ErrTokenRecordNotFound
	}

	return tokenRecord, nil
}

func (s *TokenService) createTokenRecord(c *context.WebContext, tokenRecord *model.TokenRecord) error {
	if tokenRecord.UID <= 0 {
		return errs.ErrUserIdInvalid
	}

	if tokenRecord.UserTokenID <= 0 {
		return errs.ErrInvalidUserTokenId
	}

	sess := s.repo.GetDB().NewSession()
	defer sess.Close()

	_, err := sess.Insert(tokenRecord)
	return err
}

func (s *TokenService) getUserTokenId() int64 {
	nanoSeconds := time.Now().Nanosecond()
	randomNumber, _ := utils.GetRandomInteger(math.MaxInt32)
	userTokenId := (int64(nanoSeconds) << 32) | int64(randomNumber)

	return userTokenId
}

func (s *TokenService) getUserAgent(c *context.WebContext) string {
	userAgent := ""

	if c != nil && c.Request != nil {
		userAgent = c.GetUserAgent()
	}

	if len(userAgent) > 255 {
		userAgent = utils.SubString(userAgent, 0, 255)
	}

	return userAgent
}

// Token durations - these would typically come from config
func getTokenExpiredTimeDuration() time.Duration {
	return 24 * time.Hour
}

func getTemporaryTokenExpiredTimeDuration() time.Duration {
	return 15 * time.Minute
}

func getEmailVerifyTokenExpiredTimeDuration() time.Duration {
	return 24 * time.Hour
}

func getPasswordResetTokenExpiredTimeDuration() time.Duration {
	return 1 * time.Hour
}
