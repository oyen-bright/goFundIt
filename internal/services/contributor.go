package services

import (
	"github.com/oyen-bright/goFundIt/internal/models"
	repositories "github.com/oyen-bright/goFundIt/internal/repositories/interfaces"
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type contributorService struct {
	repo   repositories.ContributorRepository
	logger logger.Logger
}

func NewContributorService(repo repositories.ContributorRepository, logger logger.Logger) services.ContributorService {
	return &contributorService{repo: repo, logger: logger}
}

func (s *contributorService) AddContributorToCampaign(contribution *models.Contributor) error {
	return nil
}

func (s *contributorService) RemoveContributorFromCampaign(contribution *models.Contributor) error {
	return nil
}

func (s *contributorService) CanContributeToCampaign(userID uint, campaignID string) (bool, error) {
	return false, nil
}

func (s *contributorService) GetContributors(campaignID string) ([]models.Contributor, error) {
	return []models.Contributor{}, nil
}

func (s *contributorService) GetContributorByID(contributorID uint) (models.Contributor, error) {
	return models.Contributor{}, nil
}

func (s *contributorService) GetContributorByUserHandle(userHandle uint) (models.Contributor, error) {
	return models.Contributor{}, nil
}

func (s *contributorService) GetEmailsOfExistingContributors(emails []string) ([]string, error) {
	existingContributorsEmail, err := s.repo.GetEmailsOfExistingContributors(emails)

	if err != nil {
		return existingContributorsEmail, errs.InternalServerError(err).Log(s.logger)
	}
	return existingContributorsEmail, nil
}

func (s *contributorService) ContributeToCampaign(userID uint, campaignID string, amount float64) error {
	return nil
}

func (s *contributorService) ProcessPayment(paymentID string) error {
	return nil
}

func (s *contributorService) RefundPayment(paymentID string) error {
	return nil
}
