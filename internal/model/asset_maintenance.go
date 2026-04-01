package model

import "time"

type AssetMaintenance struct {
	ID                  uint       `json:"id"`
	AssetID             uint       `json:"asset_id"`
	UserID              uint       `json:"user_id"`
	Type                string     `json:"type"` // scheduled, started, finished
	Notes               string     `json:"notes"`
	NextMaintenanceDate *time.Time `json:"next_maintenance_date"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (optional)
	Asset *Asset `json:"asset,omitempty"`
	User  *User  `json:"user,omitempty"`
}

type AssetMaintenancePending struct {
	ID           uint       `json:"id"`
	AssetID      uint       `json:"asset_id"`
	UserID       uint       `json:"user_id"`
	Notes        string     `json:"notes"`
	Status       string     `json:"status"` // pending, approved, rejected
	ReviewedBy   *uint      `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	RejectReason *string    `json:"reject_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Asset    *Asset `json:"asset,omitempty"`
	User     *User  `json:"user,omitempty"`
	Reviewer *User  `json:"reviewer,omitempty"`
}

type MaintenanceScheduleRequest struct {
	NextMaintenanceDate string `json:"next_maintenance_date"` // we'll parse this as a date string in the handler
}

type MaintenanceStartRequest struct {
	Notes string `json:"notes"`
}

type MaintenanceFinishRequest struct {
	Notes               string  `json:"notes"`
	NextMaintenanceDate *string `json:"next_maintenance_date,omitempty"` // optional string date
}

type MaintenancePendingReviewRequest struct {
	Status       string `json:"status"` // approved, rejected
	RejectReason string `json:"reject_reason,omitempty"`
}