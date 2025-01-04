package paystack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// InitiateTransfer
func (c *client) InitiateTransfer(transfer Transfer) (*TransactionResponse, error) {
	body, err := transfer.GetBody()
	if err != nil {
		return nil, err
	}

	resp, err := c.SetupRequest(http.MethodPost, "/transfer", body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var txnRes TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&txnRes); err != nil {
		return nil, err
	}
	return &txnRes, nil
}

// FinalizeTransfer finalize the initiated transfer by the transfer code
func (c *client) FinalizeTransfer(transferCode string) (*FinalizeTransferResponse, error) {

	data, err := json.Marshal(map[string]string{
		"transfer_code": transferCode,
	})
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(data)

	resp, err := c.SetupRequest(http.MethodPost, "transfer/finalize_transfer", body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var txnRes FinalizeTransferResponse
	if err := json.NewDecoder(resp.Body).Decode(&txnRes); err != nil {
		return nil, err
	}
	return &txnRes, nil
}

type Transfer struct {
	Source    string  `json:"source"`
	Reason    string  `json:"reason"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Recipient string  `json:"recipient"`
}

// NewTransfer
func NewTransfer(reason, recipient, currency string, amount float64) *Transfer {
	return &Transfer{
		Source:    "balance",
		Reason:    reason,
		Amount:    amount,
		Recipient: recipient,
		Currency:  currency,
	}
}

// Transfer Methods
func (t *Transfer) GetBody() (io.Reader, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(data), nil
}

type TransferResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TransferCode string `json:"transfer_code"`
		ID           string `json:"id"`
	} `json:"data"`
}

type FinalizeTransferResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference string `json:"reference"`
		Status    string `json:"status"`
		Failures  string `json:"failures"`
	} `json:"data"`
}
