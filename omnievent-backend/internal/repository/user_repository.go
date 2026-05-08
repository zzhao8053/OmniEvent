package repository

import (
	"errors"

	"omnievent-backend/cmd"
	"omnievent-backend/internal/model"

	"xorm.io/xorm"
)

var ErrNotFound = errors.New("record not found")

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetDB() *xorm.Engine {
	return cmd.GetDB()
}

func (r *UserRepository) Create(user *model.User) error {
	_, err := r.GetDB().Insert(user)
	return err
}

func (r *UserRepository) GetByUID(uid int64) (*model.User, error) {
	var user model.User
	has, err := r.GetDB().Where("uid = ?", uid).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	has, err := r.GetDB().Where("username = ?", username).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	has, err := r.GetDB().Where("email = ?", email).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *UserRepository) GetByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	var user model.User
	has, err := r.GetDB().Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	_, err := r.GetDB().Where("uid = ?", user.UID).Update(user)
	return err
}

func (r *UserRepository) UpdateLastLogin(uid int64, timestamp int64) error {
	_, err := r.GetDB().Table(&model.User{}).Where("uid = ?", uid).Update(map[string]interface{}{
		"last_login_unix_time": timestamp,
	})
	return err
}

func (r *UserRepository) CreateToken(token *model.TokenRecord) error {
	_, err := r.GetDB().Insert(token)
	return err
}

func (r *UserRepository) GetTokenByID(uid, userTokenID int64) (*model.TokenRecord, error) {
	var token model.TokenRecord
	has, err := r.GetDB().Where("uid = ? AND user_token_id = ?", uid, userTokenID).Get(&token)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return &token, nil
}

func (r *UserRepository) GetTokensByUID(uid int64) ([]*model.TokenRecord, error) {
	var tokens []*model.TokenRecord
	err := r.GetDB().Where("uid = ?", uid).Find(&tokens)
	return tokens, err
}

func (r *UserRepository) DeleteToken(uid, userTokenID int64) error {
	_, err := r.GetDB().Where("uid = ? AND user_token_id = ?", uid, userTokenID).Delete(&model.TokenRecord{})
	return err
}

func (r *UserRepository) DeleteAllTokens(uid int64) error {
	_, err := r.GetDB().Where("uid = ?", uid).Delete(&model.TokenRecord{})
	return err
}

func (r *UserRepository) UpdateTokenLastSeen(uid, userTokenID int64, timestamp int64) error {
	_, err := r.GetDB().Table(&model.TokenRecord{}).
		Where("uid = ? AND user_token_id = ?", uid, userTokenID).
		Update(map[string]interface{}{
			"last_seen_unix_time": timestamp,
		})
	return err
}
