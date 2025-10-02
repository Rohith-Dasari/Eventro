package bookingservice

import (
	"context"
	"errors"
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	"fmt"
	"regexp"
	"strings"
)

type BookingService struct {
	BookingRepo bookingrepository.BookingRepository
	ShowRepo    showrepository.ShowRepository
	venueRepo   venuerepository.VenueRepository
	eventRepo   eventsrepository.EventRepository
}

func NewBookingService(bRepo bookingrepository.BookingRepository, sRepo showrepository.ShowRepository, vRepo venuerepository.VenueRepository, eRepo eventsrepository.EventRepository) BookingService {
	return BookingService{
		BookingRepo: bRepo,
		ShowRepo:    sRepo,
		venueRepo:   vRepo,
		eventRepo:   eRepo,
	}
}

func (bs *BookingService) AddBooking(
	ctx context.Context,
	userID string,
	showID string,
	requestedSeats []string,
) (*models.BookingResponse, error) {
	show, err := bs.ShowRepo.GetByID(showID)
	if err != nil {
		return nil, fmt.Errorf("show not found: %w", err)
	}
	if show.IsBlocked {
		return nil, errors.New("cannot book tickets for a blocked show")
	}

	booked := make(map[string]bool)
	for _, s := range show.BookedSeats {
		booked[s] = true
	}

	for _, seat := range requestedSeats {
		if booked[seat] {
			return nil, fmt.Errorf("seat %s is already booked", seat)
		}
		if !bs.isValidTicket(seat, show.BookedSeats) {
			return nil, fmt.Errorf("seat %s is not valid", seat)
		}
	}

	numTickets := len(requestedSeats)
	totalPrice := float64(numTickets) * show.Price

	show.BookedSeats = append(show.BookedSeats, requestedSeats...)
	if err := bs.ShowRepo.Update(show); err != nil {
		return nil, fmt.Errorf("failed to update show bookings: %w", err)
	}

	newBooking := &models.Booking{
		UserID:            userID,
		ShowID:            showID,
		NumTickets:        numTickets,
		TotalBookingPrice: totalPrice,
		Seats:             requestedSeats,
	}

	if err := bs.BookingRepo.Create(newBooking); err != nil {
		return nil, fmt.Errorf("error creating booking: %w", err)
	}
	bookingDTO := models.BookingResponse{
		BookingID:         newBooking.BookingID,
		UserID:            newBooking.UserID,
		ShowID:            newBooking.ShowID,
		TimeBooked:        newBooking.TimeBooked,
		NumTickets:        newBooking.NumTickets,
		TotalBookingPrice: newBooking.TotalBookingPrice,
		Seats:             newBooking.Seats,
	}

	return &bookingDTO, nil
}

func (bs *BookingService) isValidTicket(userTicket string, bookedTickets []string) bool {
	userTicket = strings.ToUpper(userTicket)

	matched, err := regexp.MatchString(`^[A-J](10|[1-9])$`, userTicket)
	if err != nil || !matched {
		return false
	}

	for _, ticket := range bookedTickets {
		if strings.ToUpper(ticket) == userTicket {
			return false
		}
	}
	return true
}
func (s *BookingService) BrowseBookings(ctx context.Context, bookingID, userID, showID string) ([]models.BookingResponse, error) {
	bookings, err := s.BookingRepo.Find(bookingID, userID, showID)
	if err != nil {
		return nil, fmt.Errorf("error fetching bookings: %w", err)
	}
	bookingDTO := make([]models.BookingResponse, len(bookings))
	for i, booking := range bookings {
		bookingDTO[i] = models.BookingResponse{
			BookingID:         booking.BookingID,
			UserID:            booking.UserID,
			ShowID:            booking.ShowID,
			TimeBooked:        booking.TimeBooked,
			NumTickets:        booking.NumTickets,
			TotalBookingPrice: booking.TotalBookingPrice,
			Seats:             booking.Seats,
		}
	}
	return bookingDTO, nil
}
