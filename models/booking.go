package models

import (
	"time"

	"github.com/lib/pq"
)

type Booking struct {
	BookingID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID            string         `gorm:"type:uuid;not null"`
	User              User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ShowID            string         `gorm:"type:uuid;not null;index"`
	Show              Show           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TimeBooked        time.Time      `gorm:"autoCreateTime"`
	NumTickets        int            `gorm:"not null"`
	TotalBookingPrice float64        `gorm:"not null"`
	Seats             pq.StringArray `gorm:"type:text[]"`
}

type BookingResponse struct {
	BookingID         string    `json:"booking_id"`
	UserID            string    `json:"user_id"`
	ShowID            string    `json:"show_id"`
	TimeBooked        time.Time `json:"time_booked"`
	NumTickets        int       `json:"num_tickets"`
	TotalBookingPrice float64   `json:"total_booking_price"`
	Seats             []string  `json:"seats"`
}
