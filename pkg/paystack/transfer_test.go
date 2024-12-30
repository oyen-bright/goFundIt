package paystack

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitiateTransfer(t *testing.T) {
	tests := []struct {
		name          string
		mockRequest   *Transfer
		mockResponse  *TransactionResponse
		expectedError bool
	}{
		{
			name: "successful transfer",
			mockRequest: &Transfer{
				Source:    "balance",
				Reason:    "test transfer",
				Amount:    1000.00,
				Currency:  "NGN",
				Recipient: "test_recipient",
			},
			mockResponse: &TransactionResponse{
				Status:  true,
				Message: "Transfer initiated",
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
				if r.URL.Path != "/transfer" {
					t.Errorf("Expected path /transfer, got %s", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			client := &Client{
				secretKey: "test_key",
				baseURL:   server.URL,
			}
			resp, err := client.InitiateTransfer(*tt.mockRequest)
			if (err != nil) != tt.expectedError {
				t.Errorf("InitiateTransfer() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if err == nil && resp.Status != tt.mockResponse.Status {
				t.Errorf("InitiateTransfer() status = %v, want %v", resp.Status, tt.mockResponse.Status)
			}
		})

	}
}

func TestFinalizeTransfer(t *testing.T) {
	tests := []struct {
		name          string
		transferCode  string
		mockResponse  *FinalizeTransferResponse
		expectedError bool
	}{
		{
			name:         "successful finalize transfer",
			transferCode: "test_transfer_code",
			mockResponse: &FinalizeTransferResponse{
				Status:  true,
				Message: "Transfer finalized",
				Data: struct {
					Reference string `json:"reference"`
					Status    string `json:"status"`
					Failures  string `json:"failures"`
				}{
					Reference: "test_ref",
					Status:    "success",
					Failures:  "",
				},
			},
			expectedError: false,
		},
		{
			name:         "failed finalize transfer",
			transferCode: "invalid_code",
			mockResponse: &FinalizeTransferResponse{
				Status:  false,
				Message: "Transfer failed",
				Data: struct {
					Reference string `json:"reference"`
					Status    string `json:"status"`
					Failures  string `json:"failures"`
				}{
					Reference: "",
					Status:    "failed",
					Failures:  "Invalid transfer code",
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/transfer/finalize_transfer" {
					t.Errorf("Expected path /transfer/finalize_transfer, got %s", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			client := &Client{
				secretKey: "test_key",
				baseURL:   server.URL + "/",
			}

			resp, err := client.FinalizeTransfer(tt.transferCode)
			log.Println(err)
			if (err != nil) != tt.expectedError {
				t.Errorf("FinalizeTransfer() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if err == nil {
				if resp.Status != tt.mockResponse.Status {
					t.Errorf("FinalizeTransfer() status = %v, want %v", resp.Status, tt.mockResponse.Status)
				}
				if resp.Data.Status != tt.mockResponse.Data.Status {
					t.Errorf("FinalizeTransfer() transfer status = %v, want %v", resp.Data.Status, tt.mockResponse.Data.Status)
				}
			}
		})
	}
}
