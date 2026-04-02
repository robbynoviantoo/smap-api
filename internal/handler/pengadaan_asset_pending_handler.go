package handler

import (
	"smap-api/internal/model"
	"smap-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PengadaanAssetPendingHandler struct {
	svc *service.PengadaanAssetPendingService
}

func NewPengadaanAssetPendingHandler(svc *service.PengadaanAssetPendingService) *PengadaanAssetPendingHandler {
	return &PengadaanAssetPendingHandler{svc: svc}
}

func (h *PengadaanAssetPendingHandler) getUserID(c *fiber.Ctx) (uint, bool) {
	uid := c.Locals("user_id")
	if uid == nil {
		return 0, false
	}
	return uid.(uint), true
}

// Create godoc
// POST /api/v1/pengadaan_asset_pending
// Body: PengadaanAssetCreateRequest
func (h *PengadaanAssetPendingHandler) Create(c *fiber.Ctx) error {
	userID, ok := h.getUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req model.PengadaanAssetCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	p, err := h.svc.Create(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pengajuan pengadaan berhasil dikirim. Menunggu approval.",
		"data":    p,
	})
}

// GetAll godoc
// GET /api/v1/pengadaan_asset_pending
func (h *PengadaanAssetPendingHandler) GetAll(c *fiber.Ctx) error {
	list, err := h.svc.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch pengadaan pendings"})
	}
	if list == nil {
		list = []model.PengadaanAssetPending{}
	}
	return c.JSON(list)
}

// GetByID godoc
// GET /api/v1/pengadaan_asset_pending/:id
func (h *PengadaanAssetPendingHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	p, err := h.svc.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(p)
}

// Update godoc
// PUT /api/v1/pengadaan_asset_pending/:id
func (h *PengadaanAssetPendingHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}

	var req model.PengadaanAssetUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.svc.Update(c.Context(), uint(id), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pengadaan berhasil diperbarui."})
}

// Delete godoc
// DELETE /api/v1/pengadaan_asset_pending/:id
func (h *PengadaanAssetPendingHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	if err := h.svc.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"deleted": true})
}

// Approve godoc
// POST /api/v1/pengadaan_asset_pending/:id/approve
func (h *PengadaanAssetPendingHandler) Approve(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	reviewerID, ok := h.getUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := h.svc.Approve(c.Context(), uint(id), reviewerID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pengadaan berhasil disetujui."})
}

// Reject godoc
// POST /api/v1/pengadaan_asset_pending/:id/reject
// Body: { "reject_reason": "..." }
func (h *PengadaanAssetPendingHandler) Reject(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	reviewerID, ok := h.getUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var body struct {
		RejectReason string `json:"reject_reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.svc.Reject(c.Context(), uint(id), reviewerID, body.RejectReason); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Pengadaan berhasil ditolak."})
}
