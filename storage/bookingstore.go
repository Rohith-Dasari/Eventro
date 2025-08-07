package storage

import (
	"encoding/json"
	"eventro/config"
	"eventro/models"
	"log"
	"os"
)

func LoadBookings() []models.Booking {
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

	return bookings
}

func SaveBookings(bookings []models.Booking) error {
	data, err := json.MarshalIndent(bookings, "", "")
	if err != nil {
		return err
	}
	err = os.WriteFile(config.BookingsFile, data, 0644)
	return err
}
