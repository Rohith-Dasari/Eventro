package venueservice

import (
	"context"
	"eventro2/models"
)

//go:generate mockgen -destination=../../mocks/venue_service_mock.go -package=mocks -source=interface.go
type VenueServiceInterface interface {
	CreateVenue(ctx context.Context, hostID, name, city, state string, isSeatLayoutRequired bool) (models.VenueResponse, error)
	UpdateVenue(ctx context.Context, venueID, userID, userRole string, update models.UpdateVenueData) (models.VenueResponse, error)
	DeleteVenue(ctx context.Context, venueID, userID, userRole string) error
	BrowseVenues(ctx context.Context, filter models.VenueFilter) ([]models.VenueResponse, error)
}
