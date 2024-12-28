package paystack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) CreateRecipient(recipient Recipient) (*RecipientResponse, error) {

	body, err := recipient.getBody()
	if err != nil {
		return nil, err
	}
	resp, err := c.SetupRequest(http.MethodPost, "/transferrecipient", body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var recResp RecipientResponse
	if err := json.NewDecoder(resp.Body).Decode(&recResp); err != nil {
		return nil, err
	}
	return &recResp, nil

}

// Models

// Recipient represents the structure of a recipient
type Recipient struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency"`
}

// NewRecipient creates a new recipient instance with the provided parameters
func NewRecipient(name, accountNumber, bankCode, currency string) *Recipient {
	return &Recipient{
		Type:          "nuban",
		Name:          name,
		AccountNumber: accountNumber,
		BankCode:      bankCode,
		Currency:      currency,
	}
}

func (r *Recipient) getBody() (io.Reader, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(data), nil
}

// RecipientResponse represents the response from the Paystack API when a recipient is created
type RecipientResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		RecipientCode string `json:"recipient_code"`
	} `json:"data"`
}
