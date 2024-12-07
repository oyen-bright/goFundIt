package models

type Email struct {
	Name    string
	To      []string
	Subject string
	Body    string
}

func (e Email) PrepareBody() string {
	return "Subject: " + e.Subject + "\n\n" + e.Body

}
