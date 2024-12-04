package email

import (
	"gopkg.in/gomail.v2"
)

type smtpEmailer struct {
	config Config
}

func (s *smtpEmailer) send(m *gomail.Message) error {
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	return d.DialAndSend(m)
}

func (s *smtpEmailer) prepareMessage(from string, to []string, subject, body string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return m
}

func (s *smtpEmailer) SendEmail(email Email) error {
	m := s.prepareMessage(s.config.From, email.To, email.Subject, email.Body)
	return s.send(m)
}

func (s *smtpEmailer) SendEmailTemplate(eTemplate EmailTemplate) error {
	_, body, err := eTemplate.PrepareBody()
	if err != nil {
		return err
	}
	m := s.prepareMessage(s.config.From, eTemplate.To, eTemplate.Subject, body)
	return s.send(m)
}
