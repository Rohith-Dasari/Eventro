package bookingrepository

import (
	"encoding/json"
	"eventro2/config"
	"eventro2/models"
	"fmt"
	"log"
	"os"
)

type BookingRepository struct {
	Bookings []models.Booking
}

// func (bs *BookingRepository) LoadBookings() []models.Booking {
// 	//read json
// 	data, err := os.ReadFile(config.BookingsFile)
// 	if err != nil {
// 		log.Fatalf("failed to read file %v", err)
// 	}
// 	//unmarshal into booking class
// 	var bookings []models.Booking
// 	if err := json.Unmarshal(data, &bookings); err != nil {
// 		log.Fatalf("failed to marshal: %v", err)
// 	}

// 	return bookings
// }

func (bs *BookingRepository) SaveBookings(bookings []models.Booking) error {
	data, err := json.MarshalIndent(bookings, "", "")
	if err != nil {
		return fmt.Errorf("failed to save in bookings file: %w", err)
	}
	err = os.WriteFile(config.BookingsFile, data, 0644)
	return err
}

func (br *BookingRepository) AddBooking(booking models.Booking) error {
	br.Bookings = append(br.Bookings, booking)
	return br.SaveBookings(br.Bookings)
}

func NewBoookingStore() *BookingRepository {
	//read json
	data, err := os.ReadFile(config.BookingsFile)
	if err != nil {
		log.Fatalf("failed to read file %v", err)
	}
	//unmarshal into booking class
	var bookings []models.Booking
	if err := json.Unmarshal(data, &bookings); err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}
	return &BookingRepository{bookings}
}
