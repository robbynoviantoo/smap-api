package model

import "time"

// PengadaanAssetPending adalah approval queue untuk pengajuan pengadaan asset baru.
type PengadaanAssetPending struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	RequestDate time.Time `json:"request_date"`
	Category    string    `json:"category,omitempty"`
	Name        string    `json:"name"`
	Merk        string    `json:"merk,omitempty"`
	Spec        string    `json:"spec,omitempty"`
	Qty         int       `json:"qty"`
	Unit        string    `json:"unit,omitempty"`
	Image       string    `json:"image,omitempty"`
	Remark      string    `json:"remark,omitempty"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"` // pending, approved, rejected

	ReviewedBy   *uint      `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	RejectReason *string    `json:"reject_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (optional)
	User     *User `json:"user,omitempty"`
	Reviewer *User `json:"reviewer,omitempty"`
}

// --- Request structs ---

type PengadaanAssetCreateRequest struct {
	RequestDate string `json:"request_date"` // "YYYY-MM-DD"
	Category    string `json:"category,omitempty"`
	Name        string `json:"name"`
	Merk        string `json:"merk,omitempty"`
	Spec        string `json:"spec,omitempty"`
	Qty         int    `json:"qty"`
	Unit        string `json:"unit,omitempty"`
	Image       string `json:"image,omitempty"`
	Remark      string `json:"remark,omitempty"`
	Priority    string `json:"priority,omitempty"` // default: normal
}

type PengadaanAssetUpdateRequest struct {
	RequestDate string `json:"request_date,omitempty"`
	Category    string `json:"category,omitempty"`
	Name        string `json:"name,omitempty"`
	Merk        string `json:"merk,omitempty"`
	Spec        string `json:"spec,omitempty"`
	Qty         int    `json:"qty,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Image       string `json:"image,omitempty"`
	Remark      string `json:"remark,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

type PengadaanAssetReviewRequest struct {
	Status       string `json:"status"`                  // approved, rejected
	RejectReason string `json:"reject_reason,omitempty"` // wajib jika rejected
}
