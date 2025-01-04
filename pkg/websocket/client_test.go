package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func TestNewClient(t *testing.T) {
	hub := NewHub()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("failed to upgrade connection: %v", err)
		}
		client := NewClient(hub, conn, "campaign1", "user1")

		if client.Hub != hub {
			t.Errorf("expected hub to be %v, got %v", hub, client.Hub)
		}
		if client.campaignID != "campaign1" {
			t.Errorf("expected campaignID to be campaign1, got %v", client.campaignID)
		}
		if client.userHandle != "user1" {
			t.Errorf("expected userHandle to be user1, got %v", client.userHandle)
		}
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the test server
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
}
