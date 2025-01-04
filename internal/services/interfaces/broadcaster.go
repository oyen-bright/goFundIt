package interfaces

import "github.com/oyen-bright/goFundIt/pkg/websocket"

type EventBroadcaster interface {
	NewEvent(campaignID string, eventType websocket.EventType, data interface{})
}
