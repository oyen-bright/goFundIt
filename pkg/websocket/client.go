package websocket

import "github.com/gorilla/websocket"

type Client struct {
	Hub        *Hub
	conn       *websocket.Conn
	send       chan Message
	campaignID string
	userHandle string
}

func NewClient(hub *Hub, conn *websocket.Conn, campaignID, userHandle string) *Client {
	return &Client{
		Hub:        hub,
		conn:       conn,
		send:       make(chan Message),
		campaignID: campaignID,
		userHandle: userHandle,
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		}
	}
}
