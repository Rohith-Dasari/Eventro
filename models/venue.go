package models

type Venue struct {
	ID                   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name                 string `gorm:"type:text;not null"`
	HostID               string `gorm:"type:uuid;not null;index"`
	Host                 User   `gorm:"foreignKey:HostID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	City                 string `gorm:"type:text;not null"`
	State                string `gorm:"type:text;not null"`
	IsSeatLayoutRequired bool   `gorm:"default:false"`
}
