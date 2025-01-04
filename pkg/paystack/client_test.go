package paystack

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	secretKey := "test_secret_key"
	client := NewClient(secretKey).(*client) // Type assert to access private fields

	if client.secretKey != secretKey {
		t.Errorf("Expected secret key %s, got %s", secretKey, client.secretKey)
	}

	if client.baseURL != "https://api.paystack.co" {
		t.Errorf("Expected base URL %s, got %s", "https://api.paystack.co", client.baseURL)
	}
}

func TestSetupRequest(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		queryParam *[]map[string]string
		wantErr    bool
	}{
		{
			name:    "GET request without query params",
			method:  "GET",
			path:    "/test",
			body:    "",
			wantErr: false,
		},
		{
			name:   "POST request with query params",
			method: "POST",
			path:   "/test",
			body:   `{"test": "data"}`,
			queryParam: &[]map[string]string{
				{"key": "value"},
			},
			wantErr: false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	testClient := &client{ // Use lowercase client
		secretKey: "test_secret_key",
		baseURL:   server.URL,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.body)
			resp, err := testClient.SetupRequest(tt.method, tt.path, bodyReader, tt.queryParam)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetupRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Error("Expected response, got nil")
				return
			}

			if resp != nil {
				if resp.Request.Method != tt.method {
					t.Errorf("Expected method %s, got %s", tt.method, resp.Request.Method)
				}

				auth := resp.Request.Header.Get("Authorization")
				if auth != "Bearer "+testClient.secretKey {
					t.Errorf("Expected Authorization header Bearer %s, got %s", testClient.secretKey, auth)
				}

				contentType := resp.Request.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type header application/json, got %s", contentType)
				}
			}
		})
	}
}
