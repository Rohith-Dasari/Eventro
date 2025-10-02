package venuerepository

import "eventro2/models"

//go:generate mockgen -destination=../../mocks/venue_repository_mock.go -package=mocks -source=interface.go
type VenueRepository interface {
	Create(venue *models.Venue) error
	GetByID(id string) (*models.Venue, error)
	List() ([]models.Venue, error)
	ListByHost(hostID string) ([]models.Venue, error)
	ListByCity(city string) ([]models.Venue, error)
	Update(venue *models.Venue) error
	Delete(id string) error
	Find(filter models.VenueFilter) ([]models.Venue, error)
}
