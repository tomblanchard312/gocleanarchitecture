package interfaces

import (
	"gocleanarchitecture/frameworks/websocket"
	"net/http"
)

type WebSocketHandler struct {
	Hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		Hub: hub,
	}
}

// HandleWebSocket handles WebSocket connections
// This endpoint allows both authenticated and anonymous connections
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Try to get user ID from context (if authenticated)
	userID := ""
	if uid, ok := r.Context().Value("userID").(string); ok {
		userID = uid
	}

	// Upgrade to WebSocket and serve
	websocket.ServeWs(h.Hub, w, r, userID)
}
