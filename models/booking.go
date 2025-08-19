package models

import (
	"time"

	"github.com/lib/pq"
)

type Booking struct {
	BookingID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	UserID string `gorm:"type:uuid;not null;index"`
	User   User   `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ShowID string `gorm:"type:uuid;not null;index"`
	Show   Show   `gorm:"foreignKey:ShowID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	TimeBooked        time.Time      `gorm:"autoCreateTime"`
	NumTickets        int            `gorm:"not null"`
	TotalBookingPrice float64        `gorm:"not null"`
	Seats             pq.StringArray `gorm:"type:text[]"`
}
