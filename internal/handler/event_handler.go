package handler

import (
	"smap-api/internal/model"
	"smap-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

// GetAll godoc
// GET /api/v1/events
func (h *EventHandler) GetAll(c *fiber.Ctx) error {
	list, err := h.svc.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch events"})
	}
	if list == nil {
		list = []model.Event{}
	}
	return c.JSON(list)
}

// Create godoc
// POST /api/v1/events
func (h *EventHandler) Create(c *fiber.Ctx) error {
	var req model.EventCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	event, err := h.svc.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(event)
}

// Update godoc
// PUT /api/v1/events/:id
func (h *EventHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid event ID"})
	}

	var req model.EventUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	event, err := h.svc.Update(c.Context(), uint(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(event)
}

// Delete godoc
// DELETE /api/v1/events/:id
func (h *EventHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid event ID"})
	}

	if err := h.svc.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"deleted": true})
}

// SyncHolidays godoc
// POST /api/v1/events/sync-holidays
func (h *EventHandler) SyncHolidays(c *fiber.Ctx) error {
	count, err := h.svc.SyncHolidays(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"count":   count,
		"message": "Successfully synced " + strconv.Itoa(count) + " holidays",
	})
}
