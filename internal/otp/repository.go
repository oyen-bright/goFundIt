package otp

import (
	"gorm.io/gorm"
)

type OTPRepository interface {
	GetByEmailAndReference(email, requestId string) (*Otp, error)
	Add(otp *Otp) error
	InvalidateOtherOTPs(email, code, requestId string) error
	Delete(otp *Otp) error
}

type otpRepository struct {
	db *gorm.DB
}

func Repository(db *gorm.DB) OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) GetByEmailAndReference(email, requestId string) (*Otp, error) {
	var otp Otp
	if err := r.db.Where("email = ? AND request_id = ?", email, requestId).First(&otp).Error; err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) Add(otp *Otp) error {
	return r.db.Create(otp).Error
}

func (r *otpRepository) Delete(otp *Otp) error {
	return r.db.Delete(otp).Error
}

func (r *otpRepository) InvalidateOtherOTPs(email, code, requestId string) error {
	return r.db.Where("email = ? AND code != ? AND request_id != ?", email, code, requestId).Delete(&Otp{}).Error
}
