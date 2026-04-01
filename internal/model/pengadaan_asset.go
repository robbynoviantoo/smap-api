package model

import (
    "time"
)

type PengadaanAsset struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UserID      uint      `gorm:"not null" json:"user_id"`
    RequestDate time.Time `json:"request_date"`

    Category string `json:"category"`
    Name     string `json:"name"`
    Merk     string `json:"merk"`
    Spec     string `json:"spec"`

    Qty  int    `json:"qty"`
    Unit string `json:"unit"`

    Image  string `json:"image"`
    Remark string `json:"remark"`

    Priority  string `json:"priority"`
    UpdatedBy *uint  `json:"updated_by"`

    // Relation
    User User `gorm:"foreignKey:UserID" json:"user"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}