package service

import (
	"time"

	"omnievent-backend/internal/model"
	"omnievent-backend/internal/repository"
)

type TokenService struct {
	repo *repository.UserRepository
}

func NewTokenService(repo *repository.UserRepository) *TokenService {
	return &TokenService{repo: repo}
}

func (s *TokenService) CreateToken(uid int64, tokenType int) (string, *model.TokenRecord, error) {
	return GenerateToken(uid, tokenType)
}

func (s *TokenService) GetTokensByUID(uid int64) ([]*model.TokenRecord, error) {
	return s.repo.GetTokensByUID(uid)
}

func (s *TokenService) GetTokenByID(uid, userTokenID int64) (*model.TokenRecord, error) {
	return s.repo.GetTokenByID(uid, userTokenID)
}

func (s *TokenService) RevokeToken(uid, userTokenID int64) error {
	return s.repo.DeleteToken(uid, userTokenID)
}

func (s *TokenService) RevokeAllTokens(uid int64) error {
	return s.repo.DeleteAllTokens(uid)
}

func (s *TokenService) UpdateLastSeen(uid, userTokenID int64) error {
	return s.repo.UpdateTokenLastSeen(uid, userTokenID, time.Now().Unix())
}
