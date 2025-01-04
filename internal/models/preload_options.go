package models

type PreloadOption struct {
	Images                 bool
	Payout                 bool
	Activities             bool
	ActivitiesContributors bool
	ActivitiesComments     bool
	Contributors           bool
	ContributorsActivities bool
	CreatedBy              bool
}
