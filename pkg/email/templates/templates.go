package emailTemplates

import (
	"path/filepath"

	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/pkg/email/models"
)

func generateFile(fileName string) string {
	//TODO:implement better file path generation and handling
	return filepath.Join(config.BaseDir, "pkg", "email", "templates", fileName)
}

func Verification(to []string, name, verificationCode string) *models.EmailTemplate {

	return &models.EmailTemplate{
		To:      to,
		Subject: "Email Verification - GoFund It",
		Path:    generateFile("email_verification.html"),
		Data: map[string]interface{}{
			"verificationCode": verificationCode,
		},
	}
}
