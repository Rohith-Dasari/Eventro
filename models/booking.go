package models

type Booking struct {
	BookingID         int      `json:"booking_id"`
	UserID            string   `json:"user_id"`
	ShowID            string   `json:"event_id"`
	TimeBooked        string   `json:"time_booked"`
	NumTickets        int      `json:"num_tickets"`
	TotalBookingPrice float64  `json:"total_booking_price"`
	Seats             []string `json:"seats"`
}
