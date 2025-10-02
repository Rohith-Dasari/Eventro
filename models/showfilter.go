package models

type ShowFilter struct {
	ShowID  string `json:"showId"`
	EventID string `json:"eventId"`
	HostID  string `json:"hostId"`
	VenueID string `json:"venueId"`
}
