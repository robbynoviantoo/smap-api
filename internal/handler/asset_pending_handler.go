package handler

import (
	"smap-api/internal/model"
	"smap-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AssetPendingHandler struct {
	svc *service.AssetPendingService
}

func NewAssetPendingHandler(svc *service.AssetPendingService) *AssetPendingHandler {
	return &AssetPendingHandler{svc: svc}
}

func (h *AssetPendingHandler) CreatePending(c *fiber.Ctx) error {
	var req model.AssetPending
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Get user_id from token
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.UserID = userID.(uint)

	if err := h.svc.CreatePending(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to create pending asset",
			"detail": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "asset submitted for pending review",
		"data":    req,
	})
}

func (h *AssetPendingHandler) GetAllPendings(c *fiber.Ctx) error {
	list, err := h.svc.GetAllPendings(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to fetch asset pendings",
			"detail": err.Error(),
		})
	}
	return c.JSON(list)
}

func (h *AssetPendingHandler) ReviewPending(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format"})
	}

	var req model.AssetPendingReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// reviewer must be from local context
	reviewerID := c.Locals("user_id")
	if reviewerID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.svc.ReviewPendingAsset(c.Context(), uint(id), &req, reviewerID.(uint)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to review pending asset",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "asset review successful",
	})
}
