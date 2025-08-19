package models

import (
	"time"

	"github.com/lib/pq"
)

type Show struct {
	ID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	HostID string `gorm:"type:uuid;not null;index"`
	Host   User   `gorm:"foreignKey:HostID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	VenueID string `gorm:"type:uuid;not null;index"`
	Venue   Venue  `gorm:"foreignKey:VenueID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	EventID string `gorm:"type:uuid;not null;index"`
	Event   Event  `gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	IsBlocked   bool           `gorm:"default:false"`
	Price       float64        `gorm:"not null"`
	ShowDate    time.Time      `gorm:"not null"`
	ShowTime    string         `gorm:"type:varchar(5);not null"`
	BookedSeats pq.StringArray `gorm:"type:text[]"`
}
