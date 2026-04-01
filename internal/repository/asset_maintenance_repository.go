package repository

import (
	"context"
	"database/sql"
	"smap-api/internal/model"
	"time"
)

type AssetMaintenanceRepository struct {
	db *sql.DB
}

func NewAssetMaintenanceRepository(db *sql.DB) *AssetMaintenanceRepository {
	return &AssetMaintenanceRepository{db: db}
}

// Transactional: Update asset's next_maintenance and log schedule
func (r *AssetMaintenanceRepository) ScheduleMaintenanceTx(ctx context.Context, assetID uint, nextDate time.Time, userID uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update asset
	_, err = tx.ExecContext(ctx, `UPDATE assets SET next_maintenance = ?, updated_at = ? WHERE id = ?`, nextDate, time.Now(), assetID)
	if err != nil {
		return err
	}

	// Create maintenance log
	_, err = tx.ExecContext(ctx, `INSERT INTO asset_maintenances (asset_id, user_id, type, next_maintenance_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		assetID, userID, "scheduled", nextDate, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Check if a pending start request exists
func (r *AssetMaintenanceRepository) GetPendingStartByAssetID(ctx context.Context, assetID uint) (*model.AssetMaintenancePending, error) {
	query := `SELECT id, asset_id, user_id, notes, status FROM asset_maintenance_pendings WHERE asset_id = ? AND status = 'pending' LIMIT 1`
	var p model.AssetMaintenancePending
	err := r.db.QueryRowContext(ctx, query, assetID).Scan(&p.ID, &p.AssetID, &p.UserID, &p.Notes, &p.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // not found
		}
		return nil, err
	}
	return &p, nil
}

// Create a pending start request
func (r *AssetMaintenanceRepository) CreatePendingRequest(ctx context.Context, p *model.AssetMaintenancePending) error {
	now := time.Now()
	query := `INSERT INTO asset_maintenance_pendings (asset_id, user_id, notes, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, p.AssetID, p.UserID, p.Notes, "pending", now, now)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		p.ID = uint(id)
	}
	p.Status = "pending"
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

// Direct Start Maintenance Transaction
func (r *AssetMaintenanceRepository) StartMaintenanceTx(ctx context.Context, assetID uint, userID uint, notes string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// Update asset status
	_, err = tx.ExecContext(ctx, `UPDATE assets SET status = 'maintenance', updated_at = ? WHERE id = ?`, now, assetID)
	if err != nil {
		return err
	}

	// Insert log
	_, err = tx.ExecContext(ctx, `INSERT INTO asset_maintenances (asset_id, user_id, type, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		assetID, userID, "started", notes, now, now)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Transactional: Finish Maintenance
func (r *AssetMaintenanceRepository) FinishMaintenanceTx(ctx context.Context, assetID uint, originalNextDate *time.Time, newNextDate *time.Time, userID uint, notes string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	// last_maintenance is today startOfDay, let's just use time.Now() 
	
	finalNextDate := originalNextDate
	if newNextDate != nil {
		finalNextDate = newNextDate
	}

	// Update asset status
	_, err = tx.ExecContext(ctx, `UPDATE assets SET status = 'good', last_maintenance = ?, next_maintenance = ?, updated_at = ? WHERE id = ?`, now, finalNextDate, now, assetID)
	if err != nil {
		return err
	}

	// Insert log
	_, err = tx.ExecContext(ctx, `INSERT INTO asset_maintenances (asset_id, user_id, type, notes, next_maintenance_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		assetID, userID, "finished", notes, finalNextDate, now, now)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Pending List
func (r *AssetMaintenanceRepository) GetAllPendingRequests(ctx context.Context) ([]model.AssetMaintenancePending, error) {
	query := `SELECT id, asset_id, user_id, notes, status, created_at FROM asset_maintenance_pendings ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AssetMaintenancePending
	for rows.Next() {
		var p model.AssetMaintenancePending
		if err := rows.Scan(&p.ID, &p.AssetID, &p.UserID, &p.Notes, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// Update pending status transactionally with log
func (r *AssetMaintenanceRepository) ReviewPendingTx(ctx context.Context, id uint, status string, rejectReason string, reviewerID uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := time.Now()

	var reason interface{} = rejectReason
	if rejectReason == "" {
		reason = nil
	}

	// update pending table
	_, err = tx.ExecContext(ctx, `UPDATE asset_maintenance_pendings SET status = ?, reviewed_by = ?, reviewed_at = ?, reject_reason = ?, updated_at = ? WHERE id = ?`,
		status, reviewerID, now, reason, now, id)
	if err != nil {
		return err
	}

	// if approved, also do the "StartMaintenance" logic
	if status == "approved" {
		// we need asset ID and user ID and notes from the pending request
		var assetID, originalUserID uint
		var notes string
		err = tx.QueryRowContext(ctx, `SELECT asset_id, user_id, notes FROM asset_maintenance_pendings WHERE id = ?`, id).Scan(&assetID, &originalUserID, &notes)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `UPDATE assets SET status = 'maintenance', updated_at = ? WHERE id = ?`, now, assetID)
		if err != nil {
			return err
		}

		// Insert log authored by original user
		_, err = tx.ExecContext(ctx, `INSERT INTO asset_maintenances (asset_id, user_id, type, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			assetID, originalUserID, "started", notes, now, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
