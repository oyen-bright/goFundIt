package email_test

import (
	"os"
	"testing"

	"github.com/oyen-bright/goFundIt/pkg/email"
)

var emailCfg *email.EmailConfig

func TestMain(m *testing.M) {
	emailCfg = &email.EmailConfig{
		Host:     "localhost",
		Port:     1025,
		From:     "test",
		Username: "test",
		Password: "test",
	}
	code := m.Run()
	os.Exit(code)
}

func TestSendEmail(t *testing.T) {
	sender := email.NewMockEmailer()
	email := email.Email{
		To:      []string{"test@gmail.com"},
		Subject: "Test email",
		Body:    "This is a test email",
	}
	err := sender.SendEmail(email)
	if err != nil {
		t.Errorf("Error sending email: %v", err)
	}
}

func TestSendEmailTemplate(t *testing.T) {
	sender := email.NewMockEmailer()
	template := email.EmailTemplate{
		To:      []string{"test@gmail.com"},
		Subject: "Test template email",
		Data:    map[string]interface{}{"Name": "Test User"},
		Path:    "/path/to/template",
	}
	err := sender.SendEmailTemplate(template)
	if err != nil {
		t.Errorf("Error sending email template: %v", err)
	}
}

func TestSendOTPEmail(t *testing.T) {
	sender := email.NewMockEmailer()
	otpEmail := email.Email{
		To:      []string{"test@gmail.com"},
		Subject: "Your OTP Code",
		Body:    "Your OTP code is 123456",
	}
	err := sender.SendEmail(otpEmail)
	if err != nil {
		t.Errorf("Error sending OTP email: %v", err)
	}
}
