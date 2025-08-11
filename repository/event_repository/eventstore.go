package eventsrepository

import (
	"encoding/json"
	"eventro2/config"
	"eventro2/models"
	"fmt"
	"log"
	"os"
)

type EventRepository struct {
	Events []models.Event
}

// func LoadEvents() []models.Event {
// 	//read json
// 	data, err := os.ReadFile("data/events.json")
// 	if err != nil {
// 		log.Fatalf("failed to read file %v", err)
// 	}
// 	//unmarshal into user class
// 	var events []models.Event
// 	if err := json.Unmarshal(data, &events); err != nil {
// 		log.Fatalf("failed to marshal: %v", err)
// 	}

// 	return events
// }

func (*EventRepository) SaveEvents(events []models.Event) error {
	data, err := json.MarshalIndent(events, "", "")
	if err != nil {
		return fmt.Errorf("failed to serialise %w", err)
	}
	err = os.WriteFile(config.EventsFile, data, 0644)
	return err
}
func NewEventRepository() *EventRepository {
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
	return &EventRepository{events}
}

func (er *EventRepository) GetEvents() ([]models.Event, error) {
	data, err := os.ReadFile(config.EventsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read events file: %w", err)
	}

	var events []models.Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events: %w", err)
	}

	return events, nil
}
