package services

import (
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/email"
)

type notificationService struct {
	emailer email.Emailer
}

func NewNotificationService(emailer email.Emailer) interfaces.NotificationService {
	return &notificationService{emailer: emailer}
}
