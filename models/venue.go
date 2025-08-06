package models

type Venue struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	HostID               string `json:"host_id"`
	City                 string `json:"city"`
	State                string `json:"state"`
	IsSeatLayoutRequired bool   `json:"is_seat_layout_required"`
}
