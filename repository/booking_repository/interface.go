package bookingrepository

import "eventro2/models"

//go:generate mockgen -destination=../../mocks/booking_repository_mock.go -package=mocks -source=interface.go
type BookingRepository interface {
	Create(booking *models.Booking) error
	GetByID(id string) (*models.Booking, error)
	List() ([]models.Booking, error)
	ListByUser(userID string) ([]models.Booking, error)
	ListByShow(showID string) ([]models.Booking, error)
	Update(booking *models.Booking) error
	Delete(id string) error
	Find(bookingID, userID, showID string) ([]models.Booking, error)
}
