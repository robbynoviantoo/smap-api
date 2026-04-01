package model

import "time"

type Asset struct {
	ID               uint       `json:"id"`
	Name             string     `json:"name"`
	Image            string     `json:"image"`
	AssetCode        string     `json:"asset"`        // dari field 'asset'
	NoAsset          string     `json:"no_asset"`
	Location         string     `json:"location"`
	Building         string     `json:"building"`
	Category         string     `json:"category"`
	SubCategory      string     `json:"sub_category"`
	Merk             string     `json:"merk"`
	Size             string     `json:"size"`
	Unit             string     `json:"unit"`
	Status           string     `json:"status"`
	AvailableStatus  string     `json:"available_status"`
	LastMaintenance  *time.Time `json:"last_maintenance"`
	NextMaintenance  *time.Time `json:"next_maintenance"`
	Remarks          string     `json:"remarks"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (optional, jangan selalu dipakai di response)
	Maintenances []AssetMaintenance `json:"maintenances,omitempty"`
}

type AssetFilter struct {
	Name            string
	AssetCode       string
	NoAsset         string
	Location        string
	Building        string
	Category        string
	Status          string
	AvailableStatus string
}