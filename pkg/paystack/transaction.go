package paystack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/oyen-bright/goFundIt/pkg/utils"
)

// InitiateTransaction initiates a transaction on Paystack
func (c *client) InitiateTransaction(email, currency string, amount float64) (*TransactionResponse, error) {
	txn := NewTransaction(email, currency, amount)
	reqBody, err := txn.GetBody()
	if err != nil {
		return nil, err
	}
	resp, err := c.SetupRequest(http.MethodPost, "/transaction/initialize", reqBody, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var txnResp TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&txnResp); err != nil {
		return nil, err
	}
	return &txnResp, nil

}

func (c *client) VerifyTransaction(reference string) (*VerifyTransactionResponse, error) {
	resp, err := c.SetupRequest(http.MethodGet, "/transaction/verify/"+reference, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var txnResp VerifyTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&txnResp); err != nil {
		return nil, err
	}
	return &txnResp, nil
}

// Models

// InitiateTransaction represents the request body for initiating a transaction
type Transaction struct {
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
	Currency  string  `json:"currency"`
}

// NewTransaction creates a new transaction instance with the provided parameters
func NewTransaction(email, currency string, amount float64) *Transaction {
	return &Transaction{

		Email:     email,
		Amount:    amount * 100,
		Reference: generateReference(),
		Currency:  currency,
	}
}

// Get Body returns the body of the transaction
func (t *Transaction) GetBody() (io.Reader, error) {
	body, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

// TransactionResponse represents the response from the Paystack API when a transaction is initiated
type TransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	}
}

// VerifyTransactionResponse represents the response from the Paystack API when a transaction is verified
type VerifyTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              int    `json:"id"`
		Status          string `json:"status"`
		Message         string `json:"message"`
		GatewayResponse string `json:"gateway_response"`
		Fees            int    `json:"fees"`
		Reference       string `json:"reference"`
		Amount          int    `json:"amount"`
		Currency        string `json:"currency"`
		PaidAt          string `json:"paid_at"`
		CreatedAt       string `json:"created_at"`
	}
}

func (v *VerifyTransactionResponse) IsPaymentSuccessful() bool {
	return v.Data.Status == "success" && v.Data.GatewayResponse == "Successful"
}

func (v *VerifyTransactionResponse) ToString() string {
	byte, err := json.Marshal(v.Data)
	if err != nil {
		return ""

	}
	value := string(byte)
	return value
}

// Helper Methods
func generateReference() string {
	return utils.GenerateRandomAlphaNumeric("GF-", 13)
}
