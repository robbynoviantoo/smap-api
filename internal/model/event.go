package model

import "time"

type Event struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Start  *time.Time `json:"start"`
	End    *time.Time `json:"end"`
	AllDay bool   `json:"all_day"`
	Color  string `json:"color,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventCreateRequest struct {
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end,omitempty"`
	AllDay *bool  `json:"all_day,omitempty"`
	Color  string `json:"color,omitempty"`
}

type EventUpdateRequest struct {
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end,omitempty"`
	AllDay *bool  `json:"all_day,omitempty"`
	Color  string `json:"color,omitempty"`
}
