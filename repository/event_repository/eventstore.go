package eventsrepository

import (
	"encoding/json"
	"eventro2/config"
	"eventro2/models"
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
		return err
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
