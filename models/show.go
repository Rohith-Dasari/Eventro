package models

type Show struct {
	ID          string   `json:"id"`
	HostID      string   `json:"host_id"`
	VenueID     string   `json:"venue_id"`
	EventID     string   `json:"event_id"`
	CreatedAt   string   `json:"created_at"`
	IsBlocked   bool     `json:"is_blocked"`
	Price       float64  `json:"price"`
	ShowDate    string   `json:"show_date"`
	ShowTime    string   `json:"show_time"`
	BookedSeats []string `json:"booked_seats"`
}
