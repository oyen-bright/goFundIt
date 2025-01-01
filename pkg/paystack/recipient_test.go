package paystack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRecipient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/transferrecipient" {
			t.Errorf("Expected path '/transferrecipient', got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Read and verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var recipient Recipient
		if err := json.Unmarshal(body, &recipient); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		// Verify required fields
		if recipient.Type != "nuban" {
			t.Errorf("Expected type 'nuban', got %s", recipient.Type)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": true,
			"message": "Recipient created",
			"data": {
				"recipient_code": "RCP_1234567890"
			}
		}`))
	}))
	defer server.Close()

	testClient := &client{ // Use lowercase client
		baseURL:   server.URL,
		secretKey: "test_key",
	}

	recipient := NewRecipient("John Doe", "1234567890", "123", "NGN")
	resp, err := testClient.CreateRecipient(*recipient)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !resp.Status {
		t.Error("Expected status to be true")
	}
	if resp.Data.RecipientCode != "RCP_1234567890" {
		t.Errorf("Expected recipient code 'RCP_1234567890', got %s", resp.Data.RecipientCode)
	}
}
