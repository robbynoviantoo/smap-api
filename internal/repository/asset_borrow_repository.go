package repository

import (
	"context"
	"database/sql"
	"smap-api/internal/model"
	"time"
)

type AssetBorrowRepository struct {
	db *sql.DB
}

func NewAssetBorrowRepository(db *sql.DB) *AssetBorrowRepository {
	return &AssetBorrowRepository{db: db}
}

// GetPendingByAssetID cek apakah ada pengajuan borrow yang masih pending untuk asset ini.
func (r *AssetBorrowRepository) GetPendingByAssetID(ctx context.Context, assetID uint) (*model.AssetBorrowPending, error) {
	query := `SELECT id, asset_id, user_id, remark, status FROM asset_borrow_pendings WHERE asset_id = ? AND status = 'pending' LIMIT 1`
	var p model.AssetBorrowPending
	err := r.db.QueryRowContext(ctx, query, assetID).Scan(&p.ID, &p.AssetID, &p.UserID, &p.Remark, &p.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// CreatePendingBorrow menyimpan pengajuan borrow baru dengan status pending.
func (r *AssetBorrowRepository) CreatePendingBorrow(ctx context.Context, p *model.AssetBorrowPending) error {
	now := time.Now()
	query := `INSERT INTO asset_borrow_pendings (asset_id, user_id, remark, status, created_at, updated_at) VALUES (?, ?, ?, 'pending', ?, ?)`
	res, err := r.db.ExecContext(ctx, query, p.AssetID, p.UserID, p.Remark, now, now)
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

// GetAllPendingBorrows mengembalikan semua pengajuan borrow, join dengan asset dan user.
func (r *AssetBorrowRepository) GetAllPendingBorrows(ctx context.Context) ([]model.AssetBorrowPending, error) {
	query := `
		SELECT
			abp.id, abp.asset_id, abp.user_id, abp.remark, abp.status,
			abp.reviewed_by, abp.reviewed_at, abp.reject_reason,
			abp.created_at, abp.updated_at,
			a.name       AS asset_name,
			u.first_name AS user_first_name
		FROM asset_borrow_pendings abp
		LEFT JOIN assets a ON a.id = abp.asset_id
		LEFT JOIN users  u ON u.id = abp.user_id
		ORDER BY abp.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.AssetBorrowPending
	for rows.Next() {
		var p model.AssetBorrowPending
		var assetName, userFirstName string
		if err := rows.Scan(
			&p.ID, &p.AssetID, &p.UserID, &p.Remark, &p.Status,
			&p.ReviewedBy, &p.ReviewedAt, &p.RejectReason,
			&p.CreatedAt, &p.UpdatedAt,
			&assetName, &userFirstName,
		); err != nil {
			return nil, err
		}
		p.Asset = &model.Asset{ID: p.AssetID, Name: assetName}
		p.User = &model.User{ID: p.UserID, FirstName: userFirstName}
		list = append(list, p)
	}
	return list, nil
}

// ReviewBorrowTx secara transaksional memperbarui status pending.
// Jika approved: update available_status asset → 'borrowed' dan log ke asset_transactions.
// Jika rejected: hanya update pending record.
func (r *AssetBorrowRepository) ReviewBorrowTx(ctx context.Context, pendingID uint, status, rejectReason string, reviewerID uint) error {
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

	_, err = tx.ExecContext(ctx,
		`UPDATE asset_borrow_pendings SET status = ?, reviewed_by = ?, reviewed_at = ?, reject_reason = ?, updated_at = ? WHERE id = ?`,
		status, reviewerID, now, reason, now, pendingID,
	)
	if err != nil {
		return err
	}

	if status == "approved" {
		var assetID, originalUserID uint
		var remark string
		err = tx.QueryRowContext(ctx,
			`SELECT asset_id, user_id, remark FROM asset_borrow_pendings WHERE id = ?`, pendingID,
		).Scan(&assetID, &originalUserID, &remark)
		if err != nil {
			return err
		}

		// Set asset available_status → borrowed
		_, err = tx.ExecContext(ctx,
			`UPDATE assets SET available_status = 'borrowed', updated_at = ? WHERE id = ?`,
			now, assetID,
		)
		if err != nil {
			return err
		}

		// Log borrow transaction (authored by the original requester)
		_, err = tx.ExecContext(ctx,
			`INSERT INTO asset_transactions (asset_id, borrower_id, type, remark, action_at, created_at, updated_at) VALUES (?, ?, 'borrow', ?, ?, ?, ?)`,
			assetID, originalUserID, remark, now, now, now,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ReturnAssetTx secara transaksional mencatat return dan mengembalikan available_status → 'available'.
func (r *AssetBorrowRepository) ReturnAssetTx(ctx context.Context, assetID, userID uint, remark string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// Log return transaction
	_, err = tx.ExecContext(ctx,
		`INSERT INTO asset_transactions (asset_id, borrower_id, type, remark, action_at, created_at, updated_at) VALUES (?, ?, 'return', ?, ?, ?, ?)`,
		assetID, userID, remark, now, now, now,
	)
	if err != nil {
		return err
	}

	// Set asset available_status → available
	_, err = tx.ExecContext(ctx,
		`UPDATE assets SET available_status = 'available', updated_at = ? WHERE id = ?`,
		now, assetID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
