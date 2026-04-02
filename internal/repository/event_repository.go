package repository

import (
	"context"
	"database/sql"
	"smap-api/internal/model"
	"time"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetAll(ctx context.Context) ([]model.Event, error) {
	query := `SELECT id, title, start, ` + "`end`" + `, all_day, color, created_at, updated_at FROM events ORDER BY start ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Event
	for rows.Next() {
		var e model.Event
		var start, end sql.NullTime
		var color sql.NullString
		if err := rows.Scan(&e.ID, &e.Title, &start, &end, &e.AllDay, &color, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		if start.Valid {
			e.Start = &start.Time
		}
		if end.Valid {
			e.End = &end.Time
		}
		if color.Valid {
			e.Color = color.String
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *EventRepository) GetByID(ctx context.Context, id uint) (*model.Event, error) {
	query := `SELECT id, title, start, ` + "`end`" + `, all_day, color, created_at, updated_at FROM events WHERE id = ?`
	var e model.Event
	var start, end sql.NullTime
	var color sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Title, &start, &end, &e.AllDay, &color, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if start.Valid {
		e.Start = &start.Time
	}
	if end.Valid {
		e.End = &end.Time
	}
	if color.Valid {
		e.Color = color.String
	}
	return &e, nil
}

func (r *EventRepository) Create(ctx context.Context, e *model.Event) error {
	now := time.Now()
	query := `INSERT INTO events (title, start, ` + "`end`" + `, all_day, color, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	var colorVal interface{} = nil
	if e.Color != "" {
		colorVal = e.Color
	}

	res, err := r.db.ExecContext(ctx, query, e.Title, e.Start, e.End, e.AllDay, colorVal, now, now)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		e.ID = uint(id)
	}
	e.CreatedAt = now
	e.UpdatedAt = now
	return nil
}

func (r *EventRepository) Update(ctx context.Context, e *model.Event) error {
	now := time.Now()
	query := `UPDATE events SET title = ?, start = ?, ` + "`end`" + ` = ?, all_day = ?, color = ?, updated_at = ? WHERE id = ?`

	var colorVal interface{} = nil
	if e.Color != "" {
		colorVal = e.Color
	}

	_, err := r.db.ExecContext(ctx, query, e.Title, e.Start, e.End, e.AllDay, colorVal, now, e.ID)
	if err != nil {
		return err
	}
	e.UpdatedAt = now
	return nil
}

func (r *EventRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = ?`, id)
	return err
}

// ExistsForHoliday cek apakah event dengan tanggal dan judul yang sama sudah ada (cegah duplikat holiday).
func (r *EventRepository) ExistsForHoliday(ctx context.Context, date, title string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM events WHERE DATE(start) = ? AND title = ?`, date, title,
	).Scan(&count)
	return count > 0, err
}

// BulkCreateHolidays menyimpan banyak event hari libur sekaligus.
func (r *EventRepository) BulkCreateHolidays(ctx context.Context, events []model.Event) (int, error) {
	count := 0
	now := time.Now()
	for _, e := range events {
		exists, err := r.ExistsForHoliday(ctx, e.Start.Format("2006-01-02"), e.Title)
		if err != nil {
			return count, err
		}
		if exists {
			continue
		}

		var colorVal interface{} = nil
		if e.Color != "" {
			colorVal = e.Color
		}

		_, err = r.db.ExecContext(ctx,
			`INSERT INTO events (title, start, `+"`end`"+`, all_day, color, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			e.Title, e.Start, e.End, e.AllDay, colorVal, now, now,
		)
		if err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
