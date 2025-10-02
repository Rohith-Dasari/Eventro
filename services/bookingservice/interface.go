package bookingservice

import (
	"context"
	"eventro2/models"
)

//go:generate mockgen -destination=../../mocks/booking_service_mock.go -package=mocks -source=interface.go
type BookingServiceInterface interface {
	AddBooking(ctx context.Context, userID string, showID string, seats []string) (*models.BookingResponse, error)
	BrowseBookings(ctx context.Context, bookingID string, userID string, showID string) ([]models.BookingResponse, error)
}
