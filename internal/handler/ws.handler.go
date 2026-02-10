package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nguyen1302/realtime-quiz/internal/realtime"
	"github.com/nguyen1302/realtime-quiz/internal/service"
)

type WebSocketHandler interface {
	HandleConnection(c *gin.Context)
}

type webSocketHandler struct {
	realtimeService service.RealtimeService
}

func NewWebSocketHandler(realtimeService service.RealtimeService) WebSocketHandler {
	return &webSocketHandler{
		realtimeService: realtimeService,
	}
}

// HandleConnection upgrades the HTTP connection to WebSocket
// GET /ws
func (h *webSocketHandler) HandleConnection(c *gin.Context) {
	// Check if user is authenticated (from middleware)
	userIDVal, exists := c.Get("userID")
	var userID string
	if exists {
		userID = userIDVal.(uuid.UUID).String()
	} else {
		// do something...
	}

	realtime.ServeWs(h.realtimeService.GetManager().Hub, c, userID)
}
