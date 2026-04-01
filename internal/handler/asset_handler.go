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