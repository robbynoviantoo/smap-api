package handler

import (
	"fmt"
	"math/rand"
	"smap-api/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func (h *AssetHandler) ExportAssets(c *fiber.Ctx) error {
	// For export, we might want to get all assets without pagination, or a very large limit
	// Currently we use GetAssetsWithCount with large limit
	filter := model.AssetFilter{} // grab everything for now
	data, _, err := h.assetSvc.GetAssetsWithCount(c.Context(), 100000, 0, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch assets for export"})
	}

	f := excelize.NewFile()
	sheetName := "Assets"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"id", "name", "image", "asset", "no_asset", "location", "building",
		"category", "sub_category", "merk", "size", "unit", "status",
		"last_maintenance", "next_maintenance", "remarks",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, v := range data {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), v.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), v.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), v.Image)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), v.AssetCode)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), v.NoAsset)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), v.Location)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), v.Building)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), v.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), v.SubCategory)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), v.Merk)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), v.Size)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), v.Unit)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), v.Status)

		if v.LastMaintenance != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), v.LastMaintenance.Format("2006-01-02"))
		}
		if v.NextMaintenance != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), v.NextMaintenance.Format("2006-01-02"))
		}
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), v.Remarks)
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=assets_export.xlsx")

	return f.Write(c.Response().BodyWriter())
}

func (h *AssetHandler) ImportAssets(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no file uploaded"})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to open file"})
	}
	defer f.Close()

	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse excel file"})
	}

	sheets := excelFile.GetSheetList()
	if len(sheets) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no sheets found in excel"})
	}

	rows, err := excelFile.GetRows(sheets[0])
	if err != nil || len(rows) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "empty data or missing headers"})
	}

	headers := rows[0]
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	ctx := c.Context()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		
		getCol := func(colName string) string {
			idx, ok := headerMap[colName]
			if ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		idStr := getCol("id")
		noAsset := getCol("no_asset")

		name := getCol("name")
		image := getCol("image")
		assetCode := getCol("asset")
		location := getCol("location")
		building := getCol("building")
		category := getCol("category")
		subCategory := getCol("sub_category")
		merk := getCol("merk")
		size := getCol("size")
		unit := getCol("unit")
		status := getCol("status")
		remarks := getCol("remarks")

		lastMaint := parseDate(getCol("last_maintenance"))
		nextMaint := parseDate(getCol("next_maintenance"))

		var id uint
		if idStr != "" {
			parsedID, _ := strconv.Atoi(idStr)
			id = uint(parsedID)
		}

		if id > 0 {
			// Try update
			existing, _ := h.assetSvc.GetAssetByID(ctx, id)
			if existing != nil {
				updateReq := &model.AssetUpdateRequest{
					ID:              id,
					Name:            name,
					Image:           image,
					AssetCode:       assetCode,
					NoAsset:         noAsset,
					Location:        location,
					Building:        building,
					Category:        category,
					SubCategory:     subCategory,
					Merk:            merk,
					Size:            size,
					Unit:            unit,
					Status:          status,
					LastMaintenance: lastMaint,
					NextMaintenance: nextMaint,
					Remarks:         remarks,
				}
				_ = h.assetSvc.UpdateAsset(ctx, updateReq)
				continue
			}
		}

		// CREATE path
		if noAsset == "" {
			noAsset = generateAssetUniqueNo()
		}

		createReq := &model.AssetCreateRequest{
			Name:            name,
			Image:           image,
			AssetCode:       assetCode,
			NoAsset:         noAsset,
			Location:        location,
			Building:        building,
			Category:        category,
			SubCategory:     subCategory,
			Merk:            merk,
			Size:            size,
			Unit:            unit,
			Status:          status,
			AvailableStatus: "Tersedia",
			LastMaintenance: lastMaint,
			NextMaintenance: nextMaint,
			Remarks:         remarks,
		}
		_ = h.assetSvc.CreateAsset(ctx, createReq)
	}

	return c.JSON(fiber.Map{"message": "import completed"})
}

func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	// Try parsing standard formats or excel serial format
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		t, err := excelize.ExcelDateToTime(f, false)
		if err == nil {
			return &t
		}
	}

	formats := []string{"2006-01-02", "02/01/2006", "2006/01/02", time.RFC3339}
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return &t
		}
	}
	return nil
}

func generateAssetUniqueNo() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return "AST-" + string(b)
}
