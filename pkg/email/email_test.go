package email

import (
	"os"
	"testing"

	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/stretchr/testify/assert"
)

var emailCfg *EmailConfig

func TestMain(m *testing.M) {
	emailCfg = &EmailConfig{
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
	sender := NewMockEmailer()
	email := Email{
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
	sender := NewMockEmailer()
	template := EmailTemplate{
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
	sender := NewMockEmailer()
	otpEmail := Email{
		To:      []string{"test@gmail.com"},
		Subject: "Your OTP Code",
		Body:    "Your OTP code is 123456",
	}
	err := sender.SendEmail(otpEmail)
	if err != nil {
		t.Errorf("Error sending OTP email: %v", err)
	}
}

func TestEmail_PrepareBody(t *testing.T) {
	email := Email{
		Name:    "Test User",
		To:      []string{"test@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	expected := "Subject: Test Subject\n\nTest Body"
	result := email.PrepareBody()
	assert.Equal(t, expected, result)
}

//TODO:
// func TestEmailTemplate_PrepareBody(t *testing.T) {
// 	template := EmailTemplate{
// 		To:      []string{"test@example.com"},
// 		Name:    "Test User",
// 		Subject: "Test Subject",
// 		Path:    "testdata/test_template.html",
// 		Data: map[string]interface{}{
// 			"Name": "John Doe",
// 		},
// 	}

// 	// Test with invalid template path
// 	_, _, err := template.PrepareBody()
// 	assert.Error(t, err)

// 	// TODO: Add test with valid template file
// }

func TestNew(t *testing.T) {
	cfg := EmailConfig{
		From:           "test@example.com",
		Host:           "smtp.example.com",
		Port:           587,
		Username:       "testuser",
		Password:       "testpass",
		SendGridAPIKey: "test-api-key",
	}

	sg := sendGridEmailer{}
	smtp := smtpEmailer{}
	// Test SMTP provider
	smtpEmailer := New(providers.EmailSMTP, cfg)
	assert.IsType(t, &smtp, smtpEmailer)

	// Test SendGrid provider (default)
	sendGridEmailer := New(providers.EmailSendGrid, cfg)
	assert.IsType(t, &sg, sendGridEmailer)
}
