package venuerepository

import "eventro2/models"

type VenueStorageI interface {
	GetVenues() ([]models.Venue, error)
	SaveVenues(venues []models.Venue) error
}
