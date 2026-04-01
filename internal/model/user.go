package model

import (
    "time"
)

type User struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    Email         string    `gorm:"unique;not null" json:"email"`
    Password      string    `gorm:"not null" json:"-"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    NoHandphone   string    `json:"no_handphone"`
    Negara        string    `json:"negara"`
    Kota          string    `json:"kota"`
    Kodepos       string    `json:"kodepos"`
    Bio           string    `json:"bio"`
    Image         string    `json:"image"`
    IsFirstLogin  bool      `gorm:"default:true" json:"is_first_login"`

    EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`

    // Relations
    Roles           []Role            `gorm:"many2many:user_roles;" json:"roles"`
    SentMessages    []Message         `gorm:"foreignKey:SenderID" json:"sent_messages"`
    ReceivedMessages []Message        `gorm:"foreignKey:RecipientID" json:"received_messages"`
    RequestAssets   []PengadaanAsset  `json:"request_assets"`
    FcmTokens       []FcmToken        `json:"fcm_tokens"`

    CreatedAt time.Time
    UpdatedAt time.Time
}

type UserResponse struct {
	ID          uint      `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	NoHandphone string    `json:"no_handphone"`
	CreatedAt   time.Time `json:"created_at"`
    IsFirstLogin bool      `json:"is_first_login"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
    Roles []string `json:"roles"`
}