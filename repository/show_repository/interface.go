package showrepository

import "eventro2/models"

type ShowStorageI interface {
	SaveShows(shows []models.Show) error
	GetShows() ([]models.Show, error)
}
