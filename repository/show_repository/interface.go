package showrepository

import "eventro2/models"

//go:generate mockgen -destination=../../mocks/show_repository_mock.go -package=mocks -source=interface.go
type ShowRepository interface {
	Create(show *models.Show) error
	GetByID(id string) (*models.Show, error)
	List() ([]models.Show, error)
	ListByEvent(eventID string) ([]models.Show, error)
	Update(show *models.Show) error
	Delete(id string) error
	Find(filter models.ShowFilter) ([]models.Show, error)
}
