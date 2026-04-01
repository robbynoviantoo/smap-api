package model

import "time"

type AssetMaintenance struct {
	ID                   uint       `json:"id"`
	AssetID              uint       `json:"asset_id"`
	UserID               uint       `json:"user_id"`
	Type                 string     `json:"type"` // scheduled, started, finished
	Notes                string     `json:"notes"`
	NextMaintenanceDate  *time.Time `json:"next_maintenance_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (optional)
	Asset *Asset `json:"asset,omitempty"`
	User  *User  `json:"user,omitempty"`
}