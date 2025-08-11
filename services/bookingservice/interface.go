package bookingservice

import (
	"context"
	"eventro2/models"
)

type BookingServiceI interface {
	ViewBookingHistory(ctx context.Context)
	MakeBooking(ctx context.Context, userID string, showID string)
	AddBooking(booking models.Booking)
}
