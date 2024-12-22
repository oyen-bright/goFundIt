package postgress

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	"github.com/oyen-bright/goFundIt/internal/repositories/interfaces"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) interfaces.PaymentRepository {
	return &paymentRepository{db: db}
}
func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetByReference(reference string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Preload("Contributor").First(&payment, "reference = ?", reference).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Model(&models.Payment{}).Where("reference = ?", payment.Reference).Updates(
		&models.Payment{PaymentStatus: payment.PaymentStatus, GatewayResponse: payment.GatewayResponse}).Error
}

func (r *paymentRepository) Delete(reference string) error {
	return r.db.Delete(&models.Payment{}, "reference = ?", reference).Error
}

func (r *paymentRepository) GetByCampaign(campaignID string, limit, offset int) ([]*models.Payment, int64, error) {
	var payments []*models.Payment
	var total int64

	if err := r.db.Model(&models.Payment{}).Where("campaign_id = ?", campaignID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("campaign_id = ?", campaignID).Limit(limit).Offset(offset).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *paymentRepository) GetByContributor(contributorID uint, limit, offset int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	if err := r.db.Model(&models.Payment{}).Where("contributor_id = ?", contributorID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("contributor_id = ?", contributorID).Limit(limit).Offset(offset).Find(&payments).Error; err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}
