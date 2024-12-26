package services

import (
	services "github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
)

type eventBroadcasterImpl struct {
	hub *websocket.Hub
}

func NewEventBroadcaster(hub *websocket.Hub) services.EventBroadcaster {
	return &eventBroadcasterImpl{
		hub: hub,
	}
}

func (e *eventBroadcasterImpl) NewEvent(campaignID string, eventType websocket.EventType, data interface{}) {
	e.hub.BroadcastToCampaign(campaignID, websocket.Message{
		Type: eventType,
		Data: data,
	})
}
