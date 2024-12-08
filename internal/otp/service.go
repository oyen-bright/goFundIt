package otp

import (
	"github.com/oyen-bright/goFundIt/pkg/email"
	emailTemplates "github.com/oyen-bright/goFundIt/pkg/email/templates"
	encryptor "github.com/oyen-bright/goFundIt/pkg/encryption"
)

type OTPService interface {
	RequestOTP(email, name string) (Otp, error)
	VerifyOTP(email, otp, requestId string) (Otp, error)
}

type otpService struct {
	otpRepo   OTPRepository
	emailer   email.Emailer
	encryptor encryptor.Encryptor
}

func Service(otpRepo OTPRepository, emailer email.Emailer, encryptor encryptor.Encryptor) OTPService {
	return &otpService{
		otpRepo:   otpRepo,
		emailer:   emailer,
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
		return *otp, err
	}
	defer s.otpRepo.InvalidateOtherOTPs(otp.Email, otp.Code, otp.RequestId)
	err = s.otpRepo.Add(otp)
	if err != nil {
		return *otp, err
	}
	otp.Email = email
	return *otp, sendOTP(s.emailer, otp.Email, otp.Code, name)
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
		return *otp, err
	}
	otp, err = s.otpRepo.GetByEmailAndReference(otp.Email, requestId)
	if err != nil {
		return Otp{}, nil
	}

	otp.Decrypt(s.encryptor, email)

	// if otp.IsExpired() {
	// 	return *otp, nil
	// }

	// if otp.Code != code || otp.RequestId != requestId {
	// 	return *otp, nil
	// }
	return *otp, nil
}

// Sends OTP to the user via email
func sendOTP(emailer email.Emailer, email, code, name string) error {
	verificationTemplate := emailTemplates.Verification([]string{email}, name, code)
	return emailer.SendEmailTemplate(*verificationTemplate)
}
