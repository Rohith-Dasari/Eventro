package userrepository

import "eventro2/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	GetBlockedUsers() ([]models.User, error)
}
