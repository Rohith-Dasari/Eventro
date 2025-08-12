package bookingrepository

import "eventro2/models"

type BookingStorageI interface {
	SaveBookings(bookings []models.Booking) error
	AddBooking(booking models.Booking) error
	GetBookings() ([]models.Booking, error)
}
