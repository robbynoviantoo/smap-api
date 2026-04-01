package ws

import (
	"log"
	"strconv"

	"github.com/gofiber/websocket/v2"
)

func ChatHandler(c *websocket.Conn) {
	// ambil user_id dari query / token
	userIDStr := c.Query("user_id")
	userID, _ := strconv.Atoi(userIDStr)

	GlobalHub.Register(c, uint(userID))
	defer GlobalHub.Unregister(c)

	for {
		var msg Message

		if err := c.ReadJSON(&msg); err != nil {
			log.Println("[WS] Read error:", err)
			break
		}

		// kirim ke recipient
		GlobalHub.SendToUser(msg.RecipientID, msg)
	}
}