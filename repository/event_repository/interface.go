package eventsrepository

import "eventro2/models"

type EventStorageI interface {
	SaveEvents(events []models.Event) error
	GetEvents() ([]models.Event, error)
}
