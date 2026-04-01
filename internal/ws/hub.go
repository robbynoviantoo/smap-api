package ws

import (
	"context"
	"encoding/json"
	"log"
	"smap-api/internal/model"
	"smap-api/internal/service"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Message adalah struktur chat
type Message struct {
	Type        string `json:"type"` // message | join | leave
	SenderID    uint   `json:"sender_id"`
	RecipientID uint   `json:"recipient_id"`
	Content     string `json:"content"`
}

// Client represent 1 user connection
type Client struct {
	Conn   *websocket.Conn
	UserID uint
}

type Hub struct {
	mu         sync.RWMutex
	clients    map[*websocket.Conn]*Client
	MessageSvc *service.MessageService
}

// global instance
var GlobalHub = &Hub{
	clients: make(map[*websocket.Conn]*Client),
}

// Register client
func (h *Hub) Register(conn *websocket.Conn, userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = &Client{
		Conn:   conn,
		UserID: userID,
	}

	log.Printf("[WS] User %d connected, total: %d", userID, len(h.clients))
}

// Unregister client
func (h *Hub) Unregister(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, conn)
	log.Printf("[WS] Client disconnected, total: %d", len(h.clients))
}

// Kirim ke user tertentu (PRIVATE CHAT)
func (h *Hub) SendToUser(userID uint, message Message) {
	// Persistence
	if h.MessageSvc != nil {
		dbMsg := &model.Message{
			SenderID:    message.SenderID,
			RecipientID: message.RecipientID,
			Content:     message.Content,
			IsRead:      false,
		}
		if err := h.MessageSvc.SaveMessage(context.Background(), dbMsg); err != nil {
			log.Println("[WS] db save message error:", err)
		}
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log.Println("[WS] Marshal error:", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.UserID == userID {
			err := client.Conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				log.Println("[WS] Write error:", err)
			}
		}
	}
}

// Broadcast ke semua user (optional)
func (h *Hub) Broadcast(message Message) {
	payload, err := json.Marshal(message)
	if err != nil {
		log.Println("[WS] Marshal error:", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Println("[WS] Write error:", err)
		}
	}
}