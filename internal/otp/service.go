package otp

import (
	"github.com/oyen-bright/goFundIt/pkg/email"
	emailTemplates "github.com/oyen-bright/goFundIt/pkg/email/templates"
	encryptor "github.com/oyen-bright/goFundIt/pkg/encryption"
	"gorm.io/gorm"
)

type OTpServiceInterface interface {
	RequestOTP(email, name string) error
	VerifyOTP(email, otp string) (bool, error)
}

type OtpService struct {
	DB        *gorm.DB
	Emailer   *email.Emailer
	Encryptor *encryptor.Encryptor
}

// RequestOTP generates a new OTP and sends it to the user.
//
//  1. Clears any previous OTPs for the user from the database.
//  2. Inserts a new OTP into the database.
//  3. Sends the OTP to the user via email.
//
// It returns an error if any of the steps fail.
func (s *OtpService) RequestOTP(email, name string) error {

	otp := New(email)
	err := otp.Encrypt(*s.Encryptor)

	if err != nil {
		return err
	}
	defer clearUserOTPs(s.DB, otp.Email, otp.Code)
	err = insertOTP(s.DB, otp)
	if err != nil {
		return err
	}
	otp.Email = email
	return sendOTP(*s.Emailer, otp, name)
}

// VerifyOTP checks if the OTP provided by the user is valid.
//
// It returns a boolean indicating if the OTP is valid and an error if any.
// If the OTP is valid, the function also checks if the OTP has expired.
//
// If the OTP has expired or is invalid, it returns false.
func (s *OtpService) VerifyOTP(email, code string) (bool, error) {

	otp := New(email)
	err := otp.Encrypt(*s.Encryptor)

	if err != nil {
		return false, err
	}
	if err := s.DB.Where("email = ?", otp.Email).First(&otp).Error; err != nil {
		return false, err
	}
	otp.Decrypt(*s.Encryptor)
	if otp.IsExpired() {
		return false, nil
	}

	if otp.Code != code {
		return false, nil
	}
	return true, nil
}

// Clears users previous OTPs from the database
func clearUserOTPs(dB *gorm.DB, email, code string) error {
	return dB.Where("email = ? AND code != ?", email, code).Delete(&Otp{}).Error
}

// Inserts OTP into the database
func insertOTP(dB *gorm.DB, otp Otp) error {
	return dB.Create(&otp).Error
}

// Sends OTP to the user via email
func sendOTP(emailer email.Emailer, otp Otp, name string) error {
	verificationTemplate := emailTemplates.Verification([]string{otp.Email}, name, otp.Code)
	return emailer.SendEmailTemplate(*verificationTemplate)
}
