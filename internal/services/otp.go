package services

import (
	"net/http"

	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"

	"github.com/oyen-bright/goFundIt/pkg/email"
	emailTemplates "github.com/oyen-bright/goFundIt/pkg/email/templates"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type otpService struct {
	repo     repositories.OTPRepository
	emailer  email.Emailer
	logger   logger.Logger
	runAsync bool
}

func NewOTPService(repo repositories.OTPRepository, emailer email.Emailer, logger logger.Logger) services.OTPService {
	return &otpService{
		repo:     repo,
		emailer:  emailer,
		logger:   logger,
		runAsync: true,
	}
}

// RequestOTP generates a new OTP and returns
func (s *otpService) RequestOTP(email, name string) (models.Otp, error) {

	otp := models.NewOTP(email)
	otp.Name = name

	defer s.repo.InvalidateOtherOTPs(otp.Email, otp.Code, otp.RequestId)
	err := s.repo.Add(otp)
	if err != nil {
		return *otp, err
	}
	otp.Email = email

	if s.runAsync {
		go s.sendOTP(s.emailer, otp.Email, otp.Code, name)
	} else {
		s.sendOTP(s.emailer, otp.Email, otp.Code, name)
	}

	return *otp, nil
}

// VerifyOTP checks if the OTP provided  is valid.
func (s *otpService) VerifyOTP(email, code, requestId string) (models.Otp, error) {

	otp := models.NewOTP(email)

	fetchedOtp, err := s.repo.GetByEmailAndReference(otp.Email, requestId)
	if err != nil {
		return models.Otp{}, errs.New("Invalid OTP", http.StatusNotFound)
	}

	if fetchedOtp.IsExpired() {
		return *fetchedOtp, errs.New("OTP has expired", http.StatusBadRequest)
	}

	if fetchedOtp.Code != code || fetchedOtp.RequestId != requestId {
		return models.Otp{}, errs.New("Invalid OTP", http.StatusBadRequest)
	}

	defer s.repo.Delete(fetchedOtp)
	return *fetchedOtp, nil
}

// Helper Methods -------------------------------------------

// Sends OTP to the user via email
func (s *otpService) sendOTP(emailer email.Emailer, email, code, name string) {
	verificationTemplate := emailTemplates.Verification([]string{email}, name, code)
	err := emailer.SendEmailTemplate(*verificationTemplate)
	if err != nil {
		errs.InternalServerError(err).Log(s.logger)
	}
}
