package paystack

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitiateTransaction(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		currency      string
		amount        float64
		mockResponse  *TransactionResponse
		expectedError bool
	}{
		{
			name:     "successful transaction",
			email:    "test@example.com",
			currency: "NGN",
			amount:   1000.00,
			mockResponse: &TransactionResponse{
				Status:  true,
				Message: "Authorization URL created",
				Data: struct {
					AuthorizationURL string `json:"authorization_url"`
					AccessCode       string `json:"access_code"`
					Reference        string `json:"reference"`
				}{
					AuthorizationURL: "https://checkout.paystack.com/test",
					AccessCode:       "test_code",
					Reference:        "test_ref",
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/transaction/initialize" {
					t.Errorf("Expected path /transaction/initialize, got %s", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			testClient := &client{
				secretKey: "test_key",
				baseURL:   server.URL,
			}

			resp, err := testClient.InitiateTransaction(tt.email, tt.currency, tt.amount)
			if (err != nil) != tt.expectedError {
				t.Errorf("InitiateTransaction() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if err == nil && resp.Status != tt.mockResponse.Status {
				t.Errorf("InitiateTransaction() status = %v, want %v", resp.Status, tt.mockResponse.Status)
			}
		})
	}
}

func TestVerifyTransaction(t *testing.T) {
	tests := []struct {
		name          string
		reference     string
		mockResponse  *VerifyTransactionResponse
		expectedError bool
	}{
		{
			name:      "successful verification",
			reference: "test_ref",
			mockResponse: &VerifyTransactionResponse{
				Status:  true,
				Message: "Verification successful",
				Data: struct {
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
				}{
					Status:          "success",
					GatewayResponse: "Successful",
					Reference:       "test_ref",
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/transaction/verify/" + tt.reference
				if r.URL.Path != expectedPath {
					t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			testClient := &client{ // Use lowercase client
				secretKey: "test_key",
				baseURL:   server.URL,
			}

			resp, err := testClient.VerifyTransaction(tt.reference)
			if (err != nil) != tt.expectedError {
				t.Errorf("VerifyTransaction() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if err == nil {
				if !resp.IsPaymentSuccessful() {
					t.Error("Expected payment to be successful")
				}
				if resp.Status != tt.mockResponse.Status {
					t.Errorf("VerifyTransaction() status = %v, want %v", resp.Status, tt.mockResponse.Status)
				}
			}
		})
	}
}

func TestTransaction_GetBody(t *testing.T) {
	txn := NewTransaction("test@example.com", "NGN", 1000.00)
	body, err := txn.GetBody()

	if err != nil {
		t.Errorf("GetBody() error = %v", err)
		return
	}

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("Failed to read body: %v", err)
		return
	}

	var decodedTxn Transaction
	if err := json.Unmarshal(bodyBytes, &decodedTxn); err != nil {
		t.Errorf("Failed to unmarshal body: %v", err)
		return
	}

	if decodedTxn.Email != txn.Email {
		t.Errorf("Email = %v, want %v", decodedTxn.Email, txn.Email)
	}
	if decodedTxn.Amount != txn.Amount {
		t.Errorf("Amount = %v, want %v", decodedTxn.Amount, txn.Amount)
	}
	if decodedTxn.Currency != txn.Currency {
		t.Errorf("Currency = %v, want %v", decodedTxn.Currency, txn.Currency)
	}
}
