package models

type Venue struct {
	ID                   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name                 string `gorm:"type:text;not null"`
	HostID               string `gorm:"type:uuid;not null;index"`
	Host                 User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	City                 string `gorm:"type:text;not null"`
	State                string `gorm:"type:text;not null"`
	IsBlocked            bool   `gorm:"default:false"`
	IsSeatLayoutRequired bool   `gorm:"default:false"`
}

type VenueResponse struct {
	ID                   string
	Name                 string
	HostID               string
	City                 string
	State                string
	IsSeatLayoutRequired bool
	IsBlocked            bool
}
