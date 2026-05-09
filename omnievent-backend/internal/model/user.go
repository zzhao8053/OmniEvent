package model

// TransactionEditScope represents the scope which transaction can be edited
type TransactionEditScope byte

// Editable Transaction Ranges
const (
	TransactionEditScopeNone              TransactionEditScope = 0
	TransactionEditScopeAll               TransactionEditScope = 1
	TransactionEditScopeTodayOrLater      TransactionEditScope = 2
	TransactionEditScopeLast24HourOrLater TransactionEditScope = 3
	TransactionEditScopeThisWeekOrLater   TransactionEditScope = 4
	TransactionEditScopeThisMonthOrLater  TransactionEditScope = 5
	TransactionEditScopeThisYearOrLater   TransactionEditScope = 6
	TransactionEditScopeInvalid           TransactionEditScope = 255
)

// AmountColorType represents the type of amount color in frontend
type AmountColorType byte

// Amount Color Types
const (
	AmountColorTypeDefault        AmountColorType = 0
	AmountColorTypeGreen          AmountColorType = 1
	AmountColorTypeRed            AmountColorType = 2
	AmountColorTypeYellow         AmountColorType = 3
	AmountColorTypeBlackOrWhite   AmountColorType = 4
	AmountColorTypeInvalid        AmountColorType = 255
)

// User represents user data stored in database
type User struct {
	UID                   int64                `xorm:"pk(uid)"`
	Username              string               `xorm:"unique(s) not null"`
	Email                 string               `xorm:"unique(s) not null"`
	Nickname              string               `xorm:"not null"`
	Password              string               `xorm:"not null"`
	Salt                  string               `xorm:"not null"`
	CustomAvatarType      string               `xorm:"size(10)"`
	AvatarURL             string               `xorm:"-"`
	TransactionEditScope  TransactionEditScope `xorm:"tinyint not null"`
	Language              string               `xorm:"size(10)"`
	DefaultCurrency       string               `xorm:"size(3) not null"`
	FirstDayOfWeek        byte                 `xorm:"tinyint not null"`
	FiscalYearStart       byte                 `xorm:"smallint"`
	CalendarDisplayType   byte                 `xorm:"tinyint"`
	DateDisplayType       byte                 `xorm:"tinyint"`
	LongDateFormat        byte                 `xorm:"tinyint"`
	ShortDateFormat       byte                 `xorm:"tinyint"`
	LongTimeFormat        byte                 `xorm:"tinyint"`
	ShortTimeFormat       byte                 `xorm:"tinyint"`
	FiscalYearFormat      byte                 `xorm:"tinyint"`
	CurrencyDisplayType   byte                 `xorm:"tinyint"`
	NumeralSystem         byte                 `xorm:"tinyint"`
	DecimalSeparator      byte                 `xorm:"tinyint"`
	DigitGroupingSymbol   byte                 `xorm:"tinyint"`
	DigitGrouping         byte                 `xorm:"tinyint"`
	CoordinateDisplayType byte                 `xorm:"tinyint"`
	ExpenseAmountColor    AmountColorType      `xorm:"tinyint"`
	IncomeAmountColor     AmountColorType      `xorm:"tinyint"`
	FeatureRestriction    byte                 `xorm:"tinyint"`
	Disabled              bool                 `xorm:"default(false)"`
	Deleted               bool                 `xorm:"default(false)"`
	EmailVerified         bool                 `xorm:"default(false)"`
	CreatedAt             int64                `xorm:"created"`
	UpdatedAt             int64                `xorm:"updated"`
	DeletedAt             int64                `xorm:"deleted_unix_time"`
	LastLoginAt           int64                `xorm:"last_login_unix_time"`
}

func (User) TableName() string {
	return "users"
}

// UserBasicInfo represents a view-object of user basic info
type UserBasicInfo struct {
	Username              string                     `json:"username"`
	Email                 string                     `json:"email"`
	Nickname              string                     `json:"nickname"`
	AvatarURL             string                     `json:"avatar"`
	AvatarProvider        string                     `json:"avatarProvider,omitempty"`
	TransactionEditScope  TransactionEditScope       `json:"transactionEditScope"`
	Language              string                     `json:"language"`
	DefaultCurrency       string                     `json:"defaultCurrency"`
	FirstDayOfWeek        byte                       `json:"firstDayOfWeek"`
	FiscalYearStart       byte                       `json:"fiscalYearStart"`
	CalendarDisplayType   byte                       `json:"calendarDisplayType"`
	DateDisplayType       byte                       `json:"dateDisplayType"`
	LongDateFormat        byte                       `json:"longDateFormat"`
	ShortDateFormat       byte                       `json:"shortDateFormat"`
	LongTimeFormat        byte                       `json:"longTimeFormat"`
	ShortTimeFormat       byte                       `json:"shortTimeFormat"`
	FiscalYearFormat      byte                       `json:"fiscalYearFormat"`
	CurrencyDisplayType   byte                       `json:"currencyDisplayType"`
	NumeralSystem         byte                       `json:"numeralSystem"`
	DecimalSeparator      byte                       `json:"decimalSeparator"`
	DigitGroupingSymbol   byte                       `json:"digitGroupingSymbol"`
	DigitGrouping         byte                       `json:"digitGrouping"`
	CoordinateDisplayType byte                       `json:"coordinateDisplayType"`
	ExpenseAmountColor    AmountColorType            `json:"expenseAmountColor"`
	IncomeAmountColor     AmountColorType            `json:"incomeAmountColor"`
	EmailVerified         bool                       `json:"emailVerified"`
}

// TokenRecord represents token data stored in database
type TokenRecord struct {
	UID         int64  `xorm:"pk(uid)"`
	UserTokenID int64  `xorm:"pk(user_token_id)"`
	TokenType   int    `xorm:"tinyint not null"`
	Secret      string `xorm:"size(10) not null"`
	Token       string `xorm:"-"`
	UserAgent   string `xorm:"size(255)"`
	Context     []byte `xorm:"blob"`
	CreatedAt   int64  `xorm:"created"`
	ExpiresAt   int64  `xorm:"index(index_token_record_expired_time)"`
	LastSeenAt  int64
}

func (TokenRecord) TableName() string {
	return "token_record"
}

// Token types
const (
	TokenTypeNormal        = 1
	TokenTypeRequire2FA    = 2
	TokenTypeEmailVerify   = 3
	TokenTypePasswordReset = 4
	TokenTypeMCP           = 5
	TokenTypeOAuth2CallbackRequireVerify = 6
	TokenTypeOAuth2Callback = 7
	TokenTypeAPI           = 8
)

// OAuth2CallbackTokenContext represents the context data of oauth 2.0 callback token
type OAuth2CallbackTokenContext struct {
	ExternalAuthType int    `json:"externalAuthType"`
	ExternalUsername string `json:"externalUsername"`
	ExternalEmail    string `json:"externalEmail"`
}

// TokenInfoResponse represents a view-object of token
type TokenInfoResponse struct {
	TokenID   string `json:"tokenId"`
	TokenType int    `json:"tokenType"`
	UserAgent string `json:"userAgent"`
	LastSeen  int64  `json:"lastSeen"`
	IsCurrent bool   `json:"isCurrent"`
}