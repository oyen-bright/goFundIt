package email

import (
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/pkg/email/models"
)

type EmailConfig struct {
	From           string
	Host           string
	Port           int
	Username       string
	Password       string
	SendGridAPIKey string
}

type Emailer interface {
	SendEmail(models.Email) error
	SendEmailTemplate(models.EmailTemplate) error
}

func New(provider providers.EmailProvider, cfg EmailConfig) Emailer {
	switch provider {
	case providers.EmailSMTP:
		return &smtpEmailer{
			config: cfg,
		}

	default:
		return &sendGridEmailer{
			config: cfg,
			key:    cfg.SendGridAPIKey,
		}
	}
}
