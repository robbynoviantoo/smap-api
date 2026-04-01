package model

import "time"

type AssetPending struct {
	ID               uint       `json:"id"`
	UserID           uint       `json:"user_id"`
	Name             string     `json:"name"`
	Image            string     `json:"image,omitempty"`
	AssetCode        string     `json:"asset,omitempty"`
	NoAsset          string     `json:"no_asset"`
	Location         string     `json:"location,omitempty"`
	Building         string     `json:"building,omitempty"`
	Category         string     `json:"category,omitempty"`
	SubCategory      string     `json:"sub_category,omitempty"`
	Merk             string     `json:"merk,omitempty"`
	Size             string     `json:"size,omitempty"`
	Unit             string     `json:"unit,omitempty"`
	Status           string     `json:"status,omitempty"` // The actual asset status (e.g. "good", "Baru")
	LastMaintenance  *time.Time `json:"last_maintenance,omitempty"`
	NextMaintenance  *time.Time `json:"next_maintenance,omitempty"`
	Remarks          string     `json:"remarks,omitempty"`

	// Reviewer info
	ReviewedBy   *uint      `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	RejectReason *string    `json:"reject_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Optional references
	User     *User `json:"user,omitempty"`
	Reviewer *User `json:"reviewer,omitempty"`
}

type AssetPendingReviewRequest struct {
	Status       string `json:"status"` // "approved" or "rejected"
	RejectReason string `json:"reject_reason,omitempty"`
}
