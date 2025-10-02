package artistrepository

import "eventro2/models"

type ArtistRepository interface {
	Create(artist *models.Artist) error
	GetByID(id string) (*models.Artist, error)
	List() ([]models.Artist, error)
	ListByEvent(eventID string) ([]models.Artist, error)
	Update(artist *models.Artist) error
	Delete(id string) error
	GetEventsByArtistID(artistID string) ([]models.Event, error)
	SearchByName(name string) ([]models.Artist, error)
}
