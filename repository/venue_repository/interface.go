package venuerepository

import "eventro2/models"

type VenueRepository interface {
	Create(venue *models.Venue) error
	GetByID(id string) (*models.Venue, error)
	List() ([]models.Venue, error)
	ListByHost(hostID string) ([]models.Venue, error)
	ListByCity(city string) ([]models.Venue, error)
	Update(venue *models.Venue) error
	Delete(id string) error
}
