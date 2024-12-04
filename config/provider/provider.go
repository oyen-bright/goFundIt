package providers

type EmailProvider int
type PhoneProvider int

const (
	EmailSMTP EmailProvider = iota
	EmailSendGrid
)

const (
	A PhoneProvider = iota
	B
)

type Provider interface {
	String() string
	Email(string)
	Phone(string)
}

func (e EmailProvider) String() string {
	emailProviders := [...]string{"smtp", "sendgrid"}
	return emailProviders[e]
}

func (e *EmailProvider) Email(provider string) {
	switch provider {
	case "smtp":
		*e = EmailSMTP
	default:
		*e = EmailSendGrid
	}
}

func (p *PhoneProvider) Phone(provider string) {
	switch provider {
	case "A":
		*p = A
	default:
		*p = B
	}
}
func (p PhoneProvider) String() string {
	phoneProviders := [...]string{"A", "B"}
	return phoneProviders[p]
}
