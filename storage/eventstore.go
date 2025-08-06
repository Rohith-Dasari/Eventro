package storage

import (
	"encoding/json"
	"eventro/models"
	"log"
	"os"
)

func LoadEvents() []models.Event {
	//read json
	data, err := os.ReadFile("data/events.json")
	if err != nil {
		log.Fatalf("failed to read file %v", err)
	}
	//unmarshal into user class
	var events []models.Event
	if err := json.Unmarshal(data, &events); err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}

	return events
}
