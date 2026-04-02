package model

import "time"

type Asset struct {
	ID               uint       `json:"id"`
	Name             string     `json:"name"`
	Image            string     `json:"image,omitempty"`
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
	Remarks          string     `json:"remarks,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (optional, jangan selalu dipakai di response)
	Maintenances []AssetMaintenance `json:"maintenances,omitempty"`
}

type AssetCreateRequest struct {
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
	Status           string     `json:"status,omitempty"`
	AvailableStatus  string     `json:"available_status,omitempty"`
	LastMaintenance  *time.Time `json:"last_maintenance,omitempty"`
	NextMaintenance  *time.Time `json:"next_maintenance,omitempty"`
	Remarks          string     `json:"remarks,omitempty"`
}
type AssetUpdateRequest struct {
	ID               uint       `json:"id"`
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
	Status           string     `json:"status,omitempty"`
	AvailableStatus  string     `json:"available_status,omitempty"`
	LastMaintenance  *time.Time `json:"last_maintenance,omitempty"`
	NextMaintenance  *time.Time `json:"next_maintenance,omitempty"`
	Remarks          string     `json:"remarks,omitempty"`
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

// AssetDelete adalah record backup asset yang telah dihapus.
type AssetDelete struct {
	ID              uint       `json:"id"`
	AssetID         uint       `json:"asset_id"`
	Name            string     `json:"name"`
	Image           string     `json:"image,omitempty"`
	AssetCode       string     `json:"asset,omitempty"`
	NoAsset         string     `json:"no_asset,omitempty"`
	Location        string     `json:"location,omitempty"`
	Building        string     `json:"building,omitempty"`
	Category        string     `json:"category,omitempty"`
	SubCategory     string     `json:"sub_category,omitempty"`
	Merk            string     `json:"merk,omitempty"`
	Size            string     `json:"size,omitempty"`
	Unit            string     `json:"unit,omitempty"`
	Status          string     `json:"status,omitempty"`
	AvailableStatus string     `json:"available_status,omitempty"`
	LastMaintenance *time.Time `json:"last_maintenance,omitempty"`
	NextMaintenance *time.Time `json:"next_maintenance,omitempty"`
	Remarks         string     `json:"remarks,omitempty"`
	AssetCreatedAt  *time.Time `json:"asset_created_at,omitempty"`
	AssetUpdatedAt  *time.Time `json:"asset_updated_at,omitempty"`
	DeletedBy       *uint      `json:"deleted_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}