package activity

import contributor "github.com/oyen-bright/goFundIt/internal/contributor"

type ActivityService interface {
	CreateActivity(activity *Activity) error
	GetActivitiesByCampaignID(campaignID string) ([]Activity, error)
	GetActivityByID(activityID uint) (Activity, error)
	UpdateActivity(activity *Activity) error
	DeleteActivity(activityID uint) error
	JoinActivity(contributorID uint, activityID uint) error
	LeaveActivity(contributorID uint, activityID uint) error
	GetActivityParticipants(activityID uint) ([]contributor.Contributor, error)
}
