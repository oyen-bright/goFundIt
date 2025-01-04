package paystack

// PaystackWebhookEvent represents the structure of a Paystack webhook event
const (
	EventChargeSuccess = "charge.success"
	EventChargeFailed  = "charge.failed"
)

type PaystackWebhookEvent struct {
	Event string `json:"event"`
	Data  struct {
		ID        int     `json:"id"`
		Reference string  `json:"reference"`
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		Channel   string  `json:"channel"`
		Status    string  `json:"status"`
		PaidAt    string  `json:"paid_at"`
		CreatedAt string  `json:"created_at"`
		Customer  struct {
			Email string `json:"email"`
			Name  string `json:"customer_code"`
		} `json:"customer"`
		Metadata struct {
			ContributorID string `json:"contributor_id"`
			CampaignID    string `json:"campaign_id"`
		} `json:"metadata"`
	} `json:"data"`
}
