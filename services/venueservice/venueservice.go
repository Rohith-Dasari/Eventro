package venueservice

import (
	"context"
	"eventro2/models"
	venuerepository "eventro2/repository/venue_repository"
	"fmt"

	"github.com/google/uuid"
)

type VenueService struct {
	VenueRepo venuerepository.VenueRepository
}

func NewVenueService(repo venuerepository.VenueRepository) VenueService {
	return VenueService{VenueRepo: repo}
}

func (vs *VenueService) CreateVenue(ctx context.Context, hostID, name, city, state string, isSeatLayoutRequired bool) (models.VenueResponse, error) {
	venueID := uuid.New().String()

	venue := models.Venue{
		ID:                   venueID,
		HostID:               hostID,
		Name:                 name,
		City:                 city,
		State:                state,
		IsSeatLayoutRequired: isSeatLayoutRequired,
	}

	if err := vs.VenueRepo.Create(&venue); err != nil {
		return models.VenueResponse{}, fmt.Errorf("failed to create venue: %w", err)
	}
	venueDTO := models.VenueResponse{
		ID:                   venueID,
		HostID:               hostID,
		Name:                 name,
		City:                 city,
		State:                state,
		IsSeatLayoutRequired: isSeatLayoutRequired,
	}

	return venueDTO, nil
}

func (s *VenueService) UpdateVenue(ctx context.Context, venueID, userID, userRole string, update models.UpdateVenueData) (models.VenueResponse, error) {
	venue, err := s.VenueRepo.GetByID(venueID)
	if err != nil {
		return models.VenueResponse{}, err
	}

	if venue.HostID != userID && update.IsBlocked == nil {
		return models.VenueResponse{}, fmt.Errorf("forbidden: cannot update another user's venue")
	}

	if update.Name != nil {
		venue.Name = *update.Name
	}
	if update.City != nil {
		venue.City = *update.City
	}
	if update.State != nil {
		venue.State = *update.State
	}
	if update.IsSeatLayoutRequired != nil {
		venue.IsSeatLayoutRequired = *update.IsSeatLayoutRequired
	}

	if update.IsBlocked != nil {
		if venue.HostID == userID || userRole == "admin" {
			venue.IsBlocked = *update.IsBlocked
		} else {
			return models.VenueResponse{}, fmt.Errorf("forbidden: only host or admin can block/unblock")
		}
	}

	if err := s.VenueRepo.Update(venue); err != nil {
		return models.VenueResponse{}, err
	}
	venueDTO := models.VenueResponse{
		ID:                   venue.ID,
		HostID:               venue.HostID,
		Name:                 venue.Name,
		City:                 venue.City,
		State:                venue.State,
		IsSeatLayoutRequired: venue.IsSeatLayoutRequired,
		IsBlocked:            venue.IsBlocked,
	}

	return venueDTO, nil
}

func (s *VenueService) DeleteVenue(ctx context.Context, venueID, userID, userRole string) error {
	venue, err := s.VenueRepo.GetByID(venueID)
	if err != nil {
		return err
	}

	if venue.HostID != userID && userRole != "admin" {
		return fmt.Errorf("forbidden: cannot delete this venue")
	}

	if err := s.VenueRepo.Delete(venueID); err != nil {
		return err
	}

	return nil
}

func (s *VenueService) BrowseVenues(ctx context.Context, filter models.VenueFilter) ([]models.VenueResponse, error) {
	venue, err := s.VenueRepo.Find(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch venues: %w", err)
	}

	venueDTO := make([]models.VenueResponse, len(venue))
	for i, v := range venue {
		venueDTO[i] = models.VenueResponse{
			ID:                   v.ID,
			HostID:               v.HostID,
			Name:                 v.Name,
			City:                 v.City,
			State:                v.State,
			IsSeatLayoutRequired: v.IsSeatLayoutRequired,
			IsBlocked:            v.IsBlocked,
		}
	}
	return venueDTO, nil
}
