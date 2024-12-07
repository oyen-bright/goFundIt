package main

import (
	"fmt"

	appConfig "github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/config/providers"
	"github.com/oyen-bright/goFundIt/internal/database"
	"github.com/oyen-bright/goFundIt/internal/otp"
	"github.com/oyen-bright/goFundIt/pkg/email"
	"github.com/oyen-bright/goFundIt/pkg/encryption"
)

func main() {

	cfg, err := appConfig.LoadConfig()
	if err != nil {
		panic(err)

	}
	db, err := database.Init(*cfg)
	// migrations.DropOtpTable(db)
	defer database.Close(db)
	if err != nil {
		panic(err)
	}

	encryptor := encryption.New(cfg.EncryptionKey)

	emailer := email.New(providers.EmailSMTP, email.EmailConfig{
		Host:           cfg.EmailHost,
		Port:           cfg.EmailPort,
		From:           cfg.EmailName,
		Username:       cfg.EmailUsername,
		Password:       cfg.EmailPassword,
		SendGridAPIKey: cfg.SendGridAPIKey,
	})

	otpService := otp.OtpService{DB: db, Emailer: &emailer, Encryptor: encryptor}
	err = otpService.RequestOTP("bright@krotrust.com", "Bright")
	// isVerified, err := otpService.VerifyOTP("bright@krotrust.com", "MSCKKR")
	// fmt.Print(isVerified)

	fmt.Print(err)

	// router := gin.Default()

}
