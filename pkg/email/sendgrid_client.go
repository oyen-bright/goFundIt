package email

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendGridEmailer struct {
	config Config
	key    string
}

func (s *sendGridEmailer) prepareMessage(from string, to []string, name, subject, body string) *mail.SGMailV3 {
	log.Println("Preparing message")
	fromEmail := mail.NewEmail(name, from)
	toEmails := make([]*mail.Email, len(to))
	for i, recipient := range to {
		toEmails[i] = mail.NewEmail("", recipient)
	}
	message := mail.NewSingleEmail(fromEmail, subject, toEmails[0], body, body)

	log.Println(message)
	return message
}

func (s *sendGridEmailer) send(m *mail.SGMailV3) error {
	client := sendgrid.NewSendClient(s.key)
	response, err := client.Send(m)
	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)

	return err
}

func (s *sendGridEmailer) SendEmail(email Email) error {
	m := s.prepareMessage(s.config.From, email.To, email.Name, email.Subject, email.Body)
	return s.send(m)

}

func (s *sendGridEmailer) SendEmailTemplate(eTemplate EmailTemplate) error {
	_, body, err := eTemplate.PrepareBody()
	if err != nil {
		return err
	}
	m := s.prepareMessage(s.config.From, eTemplate.To, eTemplate.Name, eTemplate.Subject, body)
	return s.send(m)

}
