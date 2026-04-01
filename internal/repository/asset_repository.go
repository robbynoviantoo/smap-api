package repository

import (
	"context"
	"database/sql"
	"fmt"
	"smap-api/internal/model"
	"strings"
)

type AssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) GetAssets(ctx context.Context, limit, offset int, filter model.AssetFilter) ([]model.Asset, error) {
	var whereClauses []string
	var args []interface{}

	if filter.Name != "" {
		whereClauses = append(whereClauses, "name LIKE ?")
		args = append(args, "%"+filter.Name+"%")
	}
	if filter.AssetCode != "" {
		whereClauses = append(whereClauses, "asset LIKE ?")
		args = append(args, "%"+filter.AssetCode+"%")
	}
	if filter.NoAsset != "" {
		whereClauses = append(whereClauses, "no_asset LIKE ?")
		args = append(args, "%"+filter.NoAsset+"%")
	}
	if filter.Location != "" {
		whereClauses = append(whereClauses, "location LIKE ?")
		args = append(args, "%"+filter.Location+"%")
	}
	if filter.Building != "" {
		whereClauses = append(whereClauses, "building LIKE ?")
		args = append(args, "%"+filter.Building+"%")
	}
	if filter.Category != "" {
		whereClauses = append(whereClauses, "category LIKE ?")
		args = append(args, "%"+filter.Category+"%")
	}
	if filter.Status != "" {
		whereClauses = append(whereClauses, "status LIKE ?")
		args = append(args, "%"+filter.Status+"%")
	}
	if filter.AvailableStatus != "" {
		whereClauses = append(whereClauses, "available_status LIKE ?")
		args = append(args, "%"+filter.AvailableStatus+"%")
	}

	whereStmt := ""
	if len(whereClauses) > 0 {
		whereStmt = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
	SELECT id, name, image, asset, no_asset, location, building, category, sub_category, merk, size, unit, status, available_status, last_maintenance, next_maintenance, remarks, created_at, updated_at 
	FROM assets
	%s
	ORDER BY id DESC
	LIMIT ? OFFSET ?`, whereStmt)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []model.Asset

	for rows.Next() {
		var asset model.Asset

		var image, subCategory, remarks, status sql.NullString

		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&image,
			&asset.AssetCode,
			&asset.NoAsset,
			&asset.Location,
			&asset.Building,
			&asset.Category,
			&subCategory,
			&asset.Merk,
			&asset.Size,
			&asset.Unit,
			&status,
			&asset.AvailableStatus,
			&asset.LastMaintenance,
			&asset.NextMaintenance,
			&remarks,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// mapping
		if image.Valid {
			asset.Image = image.String
		}
		if subCategory.Valid {
			asset.SubCategory = subCategory.String
		}
		if remarks.Valid {
			asset.Remarks = remarks.String
		}
		if status.Valid {
			asset.Status = status.String
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *AssetRepository) CountAssets(ctx context.Context, filter model.AssetFilter) (int, error) {
	var whereClauses []string
	var args []interface{}

	if filter.Name != "" {
		whereClauses = append(whereClauses, "name LIKE ?")
		args = append(args, "%"+filter.Name+"%")
	}
	if filter.AssetCode != "" {
		whereClauses = append(whereClauses, "asset LIKE ?")
		args = append(args, "%"+filter.AssetCode+"%")
	}
	if filter.NoAsset != "" {
		whereClauses = append(whereClauses, "no_asset LIKE ?")
		args = append(args, "%"+filter.NoAsset+"%")
	}
	if filter.Location != "" {
		whereClauses = append(whereClauses, "location LIKE ?")
		args = append(args, "%"+filter.Location+"%")
	}
	if filter.Building != "" {
		whereClauses = append(whereClauses, "building LIKE ?")
		args = append(args, "%"+filter.Building+"%")
	}
	if filter.Category != "" {
		whereClauses = append(whereClauses, "category LIKE ?")
		args = append(args, "%"+filter.Category+"%")
	}
	if filter.Status != "" {
		whereClauses = append(whereClauses, "status LIKE ?")
		args = append(args, "%"+filter.Status+"%")
	}
	if filter.AvailableStatus != "" {
		whereClauses = append(whereClauses, "available_status LIKE ?")
		args = append(args, "%"+filter.AvailableStatus+"%")
	}

	whereStmt := ""
	if len(whereClauses) > 0 {
		whereStmt = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM assets %s`, whereStmt)
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}
