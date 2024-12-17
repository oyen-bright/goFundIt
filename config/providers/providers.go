package providers

type (
	EmailProvider int
	PhoneProvider int
)

type EmailProviderType string

type PhoneProviderType string

const (
	// Email providers
	SMTP     EmailProviderType = "smtp"
	SendGrid EmailProviderType = "sendgrid"

	// Phone providers
	ProviderA PhoneProviderType = "A"
	ProviderB PhoneProviderType = "B"
)

const (
	EmailSMTP EmailProvider = iota
	EmailSendGrid
)

const (
	PhoneA PhoneProvider = iota
	PhoneB
)

type Provider interface {
	String() string
}

type EmailProviderSetter interface {
	Provider
	SetEmailProvider(EmailProviderType)
}

type PhoneProviderSetter interface {
	Provider
	SetPhoneProvider(PhoneProviderType)
}

var emailProviders = [...]string{
	EmailSMTP:     string(SMTP),
	EmailSendGrid: string(SendGrid),
}

var phoneProviders = [...]string{
	PhoneA: string(ProviderA),
	PhoneB: string(ProviderB),
}

func (e EmailProvider) String() string {
	if e < 0 || int(e) >= len(emailProviders) {
		return string(SendGrid) // Default provider
	}
	return emailProviders[e]
}

func (e *EmailProvider) SetEmailProvider(provider EmailProviderType) {
	switch provider {
	case SMTP:
		*e = EmailSMTP
	default:
		*e = EmailSendGrid
	}
}

func (p PhoneProvider) String() string {
	if p < 0 || int(p) >= len(phoneProviders) {
		return string(ProviderB)
	}
	return phoneProviders[p]
}

func (p *PhoneProvider) SetPhoneProvider(provider PhoneProviderType) {
	switch provider {
	case ProviderA:
		*p = PhoneA
	default:
		*p = PhoneB
	}
}

func NewEmailProvider(providerType EmailProviderType) EmailProvider {
	var provider EmailProvider
	provider.SetEmailProvider(providerType)
	return provider
}

func NewPhoneProvider(providerType PhoneProviderType) PhoneProvider {
	var provider PhoneProvider
	provider.SetPhoneProvider(providerType)
	return provider
}
