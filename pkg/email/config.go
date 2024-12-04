package email

import (
	"github.com/oyen-bright/goFundIt/config"
	providers "github.com/oyen-bright/goFundIt/config/provider"
)

type Config struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

type Emailer interface {
	SendEmail(email Email) error
	SendEmailTemplate(EmailTemplate) error
}

func New(cfg config.AppConfig) Emailer {
	emailCfg := Config{
		Host:     cfg.EmailHost,
		Port:     cfg.EmailPort,
		From:     cfg.EmailName,
		Username: cfg.EmailUsername,
		Password: cfg.EmailPassword,
	}

	switch cfg.EmailProvider {
	case providers.EmailSMTP:
		return &smtpEmailer{
			config: emailCfg,
		}
	default:
		return &sendGridEmailer{
			config: emailCfg,
			key:    cfg.SendGridAPIKey,
		}
	}
}
