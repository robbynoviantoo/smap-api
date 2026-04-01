package handler

import (
	"smap-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	svc *service.MessageService
}

func NewChatHandler(svc *service.MessageService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

func (h *ChatHandler) GetHistory(c *fiber.Ctx) error {
	userID1Str := c.Query("user_id1")
	userID2Str := c.Query("user_id2")
	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	u1, err1 := strconv.Atoi(userID1Str)
	u2, err2 := strconv.Atoi(userID2Str)
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id1 and user_id2 are required"})
	}

	messages, err := h.svc.GetChatHistory(c.Context(), uint(u1), uint(u2), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to load chat history"})
	}

	return c.JSON(fiber.Map{
		"data": messages,
	})
}
