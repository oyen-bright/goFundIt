package handlers

import (
	"log"
	"net/http"

	gorilla "github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
)

type WebSocketHandler struct {
	hub             *websocket.Hub
	campaignService interfaces.CampaignService
}

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//TODO: Add origin check
		return true
	},
}

func NewWebSocketHandler(hub *websocket.Hub, campaignService interfaces.CampaignService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:             hub,
		campaignService: campaignService,
	}
}

func (h *WebSocketHandler) HandleCampaignWebSocket(c *gin.Context) {
	campaignID := GetCampaignID(c)
	claims := c.MustGet("claims").(jwt.Claims)

	//TODO: verify campaign key

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := websocket.NewClient(h.hub, conn, campaignID, claims.Handle)

	client.Hub.Register(client)

	// Start write pump in a goroutine
	go client.WritePump()
}
