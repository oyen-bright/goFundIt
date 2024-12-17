package interfaces

import (
	"github.com/oyen-bright/goFundIt/internal/models"
)

type OTPService interface {
	RequestOTP(email, name string) (models.Otp, error)
	VerifyOTP(email, otp, requestId string) (models.Otp, error)
}
