package handler

import (
	"smap-api/internal/model"
	"smap-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AssetMaintenanceHandler struct {
	svc *service.AssetMaintenanceService
}

func NewAssetMaintenanceHandler(svc *service.AssetMaintenanceService) *AssetMaintenanceHandler {
	return &AssetMaintenanceHandler{svc: svc}
}

func (h *AssetMaintenanceHandler) Schedule(c *fiber.Ctx) error {
	idStr := c.Params("id")
	assetID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid asset ID"})
	}

	var req model.MaintenanceScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.svc.ScheduleMaintenance(c.Context(), uint(assetID), &req, userID.(uint)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Jadwal Next Maintenance berhasil diatur."})
}

func (h *AssetMaintenanceHandler) Start(c *fiber.Ctx) error {
	idStr := c.Params("id")
	assetID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid asset ID"})
	}

	var req model.MaintenanceStartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	needsApproval, err := h.svc.StartMaintenance(c.Context(), uint(assetID), &req, userID.(uint))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if needsApproval {
		return c.JSON(fiber.Map{"message": "Permintaan maintenance berhasil diajukan dan menunggu approval."})
	}
	return c.JSON(fiber.Map{"message": "Asset berhasil masuk status Maintenance."})
}

func (h *AssetMaintenanceHandler) Finish(c *fiber.Ctx) error {
	idStr := c.Params("id")
	assetID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid asset ID"})
	}

	var req model.MaintenanceFinishRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.svc.FinishMaintenance(c.Context(), uint(assetID), &req, userID.(uint)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Perbaikan berhasil diselesaikan. Status kembali ke Good."})
}

func (h *AssetMaintenanceHandler) ReviewPending(c *fiber.Ctx) error {
	idStr := c.Params("id")
	pendingID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid pending ID"})
	}

	var req model.MaintenancePendingReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	reviewerID := c.Locals("user_id")
	if reviewerID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.svc.ReviewPendingStart(c.Context(), uint(pendingID), &req, reviewerID.(uint)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Maintenance pending diproses berhasil."})
}

func (h *AssetMaintenanceHandler) GetAllPendings(c *fiber.Ctx) error {
	list, err := h.svc.GetAllPendingRequests(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch pending requests"})
	}
	return c.JSON(list)
}
