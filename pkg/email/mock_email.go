package email

type MockEmailer struct{}

var _ Emailer = (*MockEmailer)(nil) // Ensure MockEmailer implements Emailer

func NewMockEmailer() *MockEmailer {
	return &MockEmailer{}
}

func (m *MockEmailer) SendEmail(email Email) error {
	// Mock sending email
	return nil
}

func (m *MockEmailer) SendEmailTemplate(template EmailTemplate) error {
	// Mock sending email template
	return nil
}
