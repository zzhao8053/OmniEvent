package mail

import (
	"omnievent-backend/cmd"
	"omnievent-backend/pkg/errs"
)

// MailerContainer contains the current mailer
type MailerContainer struct {
	current Mailer
}

// Initialize a mailer container singleton instance
var (
	Container = &MailerContainer{}
)

// InitializeMailer initializes the current mailer according to the config
func InitializeMailer() error {
	smtpConfig := cmd.GetConfig().SMTP

	if !smtpConfig.EnableSMTP {
		Container.current = nil
		return nil
	}

	mailer, err := NewDefaultMailer(&smtpConfig)

	if err != nil {
		return err
	}

	Container.current = mailer
	return nil
}

// SendMail sends an email according to argument
func (m *MailerContainer) SendMail(message *MailMessage) error {
	if m.current == nil {
		return errs.ErrSMTPServerNotEnabled
	}

	return m.current.SendMail(message)
}