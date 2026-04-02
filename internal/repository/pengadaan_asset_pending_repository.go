package repository

import (
	"context"
	"database/sql"
	"smap-api/internal/model"
	"time"
)

type PengadaanAssetPendingRepository struct {
	db *sql.DB
}

func NewPengadaanAssetPendingRepository(db *sql.DB) *PengadaanAssetPendingRepository {
	return &PengadaanAssetPendingRepository{db: db}
}

func (r *PengadaanAssetPendingRepository) Create(ctx context.Context, p *model.PengadaanAssetPending) error {
	now := time.Now()
	query := `INSERT INTO pengadaan_asset_pendings
		(user_id, request_date, category, name, merk, spec, qty, unit, image, remark, priority, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending', ?, ?)`

	res, err := r.db.ExecContext(ctx, query,
		p.UserID, p.RequestDate, nullStr(p.Category), p.Name, nullStr(p.Merk), nullStr(p.Spec),
		p.Qty, nullStr(p.Unit), nullStr(p.Image), nullStr(p.Remark), p.Priority,
		now, now,
	)
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

func (r *PengadaanAssetPendingRepository) GetAll(ctx context.Context) ([]model.PengadaanAssetPending, error) {
	query := `
		SELECT
			pap.id, pap.user_id, pap.request_date, pap.category, pap.name,
			pap.merk, pap.spec, pap.qty, pap.unit, pap.image, pap.remark,
			pap.priority, pap.status,
			pap.reviewed_by, pap.reviewed_at, pap.reject_reason,
			pap.created_at, pap.updated_at,
			u.first_name AS user_first_name
		FROM pengadaan_asset_pendings pap
		LEFT JOIN users u ON u.id = pap.user_id
		ORDER BY pap.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PengadaanAssetPending
	for rows.Next() {
		var p model.PengadaanAssetPending
		var category, merk, spec, unit, image, remark sql.NullString
		var reviewedBy sql.NullInt64
		var reviewedAt sql.NullTime
		var rejectReason sql.NullString
		var userFirstName sql.NullString

		if err := rows.Scan(
			&p.ID, &p.UserID, &p.RequestDate, &category, &p.Name,
			&merk, &spec, &p.Qty, &unit, &image, &remark,
			&p.Priority, &p.Status,
			&reviewedBy, &reviewedAt, &rejectReason,
			&p.CreatedAt, &p.UpdatedAt,
			&userFirstName,
		); err != nil {
			return nil, err
		}
		if category.Valid { p.Category = category.String }
		if merk.Valid { p.Merk = merk.String }
		if spec.Valid { p.Spec = spec.String }
		if unit.Valid { p.Unit = unit.String }
		if image.Valid { p.Image = image.String }
		if remark.Valid { p.Remark = remark.String }
		if reviewedBy.Valid { rb := uint(reviewedBy.Int64); p.ReviewedBy = &rb }
		if reviewedAt.Valid { p.ReviewedAt = &reviewedAt.Time }
		if rejectReason.Valid { rr := rejectReason.String; p.RejectReason = &rr }
		if userFirstName.Valid { p.User = &model.User{ID: p.UserID, FirstName: userFirstName.String} }

		list = append(list, p)
	}
	return list, nil
}

func (r *PengadaanAssetPendingRepository) GetByID(ctx context.Context, id uint) (*model.PengadaanAssetPending, error) {
	query := `SELECT id, user_id, request_date, category, name, merk, spec, qty, unit, image, remark,
		priority, status, reviewed_by, reviewed_at, reject_reason, created_at, updated_at
		FROM pengadaan_asset_pendings WHERE id = ?`

	var p model.PengadaanAssetPending
	var category, merk, spec, unit, image, remark sql.NullString
	var reviewedBy sql.NullInt64
	var reviewedAt sql.NullTime
	var rejectReason sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.RequestDate, &category, &p.Name,
		&merk, &spec, &p.Qty, &unit, &image, &remark,
		&p.Priority, &p.Status,
		&reviewedBy, &reviewedAt, &rejectReason,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	if category.Valid { p.Category = category.String }
	if merk.Valid { p.Merk = merk.String }
	if spec.Valid { p.Spec = spec.String }
	if unit.Valid { p.Unit = unit.String }
	if image.Valid { p.Image = image.String }
	if remark.Valid { p.Remark = remark.String }
	if reviewedBy.Valid { rb := uint(reviewedBy.Int64); p.ReviewedBy = &rb }
	if reviewedAt.Valid { p.ReviewedAt = &reviewedAt.Time }
	if rejectReason.Valid { rr := rejectReason.String; p.RejectReason = &rr }

	return &p, nil
}

// ApproveTx: pindahkan data ke pengadaan_assets lalu update status pending.
func (r *PengadaanAssetPendingRepository) ApproveTx(ctx context.Context, p *model.PengadaanAssetPending, reviewerID uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// Insert ke tabel pengadaan_assets (approved)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO pengadaan_assets
			(user_id, request_date, category, name, merk, spec, qty, unit, image, remarks, priority, updated_by, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.UserID, p.RequestDate, nullStr(p.Category), p.Name, nullStr(p.Merk), nullStr(p.Spec),
		p.Qty, nullStr(p.Unit), nullStr(p.Image), nullStr(p.Remark), p.Priority, reviewerID,
		now, now,
	)
	if err != nil {
		return err
	}

	// Update status pending
	_, err = tx.ExecContext(ctx,
		`UPDATE pengadaan_asset_pendings SET status = 'approved', reviewed_by = ?, reviewed_at = ?, updated_at = ? WHERE id = ?`,
		reviewerID, now, now, p.ID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PengadaanAssetPendingRepository) Reject(ctx context.Context, id uint, reviewerID uint, rejectReason string) error {
	now := time.Now()
	var reason interface{} = rejectReason
	if rejectReason == "" {
		reason = nil
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE pengadaan_asset_pendings SET status = 'rejected', reviewed_by = ?, reviewed_at = ?, reject_reason = ?, updated_at = ? WHERE id = ?`,
		reviewerID, now, reason, now, id,
	)
	return err
}

// UpdatePending update data pengadaan pending (hanya yang masih pending).
func (r *PengadaanAssetPendingRepository) Update(ctx context.Context, id uint, req *model.PengadaanAssetUpdateRequest) error {
	now := time.Now()
	query := `UPDATE pengadaan_asset_pendings
		SET category = COALESCE(NULLIF(?, ''), category),
		    name = COALESCE(NULLIF(?, ''), name),
		    merk = COALESCE(NULLIF(?, ''), merk),
		    spec = COALESCE(NULLIF(?, ''), spec),
		    qty = CASE WHEN ? > 0 THEN ? ELSE qty END,
		    unit = COALESCE(NULLIF(?, ''), unit),
		    image = COALESCE(NULLIF(?, ''), image),
		    remark = COALESCE(NULLIF(?, ''), remark),
		    priority = COALESCE(NULLIF(?, ''), priority),
		    updated_at = ?
		WHERE id = ? AND status = 'pending'`
	_, err := r.db.ExecContext(ctx, query,
		req.Category, req.Name, req.Merk, req.Spec,
		req.Qty, req.Qty, req.Unit, req.Image, req.Remark, req.Priority,
		now, id,
	)
	return err
}

func (r *PengadaanAssetPendingRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pengadaan_asset_pendings WHERE id = ?`, id)
	return err
}

// nullStr mengembalikan nil jika string kosong, untuk insert SQL.
func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
