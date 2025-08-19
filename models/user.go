package models

type Role string

const (
	Admin    Role = "Admin"
	Host     Role = "Host"
	Customer Role = "Customer"
)

type User struct {
	UserID      string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username    string    `gorm:"not null"`
	Email       string    `gorm:"uniqueIndex;not null"`
	PhoneNumber string    `gorm:"not null"`
	Password    string    `gorm:"not null"`
	Role        Role      `gorm:"type:text;not null"`
	IsBlocked   bool      `gorm:"default:false"`
	Bookings    []Booking `gorm:"foreignKey:UserID;references:UserID"`
}
