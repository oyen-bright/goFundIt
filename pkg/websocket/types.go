package websocket

type EventType string

const (
	EventTypeActivityCreated     EventType = "activity_created"
	EventTypeActivityUpdated     EventType = "activity_updated"
	EventTypeActivityDeleted     EventType = "activity_deleted"
	EventTypeCommentCreated      EventType = "comment_created"
	EventTypeCommentDeleted      EventType = "comment_deleted"
	EventTypeContributionCreated EventType = "contribution_created"
	EventTypeContributorUpdated  EventType = "contributor_updated"
	EventTypeContributorDeleted  EventType = "contributor_deleted"
	EventTypeCampaignUpdated     EventType = "campaign_updated"
	EventTypePayoutCreated       EventType = "payout_created"
	EventTypePayoutUpdated       EventType = "payout_updated"
)

type Message struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}
