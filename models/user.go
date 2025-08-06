package models

type Role string

const (
	Admin    Role = "Admin"
	Host     Role = "Host"
	Customer Role = "Customer"
)

type User struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Role        Role   `json:"role"`
	IsBlocked   bool   `json:"is_blocked"`
}
