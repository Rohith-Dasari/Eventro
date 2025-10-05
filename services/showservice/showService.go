package showservice

import (
	"context"
	"errors"
	"eventro2/middleware"
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ShowService struct {
	ShowRepo  showrepository.ShowRepository
	VenueRepo venuerepository.VenueRepository
	BookRepo  bookingrepository.BookingRepository
	EventRepo eventsrepository.EventRepository
}

func NewShowService(
	showRepo showrepository.ShowRepository,
	venueRepo venuerepository.VenueRepository,
	bookRepo bookingrepository.BookingRepository,
	eventRepo eventsrepository.EventRepository,
) ShowService {
	return ShowService{
		ShowRepo:  showRepo,
		VenueRepo: venueRepo,
		BookRepo:  bookRepo,
		EventRepo: eventRepo,
	}
}

func (s *ShowService) UpdateShow(ctx context.Context, showID string, userID string, update models.UpdateShowData) (models.ShowResponse, error) {
	// fetch show first
	show, err := s.ShowRepo.GetByID(showID)
	if err != nil {
		return models.ShowResponse{}, err
	}

	// authorization check (only host can update)
	if show.HostID != userID {
		return models.ShowResponse{}, fmt.Errorf("forbidden: cannot update another user's show")
	}

	// apply updates if provided
	if update.Price != nil {
		show.Price = *update.Price
	}
	if update.ShowDate != nil {
		show.ShowDate = *update.ShowDate
	}
	if update.ShowTime != nil {
		show.ShowTime = *update.ShowTime
	}
	if update.IsBlocked != nil {
		show.IsBlocked = *update.IsBlocked
	}

	if err := s.ShowRepo.Update(show); err != nil {
		return models.ShowResponse{}, err
	}

	return models.ShowResponse{
		ID:          show.ID,
		HostID:      show.HostID,
		VenueID:     show.VenueID,
		EventID:     show.EventID,
		CreatedAt:   show.CreatedAt,
		IsBlocked:   show.IsBlocked,
		Price:       show.Price,
		ShowDate:    show.ShowDate,
		ShowTime:    show.ShowTime,
		BookedSeats: show.BookedSeats,
	}, nil
}

func (s *ShowService) BrowseShows(ctx context.Context, filter models.ShowFilter) ([]models.Show, error) {
	shows, err := s.ShowRepo.Find(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shows: %w", err)
	}
	return shows, nil
}

func (s *ShowService) CreateShow(ctx context.Context, eventID string, venueID string, hostID string, price float64, showDate time.Time, showTime string) (models.ShowResponse, error) {
	showID := uuid.New().String()

	show := models.Show{
		ID:          showID,
		HostID:      hostID,
		VenueID:     venueID,
		EventID:     eventID,
		IsBlocked:   false,
		Price:       price,
		ShowDate:    showDate,
		ShowTime:    showTime,
		BookedSeats: []string{},
	}
	showDTO := models.ShowResponse{
		ID:          show.ID,
		HostID:      show.HostID,
		VenueID:     show.VenueID,
		EventID:     show.EventID,
		CreatedAt:   show.CreatedAt,
		IsBlocked:   show.IsBlocked,
		Price:       show.Price,
		ShowDate:    show.ShowDate,
		ShowTime:    show.ShowTime,
		BookedSeats: show.BookedSeats,
	}

	if err := s.ShowRepo.Create(&show); err != nil {
		return models.ShowResponse{}, fmt.Errorf("failed to create show: %w", err)
	}

	return showDTO, nil
}

func (s *ShowService) DeleteShow(ctx context.Context, showID string) error {
	show, err := s.ShowRepo.GetByID(showID)
	if err != nil {
		return fmt.Errorf("show not found: %w", err)
	}

	userID, _ := middleware.GetUserID(ctx)
	role, _ := middleware.GetUserRole(ctx)

	if strings.ToLower(role) != "admin" && show.HostID != userID {
		return errors.New("unauthorized to delete this show")
	}

	if err := s.ShowRepo.Delete(showID); err != nil {
		return fmt.Errorf("failed to delete show: %w", err)
	}

	return nil
}
