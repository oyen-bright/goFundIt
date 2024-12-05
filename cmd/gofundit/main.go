package main

import (
	"log"

	appConfig "github.com/oyen-bright/goFundIt/config"
	encryptor "github.com/oyen-bright/goFundIt/internal/encryption"
	"github.com/oyen-bright/goFundIt/pkg/email"
)

func main() {

	cfg, err := appConfig.LoadConfig()
	if err != nil {
		panic(err)

	}

	encryptedData, err := encryptor.Encrypt(cfg.EncryptionKey, encryptor.Data{
		Email: "test@email.com",
		Data:  "1",
	})
	if err != nil {
		panic(err)
	}
	log.Println(encryptedData)

	decripredData, err := encryptor.Decrypt(cfg.EncryptionKey, encryptor.Data{
		Email: "test@email.com",
		Data:  encryptedData,
	})

	log.Println(decripredData)

	if err != nil {
		panic(err)
	}

	emailer := email.New(*cfg)
	emailTemplate := email.Verification([]string{"bright@krotrust.com"}, "Aha", encryptedData)
	err = emailer.SendEmailTemplate(*emailTemplate)

	if err != nil {
		panic(err)
	}

}
