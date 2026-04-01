package repository

import (
	"context"
	"database/sql"
	"smap-api/internal/model"
	"time"
)

type AssetPendingRepository struct {
	db *sql.DB
}

func NewAssetPendingRepository(db *sql.DB) *AssetPendingRepository {
	return &AssetPendingRepository{db: db}
}

func (r *AssetPendingRepository) CreatePending(ctx context.Context, pending *model.AssetPending) error {
	now := time.Now()
	query := `INSERT INTO asset_pendings (
		user_id, name, image, asset, no_asset, location, building, category, sub_category, 
		merk, size, unit, status, last_maintenance, next_maintenance, remarks, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	res, err := r.db.ExecContext(ctx, query,
		pending.UserID, pending.Name, pending.Image, pending.AssetCode, pending.NoAsset,
		pending.Location, pending.Building, pending.Category, pending.SubCategory,
		pending.Merk, pending.Size, pending.Unit, pending.Status,
		pending.LastMaintenance, pending.NextMaintenance, pending.Remarks, now, now)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	pending.ID = uint(id)
	return nil
}

func (r *AssetPendingRepository) GetPendingByID(ctx context.Context, id uint) (*model.AssetPending, error) {
	query := `SELECT id, user_id, name, image, asset, no_asset, location, building, category, sub_category, merk, size, unit, status, last_maintenance, next_maintenance, remarks, reviewed_by, reviewed_at, reject_reason, created_at, updated_at FROM asset_pendings WHERE id = ?`
	
	row := r.db.QueryRowContext(ctx, query, id)
	var p model.AssetPending
	var image, asset, location, building, category, subCategory, merk, size, unit, status, remarks, rejectReason sql.NullString
	var lastMaint, nextMaint, reviewedAt sql.NullTime
	var reviewedBy sql.NullInt64

	err := row.Scan(
		&p.ID, &p.UserID, &p.Name, &image, &asset, &p.NoAsset, &location, &building, &category, &subCategory,
		&merk, &size, &unit, &status, &lastMaint, &nextMaint, &remarks, &reviewedBy, &reviewedAt, &rejectReason,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // not found
		}
		return nil, err
	}

	if image.Valid { p.Image = image.String }
	if asset.Valid { p.AssetCode = asset.String }
	if location.Valid { p.Location = location.String }
	if building.Valid { p.Building = building.String }
	if category.Valid { p.Category = category.String }
	if subCategory.Valid { p.SubCategory = subCategory.String }
	if merk.Valid { p.Merk = merk.String }
	if size.Valid { p.Size = size.String }
	if unit.Valid { p.Unit = unit.String }
	if status.Valid { p.Status = status.String }
	if remarks.Valid { p.Remarks = remarks.String }
	if rejectReason.Valid {
		rr := rejectReason.String
		p.RejectReason = &rr
	}
	if lastMaint.Valid { p.LastMaintenance = &lastMaint.Time }
	if nextMaint.Valid { p.NextMaintenance = &nextMaint.Time }
	if reviewedAt.Valid { p.ReviewedAt = &reviewedAt.Time }
	if reviewedBy.Valid {
		rb := uint(reviewedBy.Int64)
		p.ReviewedBy = &rb
	}

	return &p, nil
}

func (r *AssetPendingRepository) UpdatePendingStatus(ctx context.Context, id uint, status string, reviewerID uint, rejectReason string) error {
	now := time.Now()
	query := `UPDATE asset_pendings SET status = ?, reviewed_by = ?, reviewed_at = ?, reject_reason = ?, updated_at = ? WHERE id = ?`
	
	var reason interface{} = rejectReason
	if rejectReason == "" {
		reason = nil
	}

	_, err := r.db.ExecContext(ctx, query, status, reviewerID, now, reason, now, id)
	return err
}

func (r *AssetPendingRepository) GetAllPendings(ctx context.Context) ([]model.AssetPending, error) {
	query := `SELECT id, user_id, name, asset, no_asset, category, status, created_at FROM asset_pendings ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AssetPending
	for rows.Next() {
		var p model.AssetPending
		var asset, category, status sql.NullString
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &asset, &p.NoAsset, &category, &status, &p.CreatedAt); err != nil {
			return nil, err
		}
		if asset.Valid { p.AssetCode = asset.String }
		if category.Valid { p.Category = category.String }
		if status.Valid { p.Status = status.String }
		list = append(list, p)
	}
	return list, nil
}
