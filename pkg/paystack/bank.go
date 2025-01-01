package paystack

import (
	"encoding/json"
	"net/http"
)

// GetBanks returns a list of banks
func (c *client) GetBanks() (*BankListResponse, error) {
	resp, err := c.SetupRequest(http.MethodGet, "/bank", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var bankList BankListResponse

	if err := json.NewDecoder(resp.Body).Decode(&bankList); err != nil {
		return nil, err
	}
	return &bankList, nil
}

// ResolveAccount resolves an account number
func (c *client) ResolveAccount(accountNumber, bankCode string) (*ResolveAccountResponse, error) {

	params := []map[string]string{
		{
			"account_number": accountNumber,
		},
		{
			"bank_code": bankCode,
		},
	}

	resp, err := c.SetupRequest(http.MethodGet, "/bank/resolve", nil, &params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var resolveAccountResponse ResolveAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&resolveAccountResponse); err != nil {
		return nil, err
	}
	return &resolveAccountResponse, nil
}

// Models

// Bank represents the structure of a bank
type Bank struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Country  string `json:"country"`
	Currency string `json:"currency"`
}

// BankListResponse represents the response from the Paystack API when a list of banks is requested
type BankListResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []Bank `json:"data"`
}

// ResolveAccount

// ResolveAccountResponse represents the response from the Paystack API when an account number is resolved
type ResolveAccountResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
	} `json:"data"`
}
