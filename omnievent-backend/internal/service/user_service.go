package service

import (
	"encoding/base64"
	"time"

	"omnievent-backend/cmd"
	"omnievent-backend/internal/model"
	"omnievent-backend/internal/repository"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
	"omnievent-backend/pkg/utils"
)

// UserService represents user service
type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// GetUserByUsernameOrEmailAndPassword returns the user model according to login name and password
func (s *UserService) GetUserByUsernameOrEmailAndPassword(c *context.WebContext, loginname string, password string) (*model.User, error) {
	var user *model.User
	var err error

	if utils.IsValidUsername(loginname) {
		user, err = s.userRepo.GetByUsername(loginname)
	} else if utils.IsValidEmail(loginname) {
		user, err = s.userRepo.GetByEmail(loginname)
	} else {
		return nil, errs.ErrLoginNameInvalid
	}

	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	if !s.IsPasswordEqualsUserPassword(password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	return user, nil
}

// GetUserById returns the user model according to user uid
func (s *UserService) GetUserById(c *context.WebContext, uid int64) (*model.User, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	user, err := s.userRepo.GetByUID(uid)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetUserByUsername returns the user model according to user name
func (s *UserService) GetUserByUsername(c *context.WebContext, username string) (*model.User, error) {
	if username == "" {
		return nil, errs.ErrUsernameIsEmpty
	}

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetUserByEmail returns the user model according to user email
func (s *UserService) GetUserByEmail(c *context.WebContext, email string) (*model.User, error) {
	if email == "" {
		return nil, errs.ErrEmailIsEmpty
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// CreateUser saves a new user model to database
func (s *UserService) CreateUser(c *context.WebContext, user *model.User, noPassword bool) error {
	exists, err := s.ExistsUsername(c, user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errs.ErrUsernameAlreadyExists
	}

	exists, err = s.ExistsEmail(c, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errs.ErrUserEmailAlreadyExists
	}

	if !noPassword && user.Password == "" {
		return errs.ErrPasswordIsEmpty
	}

	user.Salt = utils.GenerateSalt(10)
	user.UID = 0 // Will be auto-generated

	if !noPassword {
		user.Password = utils.EncodePassword(user.Password, user.Salt)
	} else {
		user.Password = ""
	}

	user.Disabled = false
	user.Deleted = false
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	user.LastLoginAt = time.Now().Unix()

	return s.userRepo.Create(user)
}

// UpdateUser saves an existed user model to database
func (s *UserService) UpdateUser(c *context.WebContext, user *model.User) (bool, error) {
	if user.UID <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	if user.Email != "" {
		existing, err := s.userRepo.GetByEmail(user.Email)
		if err != nil && err != repository.ErrNotFound {
			return false, err
		}
		if existing != nil && existing.UID != user.UID {
			return false, errs.ErrUserEmailAlreadyExists
		}
		user.EmailVerified = false
	}

	now := time.Now().Unix()
	user.UpdatedAt = now

	err := s.userRepo.Update(user)
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateUserPassword updates the password of specified user
func (s *UserService) UpdateUserPassword(c *context.WebContext, user *model.User) error {
	if user.UID <= 0 {
		return errs.ErrUserIdInvalid
	}

	if user.Password == "" {
		return errs.ErrPasswordIsEmpty
	}

	user.Password = utils.EncodePassword(user.Password, user.Salt)
	user.UpdatedAt = time.Now().Unix()

	return s.userRepo.Update(user)
}

// UpdateUserLastLoginTime updates the last login time field
func (s *UserService) UpdateUserLastLoginTime(c *context.WebContext, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	return s.userRepo.UpdateLastLogin(uid, time.Now().Unix())
}

// ExistsUsername returns whether the given user name exists
func (s *UserService) ExistsUsername(c *context.WebContext, username string) (bool, error) {
	if username == "" {
		return false, errs.ErrUsernameIsEmpty
	}

	_, err := s.userRepo.GetByUsername(username)
	if err == repository.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ExistsEmail returns whether the given user email exists
func (s *UserService) ExistsEmail(c *context.WebContext, email string) (bool, error) {
	if email == "" {
		return false, errs.ErrEmailIsEmpty
	}

	_, err := s.userRepo.GetByEmail(email)
	if err == repository.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// IsPasswordEqualsUserPassword returns whether the given password is correct
func (s *UserService) IsPasswordEqualsUserPassword(password string, user *model.User) bool {
	return user.Password == utils.EncodePassword(password, user.Salt)
}

// GenerateToken generates a token for user
func GenerateToken(uid int64, tokenType int) (string, *model.TokenRecord, error) {
	secret := utils.GenerateSalt(32)
	combined := secret + string(rune(uid))
	hash := utils.SHA256Bytes([]byte(combined))
	token := base64.URLEncoding.EncodeToString(hash[:])

	expireHours := 24 // default
	cfg := cmd.GetConfig()
	if cfg != nil {
		expireHours = cfg.JWT.ExpireHours
		if expireHours == 0 {
			expireHours = 24
		}
	}

	record := &model.TokenRecord{
		UID:         uid,
		UserTokenID: time.Now().UnixNano(),
		TokenType:   tokenType,
		Secret:      secret,
		UserAgent:   "",
		CreatedAt:   time.Now().Unix(),
		ExpiresAt:   time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
		LastSeenAt:  time.Now().Unix(),
	}

	return token, record, nil
}