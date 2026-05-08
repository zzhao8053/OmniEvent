package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"omnievent-backend/cmd"
	"omnievent-backend/internal/model"
	"omnievent-backend/internal/repository"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, email, nickname, password string) (*model.User, error) {
	_, err := s.repo.GetByUsername(username)
	if err == nil {
		return nil, ErrUserExists
	}
	_, err = s.repo.GetByEmail(email)
	if err == nil {
		return nil, ErrUserExists
	}

	salt := generateSalt(10)
	encodedPassword := encodePassword(password, salt)

	user := &model.User{
		Username:      username,
		Email:         email,
		Nickname:      nickname,
		Password:      encodedPassword,
		Salt:          salt,
		EmailVerified: false,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByUID(uid int64) (*model.User, error) {
	return s.repo.GetByUID(uid)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UserService) GetUserByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	return s.repo.GetByUsernameOrEmail(usernameOrEmail)
}

func (s *UserService) ValidateCredentials(usernameOrEmail, password string) (*model.User, error) {
	user, err := s.repo.GetByUsernameOrEmail(usernameOrEmail)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	encodedPassword := encodePassword(password, user.Salt)
	if encodedPassword != user.Password {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) UpdateProfile(uid int64, nickname, email string) (*model.User, error) {
	user, err := s.repo.GetByUID(uid)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.Nickname = nickname
	user.Email = email
	user.UpdatedAt = time.Now().Unix()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateAvatar(uid int64, avatarType string) (*model.User, error) {
	user, err := s.repo.GetByUID(uid)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.CustomAvatarType = avatarType
	user.UpdatedAt = time.Now().Unix()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) RemoveAvatar(uid int64) (*model.User, error) {
	user, err := s.repo.GetByUID(uid)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.CustomAvatarType = ""
	user.UpdatedAt = time.Now().Unix()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateLastLogin(uid int64) error {
	return s.repo.UpdateLastLogin(uid, time.Now().Unix())
}

func generateSalt(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func encodePassword(password, salt string) string {
	combined := password + salt
	hash := sha256.Sum256([]byte(combined))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func generateTokenID() int64 {
	return time.Now().UnixNano()
}

func GenerateToken(uid int64, tokenType int) (string, *model.TokenRecord, error) {
	secret := generateSalt(32)
	combined := secret + string(rune(uid))
	hash := sha256.Sum256([]byte(combined))
	token := base64.URLEncoding.EncodeToString(hash[:])

	cfg := cmd.GetConfig()
	expireHours := cfg.JWT.ExpireHours
	if expireHours == 0 {
		expireHours = 24
	}

	record := &model.TokenRecord{
		UID:         uid,
		UserTokenID: generateTokenID(),
		TokenType:   tokenType,
		Secret:      secret,
		UserAgent:   "",
		CreatedAt:   time.Now().Unix(),
		ExpiresAt:   time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
		LastSeenAt:  time.Now().Unix(),
	}

	return token, record, nil
}
