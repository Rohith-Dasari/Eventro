package eventrepository

import "eventro2/models"

type EventRepository interface {
	Create(event *models.Event) error
	GetByID(id string) (*models.Event, error)
	List() ([]models.Event, error)
	Update(event *models.Event) error
	Delete(id string) error
	AddEventArtist(ea *models.EventArtist) error
	GetArtistsByEventID(eventID string) ([]models.Artist, error)
	GetEventsByCity(city string) ([]models.Event, error)
}
