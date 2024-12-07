package email

import (
	"github.com/oyen-bright/goFundIt/pkg/email/models"
)

type MockEmailer struct{}

var _ Emailer = (*MockEmailer)(nil) // Ensure MockEmailer implements Emailer

func NewMockEmailer() *MockEmailer {
	return &MockEmailer{}
}

func (m *MockEmailer) SendEmail(email models.Email) error {
	// Mock sending email
	return nil
}

func (m *MockEmailer) SendEmailTemplate(template models.EmailTemplate) error {
	// Mock sending email template
	return nil
}
