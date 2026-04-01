package model

import "time"

type FcmToken struct {
    ID         uint   `gorm:"primaryKey" json:"id"`
    UserID     uint   `gorm:"not null;index" json:"user_id"`
    Token      string `gorm:"type:text;not null" json:"token"`
    DeviceType string `json:"device_type"`

    // Relation
    User User `gorm:"foreignKey:UserID" json:"user"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}