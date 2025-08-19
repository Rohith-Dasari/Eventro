package showrepository

import "eventro2/models"

type ShowRepository interface {
	Create(show *models.Show) error
	GetByID(id string) (*models.Show, error)
	List() ([]models.Show, error)
	ListByEvent(eventID string) ([]models.Show, error)
	Update(show *models.Show) error
	Delete(id string) error
}
