package service

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strings"

	"gopkg.in/mail.v2"

	"omnievent-backend/cmd"
	"omnievent-backend/internal/model"
	"omnievent-backend/pkg/context"
	"omnievent-backend/pkg/errs"
)

const passwordResetUrlFormat = "%sdesktop#/resetpassword?token=%s"

// SMTPConfig represents SMTP configuration
type SMTPConfig struct {
	Host          string
	Port          int
	User          string
	Password      string
	FromAddress   string
	SkipTLSVerify bool
	EnableSMTP    bool
}

// MailMessage represents an email entity
type MailMessage struct {
	To      string
	Subject string
	Body    string
}

// Mailer is email sender interface
type Mailer interface {
	SendMail(message *MailMessage) error
}

// DefaultMailer represents default mailer
type DefaultMailer struct {
	dialer      *mail.Dialer
	fromAddress string
}

// NewDefaultMailer creates a new default mailer
func NewDefaultMailer(smtpConfig *SMTPConfig) (*DefaultMailer, error) {
	host, portStr, err := net.SplitHostPort(smtpConfig.Host)
	if err != nil {
		return nil, errs.ErrSMTPServerHostInvalid
	}

	port, err := net.LookupPort("tcp", portStr)
	if err != nil {
		return nil, errs.ErrSMTPServerHostInvalid
	}

	dialer := mail.NewDialer(host, port, smtpConfig.User, smtpConfig.Password)
	dialer.TLSConfig = &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: smtpConfig.SkipTLSVerify,
	}

	return &DefaultMailer{
		dialer:      dialer,
		fromAddress: smtpConfig.FromAddress,
	}, nil
}

// SendMail sends an email
func (m *DefaultMailer) SendMail(message *MailMessage) error {
	if m.dialer == nil {
		return errs.ErrSMTPServerNotEnabled
	}

	mailMessage := mail.NewMessage()
	mailMessage.SetHeader("From", m.fromAddress)
	mailMessage.SetHeader("To", message.To)
	mailMessage.SetHeader("Subject", message.Subject)
	mailMessage.SetBody("text/html", message.Body)

	return m.dialer.DialAndSend(mailMessage)
}

// ForgetPasswordService represents forget password service
type ForgetPasswordService struct {
	mailer Mailer
}

// Initialize a forget password service singleton instance
var forgetPasswordService *ForgetPasswordService

// GetForgetPasswordService returns the forget password service singleton
func GetForgetPasswordService() *ForgetPasswordService {
	if forgetPasswordService == nil {
		cfg := cmd.GetConfig()
		var mailer Mailer

		if cfg != nil && cfg.SMTP.EnableSMTP {
			smtpConfig := &SMTPConfig{
				Host:          cfg.SMTP.Host,
				Port:          cfg.SMTP.Port,
				User:          cfg.SMTP.User,
				Password:      cfg.SMTP.Password,
				FromAddress:   cfg.SMTP.FromAddress,
				SkipTLSVerify: cfg.SMTP.SkipTLSVerify,
				EnableSMTP:    cfg.SMTP.EnableSMTP,
			}
			defaultMailer, err := NewDefaultMailer(smtpConfig)
			if err == nil {
				mailer = defaultMailer
			}
		}

		forgetPasswordService = &ForgetPasswordService{
			mailer: mailer,
		}
	}
	return forgetPasswordService
}

// SendPasswordResetEmail sends password reset email according to specified parameters
func (s *ForgetPasswordService) SendPasswordResetEmail(c *context.WebContext, user *model.User, passwordResetToken string) error {
	if s.mailer == nil {
		return errs.ErrSMTPServerNotEnabled
	}

	cfg := cmd.GetConfig()
	if cfg == nil || !cfg.SMTP.EnableSMTP {
		return errs.ErrSMTPServerNotEnabled
	}

	expireTimeInMinutes := 60 // default, should come from config
	if cfg.Security.TokenExpireSeconds > 0 {
		expireTimeInMinutes = cfg.Security.TokenExpireSeconds / 60
	}

	rootURL := cfg.Server.RootURL
	if rootURL == "" {
		rootURL = fmt.Sprintf("%s://%s:%d", cfg.Server.Protocol, cfg.Server.HttpAddr, cfg.Server.HttpPort)
	}

	passwordResetUrl := fmt.Sprintf(passwordResetUrlFormat, strings.TrimSuffix(rootURL, "/"), url.QueryEscape(passwordResetToken))

	// Get locale text items
	locale := user.Language
	if locale == "" {
		locale = "zh_CN" // default locale
	}

	appName := cfg.Global.AppName
	if appName == "" {
		appName = "OmniEvent"
	}

	mailTitle := getLocalizedText(locale, "password_reset_mail_title", "Reset Your Password")
	salutation := fmt.Sprintf(getLocalizedText(locale, "password_reset_salutation", "Dear %s,"), user.Nickname)
	descriptionAbove := getLocalizedText(locale, "password_reset_description_above", "You have requested to reset your password. Please click the button below to proceed.")
	resetButton := getLocalizedText(locale, "password_reset_button", "Reset Password")
	descriptionBelow := fmt.Sprintf(getLocalizedText(locale, "password_reset_description_below", "This link will expire in %d minutes. If you did not request a password reset, please ignore this email."), expireTimeInMinutes)

	// Build email body
	body := buildPasswordResetEmailBody(appName, mailTitle, salutation, descriptionAbove, passwordResetUrl, resetButton, descriptionBelow)

	message := &MailMessage{
		To:      user.Email,
		Subject: mailTitle,
		Body:    body,
	}

	return s.mailer.SendMail(message)
}

func getLocalizedText(locale, key, defaultValue string) string {
	// Simple localization - in production this would use a proper i18n system
	// For now, return default values based on locale
	translations := map[string]map[string]string{
		"zh_CN": {
			"password_reset_mail_title":     "重置您的密码",
			"password_reset_salutation":     "亲爱的 %s，",
			"password_reset_description_above": "您已申请重置密码。请点击下方按钮继续操作。",
			"password_reset_button":         "重置密码",
			"password_reset_description_below": "此链接将在 %d 分钟内失效。如果您未申请重置密码，请忽略此邮件。",
		},
		"en": {
			"password_reset_mail_title":     "Reset Your Password",
			"password_reset_salutation":     "Dear %s,",
			"password_reset_description_above": "You have requested to reset your password. Please click the button below to proceed.",
			"password_reset_button":         "Reset Password",
			"password_reset_description_below": "This link will expire in %d minutes. If you did not request a password reset, please ignore this email.",
		},
	}

	if localeTranslations, ok := translations[locale]; ok {
		if text, ok := localeTranslations[key]; ok {
			return text
		}
	}

	// Fallback to English if locale not found
	if locale != "en" && locale != "zh_CN" {
		if enTranslations, ok := translations["en"]; ok {
			if text, ok := enTranslations[key]; ok {
				return text
			}
		}
	}

	return defaultValue
}

func buildPasswordResetEmailBody(appName, title, salutation, descriptionAbove, resetUrl, resetButton, descriptionBelow string) string {
	var buf bytes.Buffer
	buf.WriteString("<!DOCTYPE html>\n")
	buf.WriteString("<html>\n")
	buf.WriteString("<head>\n")
	buf.WriteString("    <meta charset=\"utf-8\">\n")
	buf.WriteString("    <meta http-equiv=\"Content-Type\" content=\"text/html;charset=utf-8\"/>\n")
	buf.WriteString("    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n")
	buf.WriteString("    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no, minimal-ui, viewport-fit=cover\">\n")
	buf.WriteString("    <title>")
	buf.WriteString(title)
	buf.WriteString("</title>\n")
	buf.WriteString("</head>\n")
	buf.WriteString("<body style=\"margin: 0; padding: 0 10px 0 10px\">\n")
	buf.WriteString("    <table width=\"360px\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\" style=\"width: 360px; border: 0; border-collapse: collapse; margin: 10px auto 5px auto;\">\n")
	buf.WriteString("        <tr>\n")
	buf.WriteString("            <td height=\"50\" style=\"font-size: 20px; line-height: 50px\"><strong>")
	buf.WriteString(appName)
	buf.WriteString("</strong></td>\n")
	buf.WriteString("        </tr>\n")
	buf.WriteString("        <tr>\n")
	buf.WriteString("            <td style=\"padding: 10px 0 10px 0; border-top: solid 1px #ccc\">\n")
	buf.WriteString("                <p>")
	buf.WriteString(salutation)
	buf.WriteString("</p>\n")
	buf.WriteString("                <p>")
	buf.WriteString(descriptionAbove)
	buf.WriteString("</p>\n")
	buf.WriteString("            </td>\n")
	buf.WriteString("        </tr>\n")
	buf.WriteString("        <tr>\n")
	buf.WriteString("            <td height=\"50\" style=\"line-height: 50px; text-align: center\">\n")
	buf.WriteString("                <a href=\"")
	buf.WriteString(resetUrl)
	buf.WriteString("\" style=\"width: 100%; color: #fff; background-color:#c67e48; display:block\">\n")
	buf.WriteString("                    <strong>")
	buf.WriteString(resetButton)
	buf.WriteString("</strong>\n")
	buf.WriteString("                </a>\n")
	buf.WriteString("            </td>\n")
	buf.WriteString("        </tr>\n")
	buf.WriteString("        <tr>\n")
	buf.WriteString("            <td style=\"padding: 10px 0 10px 0\">\n")
	buf.WriteString("                <p>")
	buf.WriteString(descriptionBelow)
	buf.WriteString("</p>\n")
	buf.WriteString("            </td>\n")
	buf.WriteString("        </tr>\n")
	buf.WriteString("        <tr>\n")
	buf.WriteString("            <td style=\"padding-bottom: 20px\">\n")
	buf.WriteString("                <small style=\"color: #888; word-break: break-all\">")
	buf.WriteString(resetUrl)
	buf.WriteString("</small>\n")
	buf.WriteString("            </td>\n")
	buf.WriteString("        </tr>\n")
	buf.WriteString("    </table>\n")
	buf.WriteString("</body>\n")
	buf.WriteString("</html>\n")
	return buf.String()
}