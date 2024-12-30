package email

import "errors"

type mockEmailer struct {
	sendEmailCalled         bool
	sendEmailTemplateCalled bool
	shouldError             bool
}

func (m *mockEmailer) SendEmail(e Email) error {
	m.sendEmailCalled = true
	if m.shouldError {
		return errors.New("mock error")
	}
	return nil
}

func (m *mockEmailer) SendEmailTemplate(e EmailTemplate) error {
	m.sendEmailTemplateCalled = true
	if m.shouldError {
		return errors.New("mock error")
	}
	return nil
}
