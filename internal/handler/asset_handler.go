package handler

import (
	"math"
	"smap-api/internal/model"
	"smap-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	assetSvc *service.AssetService
}

func NewAssetHandler(assetSvc *service.AssetService) *AssetHandler {
	return &AssetHandler{assetSvc: assetSvc}
}

func (h *AssetHandler) GetAssets(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	filter := model.AssetFilter{
		Name:            c.Query("name"),
		AssetCode:       c.Query("asset"),
		NoAsset:         c.Query("no_asset"),
		Location:        c.Query("location"),
		Building:        c.Query("building"),
		Category:        c.Query("category"),
		Status:          c.Query("status"),
		AvailableStatus: c.Query("available_status"),
	}

	data, total, err := h.assetSvc.GetAssetsWithCount(c.Context(), limit, offset, filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch assets",
			"detail": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	var req model.AssetCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	err := h.assetSvc.CreateAsset(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create asset",
			"detail": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "asset created successfully",
	})
}

func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	var req model.AssetUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	err := h.assetSvc.UpdateAsset(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update asset",
			"detail": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "asset updated successfully",
	})
}

func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "asset id is required",
		})
	}
	err := h.assetSvc.DeleteAsset(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete asset",
			"detail": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "asset deleted successfully",
	})
}