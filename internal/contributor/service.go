package contributor

import (
	"github.com/oyen-bright/goFundIt/pkg/errs"
	"github.com/oyen-bright/goFundIt/pkg/logger"
)

type ContributorService interface {
	AddContributorToCampaign(contribution *Contributor) error
	RemoveContributorFromCampaign(contribution *Contributor) error
	CanContributeToCampaign(userID uint, campaignID string) (bool, error)
	GetContributors(campaignID string) ([]Contributor, error)
	GetContributorByID(contributorID uint) (Contributor, error)
	GetContributorByUserHandle(userHandle uint) (Contributor, error)
	GetEmailsOfExistingContributors(emails []string) ([]string, error)
	ContributeToCampaign(userID uint, campaignID string, amount float64) error
	ProcessPayment(paymentID string) error
	RefundPayment(paymentID string) error
}
type contributorService struct {
	repo   ContributorRepository
	logger logger.Logger
}

func Service(repo ContributorRepository, logger logger.Logger) ContributorService {
	return &contributorService{repo: repo, logger: logger}
}

func (s *contributorService) AddContributorToCampaign(contribution *Contributor) error {
	return nil
}

func (s *contributorService) RemoveContributorFromCampaign(contribution *Contributor) error {
	return nil
}

func (s *contributorService) CanContributeToCampaign(userID uint, campaignID string) (bool, error) {
	return false, nil
}

func (s *contributorService) GetContributors(campaignID string) ([]Contributor, error) {
	return []Contributor{}, nil
}

func (s *contributorService) GetContributorByID(contributorID uint) (Contributor, error) {
	return Contributor{}, nil
}

func (s *contributorService) GetContributorByUserHandle(userHandle uint) (Contributor, error) {
	return Contributor{}, nil
}

func (s *contributorService) GetEmailsOfExistingContributors(emails []string) ([]string, error) {
	existingContributorsEmail, err := s.repo.getEmailsOfExistingContributors(emails)

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
