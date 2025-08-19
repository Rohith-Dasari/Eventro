package models

type Artist struct {
	ID   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string `gorm:"type:text;not null"`
	Bio  string `gorm:"type:text"`
}
