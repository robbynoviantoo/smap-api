package repository

import (
	"context"
	"database/sql"
	"fmt"
	"smap-api/internal/model"
	"strings"
	"time"
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

func (r *AssetRepository) CreateAsset(ctx context.Context, asset *model.AssetCreateRequest) error {
	now := time.Now()
	query := `INSERT INTO assets (name, image, asset, no_asset, location, building, category, sub_category, merk, size, unit, status, available_status, last_maintenance, next_maintenance, remarks, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		asset.Name,
		asset.Image,
		asset.AssetCode,
		asset.NoAsset,
		asset.Location,
		asset.Building,
		asset.Category,
		asset.SubCategory,
		asset.Merk,
		asset.Size,
		asset.Unit,
		asset.Status,
		asset.AvailableStatus,
		asset.LastMaintenance,
		asset.NextMaintenance,
		asset.Remarks,
		now,
		now)
	return err
}

func (r *AssetRepository) GetAssetByID(ctx context.Context, id uint) (*model.Asset, error) {
	query := `
	SELECT id, name, image, asset, no_asset, location, building, category, sub_category, merk, size, unit, status, available_status, last_maintenance, next_maintenance, remarks, created_at, updated_at 
	FROM assets
	WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	var asset model.Asset
	var image, subCategory, remarks, status sql.NullString

	err := row.Scan(
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
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if image.Valid { asset.Image = image.String }
	if subCategory.Valid { asset.SubCategory = subCategory.String }
	if remarks.Valid { asset.Remarks = remarks.String }
	if status.Valid { asset.Status = status.String }

	return &asset, nil
}

func (r *AssetRepository) UpdateAsset(ctx context.Context, asset *model.AssetUpdateRequest) error {
	now := time.Now()
	query := `UPDATE assets SET name = ?, image = ?, asset = ?, no_asset = ?, location = ?, building = ?, category = ?, sub_category = ?, merk = ?, size = ?, unit = ?, status = ?, available_status = ?, last_maintenance = ?, next_maintenance = ?, remarks = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		asset.Name,
		asset.Image,
		asset.AssetCode,
		asset.NoAsset,
		asset.Location,
		asset.Building,
		asset.Category,
		asset.SubCategory,
		asset.Merk,
		asset.Size,
		asset.Unit,
		asset.Status,
		asset.AvailableStatus,
		asset.LastMaintenance,
		asset.NextMaintenance,
		asset.Remarks,
		now,
		asset.ID)
	return err
}

// DeleteWithBackup menyimpan backup ke asset_deletes lalu menghapus asset, dalam satu transaksi.
func (r *AssetRepository) DeleteWithBackup(ctx context.Context, id uint, deletedBy *uint) error {
	// 1. Ambil data asset dulu
	asset, err := r.GetAssetByID(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil {
		return fmt.Errorf("asset with id %d not found", id)
	}

	// 2. Mulai transaksi
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 3. Insert backup
	_, err = tx.ExecContext(ctx, `
		INSERT INTO asset_deletes
		  (asset_id, name, image, asset, no_asset, location, building, category,
		   sub_category, merk, size, unit, status, available_status,
		   last_maintenance, next_maintenance, remarks,
		   asset_created_at, asset_updated_at, deleted_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		asset.ID, asset.Name, asset.Image, asset.AssetCode, asset.NoAsset,
		asset.Location, asset.Building, asset.Category, asset.SubCategory,
		asset.Merk, asset.Size, asset.Unit, asset.Status, asset.AvailableStatus,
		asset.LastMaintenance, asset.NextMaintenance, asset.Remarks,
		asset.CreatedAt, asset.UpdatedAt,
		deletedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to backup asset: %w", err)
	}

	// 4. Hapus asset asli
	_, err = tx.ExecContext(ctx, `DELETE FROM assets WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	return tx.Commit()
}


// GetDashboardStats mengembalikan ringkasan jumlah per status asset.
func (r *AssetRepository) GetDashboardStats(ctx context.Context) (total, maintenance, good, broken, borrowed int, err error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COUNT(*) AS total,
			SUM(CASE WHEN status = 'maintenance' THEN 1 ELSE 0 END) AS maintenance,
			SUM(CASE WHEN status = 'good' THEN 1 ELSE 0 END) AS good,
			SUM(CASE WHEN status = 'broken' THEN 1 ELSE 0 END) AS broken,
			SUM(CASE WHEN available_status = 'borrowed' THEN 1 ELSE 0 END) AS borrowed
		FROM assets`)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&total, &maintenance, &good, &broken, &borrowed)
	}
	return
}

// GetMonthlyMaintenance mengembalikan jumlah asset yang perlu maintenance tiap bulan (12 bulan ke depan dari sekarang).
func (r *AssetRepository) GetMonthlyMaintenance(ctx context.Context) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	now := time.Now()
	for i := 0; i < 12; i++ {
		d := now.AddDate(0, i, 0)
		var count int
		err := r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM assets WHERE next_maintenance IS NOT NULL AND MONTH(next_maintenance) = ? AND YEAR(next_maintenance) = ?`,
			d.Month(), d.Year(),
		).Scan(&count)
		if err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"month": d.Format("Jan"),
			"count": count,
		})
	}
	return result, nil
}

// GetDeletedAssets mengambil list asset yang sudah dihapus dari tabel asset_deletes.
func (r *AssetRepository) GetDeletedAssets(ctx context.Context, limit, offset int) ([]model.AssetDelete, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, asset_id, name,
		       COALESCE(image,''), COALESCE(asset,''), COALESCE(no_asset,''),
		       COALESCE(location,''), COALESCE(building,''), COALESCE(category,''),
		       COALESCE(sub_category,''), COALESCE(merk,''), COALESCE(size,''),
		       COALESCE(unit,''), COALESCE(status,''), COALESCE(available_status,''),
		       last_maintenance, next_maintenance, COALESCE(remarks,''),
		       asset_created_at, asset_updated_at, deleted_by, created_at
		FROM asset_deletes
		ORDER BY id DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AssetDelete
	for rows.Next() {
		var d model.AssetDelete
		if err := rows.Scan(
			&d.ID, &d.AssetID, &d.Name,
			&d.Image, &d.AssetCode, &d.NoAsset,
			&d.Location, &d.Building, &d.Category,
			&d.SubCategory, &d.Merk, &d.Size,
			&d.Unit, &d.Status, &d.AvailableStatus,
			&d.LastMaintenance, &d.NextMaintenance, &d.Remarks,
			&d.AssetCreatedAt, &d.AssetUpdatedAt, &d.DeletedBy, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, nil
}

// CountDeletedAssets mengembalikan total jumlah record di asset_deletes.
func (r *AssetRepository) CountDeletedAssets(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM asset_deletes`).Scan(&count)
	return count, err
}
