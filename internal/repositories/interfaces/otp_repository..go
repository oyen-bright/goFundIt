package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type OTPRepository interface {
	GetByEmailAndReference(email, requestId string) (*models.Otp, error)
	Add(otp *models.Otp) error
	InvalidateOtherOTPs(email, code, requestId string) error
	Delete(otp *models.Otp) error
}
