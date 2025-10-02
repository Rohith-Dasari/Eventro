package models

type Role string

const (
	Admin    Role = "Admin"
	Host     Role = "Host"
	Customer Role = "Customer"
)

type User struct {
	UserID      string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username    string `gorm:"type:text;not null"`
	Email       string `gorm:"type:text;uniqueIndex;not null"`
	PhoneNumber string `gorm:"type:text;not null"`
	Password    string `gorm:"type:text;not null"`
	Role        Role   `gorm:"type:text;not null"`
	IsBlocked   bool   `gorm:"default:false"`
}
