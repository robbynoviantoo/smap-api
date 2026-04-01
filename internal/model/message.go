package model

import "time"

type Message struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SenderID    uint      `json:"sender_id"`
	RecipientID uint      `json:"recipient_id"`
	Content     string    `json:"content"`
	IsRead      bool      `gorm:"default:false" json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}