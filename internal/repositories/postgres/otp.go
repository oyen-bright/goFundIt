package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) interfaces.OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) GetByEmailAndReference(email, requestId string) (*models.Otp, error) {
	var otp models.Otp
	if err := r.db.Where("email = ? AND request_id = ?", email, requestId).First(&otp).Error; err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) Add(otp *models.Otp) error {
	return r.db.Create(otp).Error
}

func (r *otpRepository) Delete(otp *models.Otp) error {
	return r.db.Delete(otp).Error
}

func (r *otpRepository) InvalidateOtherOTPs(email, code, requestId string) error {
	return r.db.Where("email = ? AND code != ? AND request_id != ?", email, code, requestId).Delete(&models.Otp{}).Error
}
