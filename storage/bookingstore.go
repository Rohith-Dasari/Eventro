package storage

import (
	"encoding/json"
	"eventro/models"
	"log"
	"os"
)

func LoadBookings() []models.Booking {
	//read json
	data, err := os.ReadFile("data/bookings.json")
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
