package websocket

import (
	"log"
	"sync"
)

// Client represents a single WebSocket connection
type Client struct {
	Hub    *Hub
	Send   chan []byte
	UserID uint
}

// Hub maintains all active WebSocket connections
type Hub struct {
	// Registered clients mapped by userID for targeted messaging
	clients map[uint]*Client

	// Broadcast to all clients
	broadcast chan []byte

	// Register a new client
	register chan *Client

	// Unregister a client
	unregister chan *Client

	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop - must be run as a goroutine
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("WebSocket: User %d connected. Total clients: %d", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket: User %d disconnected.", client.UserID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser sends a message to a specific user only
func (h *Hub) SendToUser(userID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[userID]; ok {
		select {
		case client.Send <- message:
		default:
			log.Printf("WebSocket: Failed to send message to user %d, channel full", userID)
		}
	}
}

// Register exposes the register channel
func (h *Hub) Register() chan<- *Client {
	return h.register
}

// Unregister exposes the unregister channel
func (h *Hub) Unregister() chan<- *Client {
	return h.unregister
}
