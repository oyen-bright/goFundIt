package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type OTPRepository interface {
	Add(otp *models.Otp) error
	Delete(otp *models.Otp) error
	InvalidateOtherOTPs(email, code, requestId string) error
	GetByEmailAndReference(email, requestId string) (*models.Otp, error)
}
