package paystack

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBanks(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bank" {
			t.Errorf("Expected path '/bank', got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": true,
			"message": "Banks retrieved",
			"data": [{
				"name": "Test Bank",
				"code": "123",
				"country": "Nigeria",
				"currency": "NGN"
			}]
		}`))
	}))
	defer server.Close()

	testClient := &client{
		baseURL:   server.URL,
		secretKey: "test_key",
	}

	resp, err := testClient.GetBanks()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !resp.Status {
		t.Error("Expected status to be true")
	}
	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 bank, got %d", len(resp.Data))
	}
}

func TestResolveAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bank/resolve" {
			t.Errorf("Expected path '/bank/resolve', got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": true,
			"message": "Account resolved",
			"data": {
				"account_number": "1234567890",
				"account_name": "John Doe"
			}
		}`))
	}))
	defer server.Close()

	testClient := &client{
		baseURL:   server.URL,
		secretKey: "test_key",
	}

	resp, err := testClient.ResolveAccount("1234567890", "123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !resp.Status {
		t.Error("Expected status to be true")
	}
	if resp.Data.AccountName != "John Doe" {
		t.Errorf("Expected account name 'John Doe', got %s", resp.Data.AccountName)
	}
}
