package email

import (
	"bytes"
	"text/template"

	"github.com/oyen-bright/goFundIt/config/providers"
)

type EmailConfig struct {
	From           string
	Host           string
	Port           int
	Username       string
	Password       string
	SendGridAPIKey string
}
type Email struct {
	Name        string
	To          []string
	Subject     string
	Body        string
	Attachments []string
}

func (e Email) PrepareBody() string {
	return "Subject: " + e.Subject + "\n\n" + e.Body

}

type EmailTemplate struct {
	To          []string
	Name        string
	Subject     string
	Path        string
	Attachments []string
	Data        map[string]interface{}
}

func (e *EmailTemplate) PrepareBody() (string, string, error) {
	var body bytes.Buffer
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	t, err := template.ParseFiles(e.Path)
	if err != nil {
		return "", "", err
	}
	err = t.Execute(&body, e.Data)
	if err != nil {
		return "", "", err
	}
	return "Subject:" + e.Subject + "\n" + headers + "\n\n" + body.String(), body.String(), nil
}

type Emailer interface {
	SendEmail(Email) error
	SendEmailTemplate(EmailTemplate) error
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
