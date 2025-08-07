package venueservice

import "context"

type VenueServiceInterface interface {
	AddVenue(ctx context.Context)
	BrowseHostVenues(ctx context.Context)
	RemoveVenue(ctx context.Context)
}
