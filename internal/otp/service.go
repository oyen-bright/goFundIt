package otp

import (
	"net/http"

	"github.com/oyen-bright/goFundIt/pkg/email"
	emailTemplates "github.com/oyen-bright/goFundIt/pkg/email/templates"
	encryptor "github.com/oyen-bright/goFundIt/pkg/encryption"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type OTPService interface {
	RequestOTP(email, name string) (Otp, error)
	VerifyOTP(email, otp, requestId string) (Otp, error)
}

type otpService struct {
	otpRepo   OTPRepository
	emailer   email.Emailer
	encryptor encryptor.Encryptor
	logger    logger.Logger
}

func Service(otpRepo OTPRepository, emailer email.Emailer, encryptor encryptor.Encryptor, logger logger.Logger) OTPService {
	return &otpService{
		otpRepo:   otpRepo,
		emailer:   emailer,
		logger:    logger,
		encryptor: encryptor,
	}
}

// RequestOTP generates a new OTP and sends it to the user.
//
//  1. Clears any previous OTPs for the user from the database.
//  2. Inserts a new OTP into the database.
//  3. Sends the OTP to the user via email.
//
// It returns an error if any of the steps fail.
func (s *otpService) RequestOTP(email, name string) (Otp, error) {

	otp := New(email)
	otp.Name = name
	err := otp.Encrypt(s.encryptor)

	if err != nil {
		return *otp, errs.InternalServerError(err).Log(s.logger)
	}
	defer s.otpRepo.InvalidateOtherOTPs(otp.Email, otp.Code, otp.RequestId)
	err = s.otpRepo.Add(otp)
	if err != nil {
		return *otp, err
	}
	otp.Email = email

	if err = sendOTP(s.emailer, otp.Email, otp.Code, name); err != nil {
		return *otp, errs.InternalServerError(err).Log(s.logger)
	}

	return *otp, nil
}

// VerifyOTP checks if the OTP provided by the user is valid.
//
// It returns a boolean indicating if the OTP is valid and an error if any.
// If the OTP is valid, the function also checks if the OTP has expired.
//
// If the OTP has expired or is invalid, it returns false.
func (s *otpService) VerifyOTP(email, code, requestId string) (Otp, error) {

	otp := New(email)
	err := otp.Encrypt(s.encryptor)

	if err != nil {
		return *otp, errs.InternalServerError(err).Log(s.logger)
	}
	otp, err = s.otpRepo.GetByEmailAndReference(otp.Email, requestId)
	if err != nil {
		return Otp{}, errs.New("Invalid OTP", http.StatusNotFound)
	}

	if err = otp.Decrypt(s.encryptor, email); err != nil {
		return *otp, errs.InternalServerError(err).Log(s.logger)
	}

	if otp.IsExpired() {
		return *otp, errs.New("OTP has expired", http.StatusBadRequest)
	}

	if otp.Code != code || otp.RequestId != requestId {
		return *otp, errs.New("Invalid OTP", http.StatusBadRequest)
	}
	defer s.otpRepo.Delete(otp)
	return *otp, nil
}

// Sends OTP to the user via email
func sendOTP(emailer email.Emailer, email, code, name string) error {
	verificationTemplate := emailTemplates.Verification([]string{email}, name, code)
	return emailer.SendEmailTemplate(*verificationTemplate)
}
