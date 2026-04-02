package model

import "time"

// AssetBorrowPending adalah approval queue untuk pengajuan peminjaman asset.
type AssetBorrowPending struct {
	ID           uint       `json:"id"`
	AssetID      uint       `json:"asset_id"`
	UserID       uint       `json:"user_id"`
	Remark       string     `json:"remark"`
	Status       string     `json:"status"` // pending, approved, rejected
	ReviewedBy   *uint      `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	RejectReason *string    `json:"reject_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (optional)
	Asset    *Asset `json:"asset,omitempty"`
	User     *User  `json:"user,omitempty"`
	Reviewer *User  `json:"reviewer,omitempty"`
}

// AssetTransaction adalah append-only log untuk borrow dan return.
type AssetTransaction struct {
	ID         uint      `json:"id"`
	AssetID    uint      `json:"asset_id"`
	BorrowerID uint      `json:"borrower_id"`
	Type       string    `json:"type"` // borrow, return
	Remark     string    `json:"remark"`
	ActionAt   time.Time `json:"action_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (optional)
	Asset    *Asset `json:"asset,omitempty"`
	Borrower *User  `json:"borrower,omitempty"`
}

// --- Request structs ---

type BorrowRequest struct {
	Remark string `json:"remark"`
}

type ReturnRequest struct {
	Remark string `json:"remark"`
}

type BorrowPendingReviewRequest struct {
	Status       string `json:"status"`                  // approved, rejected
	RejectReason string `json:"reject_reason,omitempty"` // wajib jika rejected
}
