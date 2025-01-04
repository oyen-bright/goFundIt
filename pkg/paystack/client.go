package paystack

import (
	"io"
	"net/http"
)

// PaystackClient represents all available Paystack operations
type PaystackClient interface {
	InitiateTransaction(email, currency string, amount float64) (*TransactionResponse, error)
	VerifyTransaction(reference string) (*VerifyTransactionResponse, error)
	CreateRecipient(recipient Recipient) (*RecipientResponse, error)
	InitiateTransfer(transfer Transfer) (*TransactionResponse, error)
	FinalizeTransfer(transferCode string) (*FinalizeTransferResponse, error)
	ResolveAccount(accountNumber, bankCode string) (*ResolveAccountResponse, error)
	GetBanks() (*BankListResponse, error)
}

type client struct {
	secretKey string
	baseURL   string
}

func NewClient(secretKey string) PaystackClient {
	return &client{
		secretKey: secretKey,
		baseURL:   "https://api.paystack.co",
	}
}

// Helper factions

// SetupRequest sets up the request to the Paystack API
func (c *client) SetupRequest(method, path string, body io.Reader, queryParam *[]map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	if queryParam != nil {
		q := req.URL.Query()
		for _, v := range *queryParam {
			for key, value := range v {
				q.Add(key, value)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Add("Authorization", "Bearer "+c.secretKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
