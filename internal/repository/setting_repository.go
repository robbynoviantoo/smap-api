package repository

import (
	"context"
	"database/sql"
)

type SettingRepository struct {
	db *sql.DB
}

func NewSettingRepository(db *sql.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) GetValue(ctx context.Context, key string, defaultValue string) string {
	query := `SELECT setting_value FROM settings WHERE setting_key = ?`
	var val string
	err := r.db.QueryRowContext(ctx, query, key).Scan(&val)
	if err != nil {
		if err == sql.ErrNoRows {
			return defaultValue
		}
		return defaultValue // or log the error
	}
	return val
}
