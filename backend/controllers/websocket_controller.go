package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/garvpatidar/food-delivery/backend/usecases"
	ws "github.com/garvpatidar/food-delivery/backend/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production: validate origin properly
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketController struct {
	hub         *ws.Hub
	cartUsecase usecases.CartUsecase
}

func NewWebSocketController(hub *ws.Hub) *WebSocketController {
	return &WebSocketController{hub: hub}
}

// ServeWS upgrades an HTTP connection to a WebSocket connection
// GET /ws?token=<jwt_token>
func (wsc *WebSocketController) ServeWS(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &ws.Client{
		Hub:    wsc.hub,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	// Register this client with the hub
	wsc.hub.Register() <- client

	// Start goroutines to read/write
	go writePump(client, conn)
	go readPump(client, conn)
}

// writePump sends messages from hub to the WebSocket connection
func writePump(client *ws.Client, conn *websocket.Conn) {
	ticker := time.NewTicker(54 * time.Second) // Ping every 54s
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump listens for messages from the client (e.g., disconnect)
func readPump(client *ws.Client, conn *websocket.Conn) {
	defer func() {
		client.Hub.Unregister() <- client
		conn.Close()
	}()

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// OrderUpdatePayload is what we send to the client via WebSocket
type OrderUpdatePayload struct {
	Type    string `json:"type"`
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BroadcastOrderUpdate sends order status update to a specific user
func BroadcastOrderUpdate(hub *ws.Hub, userID uint, orderID uint, status string) {
	payload := OrderUpdatePayload{
		Type:    "ORDER_STATUS_UPDATE",
		OrderID: orderID,
		Status:  status,
		Message: getStatusMessage(status),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal order update: %v", err)
		return
	}

	hub.SendToUser(userID, data)
}

func getStatusMessage(status string) string {
	messages := map[string]string{
		"pending":          "Your order has been placed!",
		"confirmed":        "Restaurant confirmed your order.",
		"preparing":        "Your food is being prepared!",
		"out_for_delivery": "Your order is on the way!",
		"delivered":        "Order delivered! Enjoy your meal!",
		"cancelled":        "Your order has been cancelled.",
	}
	if msg, ok := messages[status]; ok {
		return msg
	}
	return "Order status updated."
}
