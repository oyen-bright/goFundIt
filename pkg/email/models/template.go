package models

import (
	"bytes"
	"text/template"
)

type EmailTemplate struct {
	To      []string
	Name    string
	Subject string
	Path    string
	Data    map[string]interface{}
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
