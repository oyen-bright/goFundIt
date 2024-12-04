package main

import (
	"log"

	"github.com/oyen-bright/goFundIt/config"
	"github.com/oyen-bright/goFundIt/pkg/email"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)

	}
	log.Println(cfg)

	emailer := email.New(*cfg)
	emailTemplate := email.Verification([]string{"bright@krotrust.com"}, "Aha", "34922343")
	err = emailer.SendEmailTemplate(*emailTemplate)

	if err != nil {
		panic(err)
	}

}
