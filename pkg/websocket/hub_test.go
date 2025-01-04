package websocket

import (
	"testing"
	"time"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub.clients == nil {
		t.Error("expected clients map to be initialized")
	}
	if hub.broadcast == nil {
		t.Error("expected broadcast channel to be initialized")
	}
	if hub.register == nil {
		t.Error("expected register channel to be initialized")
	}
	if hub.unregister == nil {
		t.Error("expected unregister channel to be initialized")
	}
}

func TestHubRegisterAndUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Close()

	// Create a test client
	client := &Client{
		Hub:        hub,
		send:       make(chan Message),
		campaignID: "test-campaign",
	}

	// Test registration
	hub.Register(client)
	time.Sleep(100 * time.Millisecond) // Give time for the registration to process

	hub.mutex.RLock()
	if _, ok := hub.clients[client.campaignID]; !ok {
		t.Error("expected client to be registered")
	}
	hub.mutex.RUnlock()

	// Test unregistration
	hub.unregister <- client
	time.Sleep(100 * time.Millisecond) // Give time for the unregistration to process

	hub.mutex.RLock()
	if _, ok := hub.clients[client.campaignID]; ok {
		if len(hub.clients[client.campaignID]) > 0 {
			t.Error("expected client to be unregistered")
		}
	}
	hub.mutex.RUnlock()
}

func TestBroadcastToCampaign(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Close()

	campaignID := "test-campaign"
	messageType := EventTypeActivityCreated
	messageData := "test message"

	testMessage := Message{
		Type: messageType,
		Data: messageData,
	}

	client := &Client{
		Hub:        hub,
		send:       make(chan Message, 1), // Add buffer to prevent blocking
		campaignID: campaignID,
	}

	// Register client and wait for registration to complete
	hub.Register(client)
	time.Sleep(100 * time.Millisecond)

	// Broadcast in a separate goroutine
	go hub.BroadcastToCampaign(campaignID, testMessage)

	select {
	case receivedMsg := <-client.send:
		if receivedMsg.Type != messageType {
			t.Errorf("expected message type %q, got %q", messageType, receivedMsg.Type)
		}
		if receivedMsg.Data != messageData {
			t.Errorf("expected message data %q, got %v", messageData, receivedMsg.Data)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for broadcast message")
	}
}

func TestClose(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Register some test clients
	client1 := &Client{
		Hub:        hub,
		send:       make(chan Message),
		campaignID: "campaign1",
	}
	client2 := &Client{
		Hub:        hub,
		send:       make(chan Message),
		campaignID: "campaign2",
	}

	hub.Register(client1)
	hub.Register(client2)
	time.Sleep(100 * time.Millisecond)

	// Test closing
	hub.Close()

	// Verify channels are closed
	select {
	case _, ok := <-hub.broadcast:
		if ok {
			t.Error("broadcast channel should be closed")
		}
	default:
	}

	select {
	case _, ok := <-hub.register:
		if ok {
			t.Error("register channel should be closed")
		}
	default:
	}

	select {
	case _, ok := <-hub.unregister:
		if ok {
			t.Error("unregister channel should be closed")
		}
	default:
	}

	// Verify clients map is empty
	hub.mutex.RLock()
	if len(hub.clients) != 0 {
		t.Error("clients map should be empty")
	}
	hub.mutex.RUnlock()
}
