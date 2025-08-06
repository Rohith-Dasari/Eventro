package storage

import (
	"encoding/json"
	"eventro/config"
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

func SaveEvents(events []models.Event) error {
	data, err := json.MarshalIndent(events, "", "")
	if err != nil {
		return err
	}
	err = os.WriteFile(config.EventsFile, data, 0644)
	return err
}
