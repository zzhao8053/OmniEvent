package mail

import (
	"crypto/tls"

	"gopkg.in/mail.v2"

	"omnievent-backend/cmd"
	"omnievent-backend/pkg/errs"
)

// DefaultMailer represents default mailer
type DefaultMailer struct {
	dialer      *mail.Dialer
	fromAddress string
}

// NewDefaultMailer returns a new default mailer
func NewDefaultMailer(smtpConfig *cmd.SMTPConfig) (*DefaultMailer, error) {
	if smtpConfig.Host == "" {
		return nil, errs.ErrSMTPServerHostInvalid
	}

	dialer := mail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Password)
	dialer.TLSConfig = &tls.Config{
		ServerName:         smtpConfig.Host,
		InsecureSkipVerify: smtpConfig.SkipTLSVerify,
	}

	mailer := &DefaultMailer{
		dialer:      dialer,
		fromAddress: smtpConfig.FromAddress,
	}

	return mailer, nil
}

// SendMail sends an email according to argument
func (m *DefaultMailer) SendMail(message *MailMessage) error {
	if m.dialer == nil {
		return errs.ErrSMTPServerNotEnabled
	}

	mailMessage := mail.NewMessage()
	mailMessage.SetHeader("From", m.fromAddress)
	mailMessage.SetHeader("To", message.To)
	mailMessage.SetHeader("Subject", message.Subject)
	mailMessage.SetBody("text/html", message.Body)

	err := m.dialer.DialAndSend(mailMessage)

	return err
}