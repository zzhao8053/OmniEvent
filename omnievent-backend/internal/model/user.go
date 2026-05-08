package model

type User struct {
	UID              int64  `xorm:"pk(uid)"`
	Username         string `xorm:"unique(s) not null"`
	Email            string `xorm:"unique(s) not null"`
	Nickname         string `xorm:"not null"`
	Password         string `xorm:"not null"`
	Salt             string `xorm:"not null"`
	CustomAvatarType string `xorm:"size(10)"`
	AvatarURL        string `xorm:"-"`
	Disabled         bool   `xorm:"default(false)"`
	Deleted          bool   `xorm:"default(false)"`
	EmailVerified    bool   `xorm:"default(false)"`
	CreatedAt        int64  `xorm:"created"`
	UpdatedAt        int64  `xorm:"updated"`
	LastLoginAt      int64  `xorm:"last_login_unix_time"`
}

func (User) TableName() string {
	return "users"
}

type TokenRecord struct {
	UID         int64  `xorm:"pk(uid)"`
	UserTokenID int64  `xorm:"pk(user_token_id)"`
	TokenType   int    `xorm:"not null"`
	Secret      string `xorm:"not null"`
	Token       string `xorm:"-"`
	UserAgent   string `xorm:"size(255)"`
	CreatedAt   int64
	ExpiresAt   int64
	LastSeenAt  int64
}

func (TokenRecord) TableName() string {
	return "token_record"
}

const (
	TokenTypeNormal        = 1
	TokenTypeRequire2FA    = 2
	TokenTypeEmailVerify   = 3
	TokenTypePasswordReset = 4
	TokenTypeMCP           = 5
)
