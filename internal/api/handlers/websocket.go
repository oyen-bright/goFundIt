package handlers

import (
	"net/http"

	gorilla "github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/oyen-bright/goFundIt/internal/services/interfaces"
	"github.com/oyen-bright/goFundIt/pkg/jwt"
	"github.com/oyen-bright/goFundIt/pkg/websocket"
)

// WebSocketHandler handles WebSocket connections
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

// NewWebSocketHandler creates a new instance of WebSocketHandler
func NewWebSocketHandler(hub *websocket.Hub, campaignService interfaces.CampaignService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:             hub,
		campaignService: campaignService,
	}
}

// @Summary Campaign WebSocket Connection
// @Description Establishes a WebSocket connection for real-time updates about campaign activities
// @Tags websocket
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param campaignID path string true "Campaign ID"
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} BadRequestResponse "Invalid campaign ID"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized"
// @Failure 404 {object} response "Campaign not found"
// @Router /ws/campaign/{campaignID} [get]
func (h *WebSocketHandler) HandleCampaignWebSocket(c *gin.Context) {
	campaignID := GetCampaignID(c)
	claims := c.MustGet("claims").(jwt.Claims)

	// Verify campaign
	if _, err := h.campaignService.GetCampaignByID(campaignID, getCampaignKey(c)); err != nil {
		FromError(c, err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := websocket.NewClient(h.hub, conn, campaignID, claims.Handle)

	client.Hub.Register(client)

	// Start write pump in a goroutine
	go client.WritePump()
}
