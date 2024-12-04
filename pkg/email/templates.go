package email

func Verification(to []string, name, verificationCode string) *EmailTemplate {
	return &EmailTemplate{
		To:      to,
		Subject: "Email Verification - GoFund It",
		Path:    "./templates/email_verification.html",
		Data: map[string]interface{}{
			"verificationCode": verificationCode,
		},
	}
}
