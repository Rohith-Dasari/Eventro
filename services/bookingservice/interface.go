package bookingservice

import "context"

type BookingServiceI interface {
	ViewBookingHistory(ctx context.Context)
	MakeBooking(ctx context.Context, userID string, showID string)
}
