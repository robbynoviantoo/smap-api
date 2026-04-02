package handler

import (
	"smap-api/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	assetRepo *repository.AssetRepository
}

func NewDashboardHandler(assetRepo *repository.AssetRepository) *DashboardHandler {
	return &DashboardHandler{assetRepo: assetRepo}
}

// GetStats godoc
// GET /api/v1/dashboard/stats
// Mengembalikan ringkasan statistik asset untuk dashboard.
func (h *DashboardHandler) GetStats(c *fiber.Ctx) error {
	total, maintenance, good, broken, borrowed, err := h.assetRepo.GetDashboardStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal mengambil statistik dashboard"})
	}

	monthly, err := h.assetRepo.GetMonthlyMaintenance(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal mengambil data maintenance bulanan"})
	}

	return c.JSON(fiber.Map{
		"success":              true,
		"totalAssets":          total,
		"maintenanceAssets":    maintenance,
		"goodAssets":           good,
		"borrowedAssets":       borrowed,
		"monthlyMaintenanceData": monthly,
		"assetStatusData": fiber.Map{
			"Good":        good,
			"Maintenance": maintenance,
			"Broken":      broken,
		},
	})
}
