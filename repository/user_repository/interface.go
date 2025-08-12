package userrepository

import "eventro2/models"

type UserStorageI interface {
	GetUsers() ([]models.User, error)
	SaveUsers(users []models.User) error
	AddUser(user models.User) error
	UserExists(users []models.User, email string) bool
}
