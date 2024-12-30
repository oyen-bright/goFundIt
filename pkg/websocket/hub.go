package websocket

import (
	"sync"
)

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Broadcast(message Message) {
	h.broadcast <- message
}

// Add new method to broadcast to specific campaign
// func (h *Hub) BroadcastToCampaign(campaignID string, message Message) {
// 	h.mutex.RLock()
// 	if clients, ok := h.clients[campaignID]; ok {
// 		for client := range clients {
// 			select {
// 			case client.send <- message:
// 			default:
// 				close(client.send)
// 				delete(clients, client)
// 			}
// 		}
// 	}
// 	h.mutex.RUnlock()
// }

func (h *Hub) BroadcastToCampaign(campaignID string, message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, ok := h.clients[campaignID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				h.mutex.RUnlock()
				h.unregister <- client
				h.mutex.RLock()
			}
		}
	}
}

// Modify the Run method to handle campaign-specific broadcasts
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if client == nil {
				continue
			}

			h.mutex.Lock()
			if _, ok := h.clients[client.campaignID]; !ok {
				h.clients[client.campaignID] = make(map[*Client]bool)
			}
			h.clients[client.campaignID][client] = true
			h.mutex.Unlock()

		case client := <-h.unregister:

			if client == nil {
				continue
			}
			h.mutex.Lock()
			if _, ok := h.clients[client.campaignID]; ok {
				delete(h.clients[client.campaignID], client)
				close(client.send)
				if len(h.clients[client.campaignID]) == 0 {
					delete(h.clients, client.campaignID)
				}
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			// Remove the old broadcast code and use BroadcastToCampaign instead
			if campaignID, ok := message.Data.(string); ok {
				h.BroadcastToCampaign(campaignID, message)
			}
		}
	}
}

func (h *Hub) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for campaignID, clients := range h.clients {
		for client := range clients {
			close(client.send)
			delete(clients, client)
		}
		delete(h.clients, campaignID)
	}

	close(h.broadcast)
	close(h.register)
	close(h.unregister)
}
