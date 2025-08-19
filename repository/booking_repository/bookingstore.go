package bookingrepository

import (
	"eventro2/models"

	"gorm.io/gorm"
)

type BookingRepositoryPG struct {
	db *gorm.DB
}

func NewBookingRepositoryPG(db *gorm.DB) *BookingRepositoryPG {
	return &BookingRepositoryPG{db: db}
}

func (r *BookingRepositoryPG) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepositoryPG) GetByID(id string) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.First(&booking, "booking_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepositoryPG) List() ([]models.Booking, error) {
	var bookings []models.Booking
	if err := r.db.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepositoryPG) ListByUser(userID string) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := r.db.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepositoryPG) ListByShow(showID string) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := r.db.Where("show_id = ?", showID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepositoryPG) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *BookingRepositoryPG) Delete(id string) error {
	return r.db.Delete(&models.Booking{}, "booking_id = ?", id).Error
}
