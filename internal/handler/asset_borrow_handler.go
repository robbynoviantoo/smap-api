package handler

import (
	"smap-api/internal/model"
	"smap-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AssetBorrowHandler struct {
	svc *service.AssetBorrowService
}

func NewAssetBorrowHandler(svc *service.AssetBorrowService) *AssetBorrowHandler {
	return &AssetBorrowHandler{svc: svc}
}

// Borrow godoc
// POST /api/v1/asset_borrow/:id/borrow
// Body: { "remark": "..." }
func (h *AssetBorrowHandler) Borrow(c *fiber.Ctx) error {
	assetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid asset ID"})
	}

	var req model.BorrowRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Remark == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "remark wajib diisi"})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.svc.RequestBorrow(c.Context(), uint(assetID), userID.(uint), req.Remark); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pengajuan peminjaman berhasil dikirim. Menunggu approval."})
}

// Return godoc
// POST /api/v1/asset_borrow/:id/return
// Body: { "remark": "..." }
func (h *AssetBorrowHandler) Return(c *fiber.Ctx) error {
	assetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid asset ID"})
	}

	var req model.ReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.Remark == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "remark wajib diisi"})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.svc.ReturnAsset(c.Context(), uint(assetID), userID.(uint), req.Remark); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Berhasil: asset returned dan status kembali available."})
}

// GetAllPendings godoc
// GET /api/v1/asset_borrow_pending/
func (h *AssetBorrowHandler) GetAllPendings(c *fiber.Ctx) error {
	list, err := h.svc.GetAllPendings(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch pending borrow requests"})
	}
	return c.JSON(list)
}

// ReviewPending godoc
// POST /api/v1/asset_borrow_pending/:id/review
// Body: { "status": "approved"|"rejected", "reject_reason": "..." }
func (h *AssetBorrowHandler) ReviewPending(c *fiber.Ctx) error {
	pendingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid pending ID"})
	}

	var req model.BorrowPendingReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	reviewerID := c.Locals("user_id")
	if reviewerID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := h.svc.ReviewBorrow(c.Context(), uint(pendingID), &req, reviewerID.(uint)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pengajuan borrow berhasil diproses."})
}
