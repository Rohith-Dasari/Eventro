package eventrepository

import (
	"eventro2/models"
	"strings"

	"gorm.io/gorm"
)

type EventRepositoryPG struct {
	db *gorm.DB
}

func NewEventRepositoryPG(db *gorm.DB) *EventRepositoryPG {
	return &EventRepositoryPG{db: db}
}

func (r *EventRepositoryPG) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *EventRepositoryPG) GetByID(id string) (*models.Event, error) {
	var event models.Event
	if err := r.db.First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepositoryPG) List() ([]models.Event, error) {
	var events []models.Event
	if err := r.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// Update event
func (r *EventRepositoryPG) Update(event *models.Event) error {
	return r.db.Model(&models.Event{}).Where("id = ?", event.ID).Update("is_blocked", event.IsBlocked).Error
}

// Delete event
func (r *EventRepositoryPG) Delete(id string) error {
	return r.db.Delete(&models.Event{}, "id = ?", id).Error
}

func (r *EventRepositoryPG) AddEventArtist(ea *models.EventArtist) error {
	return r.db.Create(ea).Error
}

func (r *EventRepositoryPG) GetArtistsByEventID(eventID string) ([]models.Artist, error) {
	var artists []models.Artist
	err := r.db.
		Table("artists").
		Select("artists.*").
		Joins("JOIN event_artists ea ON ea.artist_id = artists.id").
		Where("ea.event_id = ?", eventID).
		Find(&artists).Error

	if err != nil {
		return nil, err
	}
	return artists, nil
}

func (r *EventRepositoryPG) GetEventsByCity(city string) ([]models.Event, error) {
	var events []models.Event

	err := r.db.
		Table("events e").
		Select("DISTINCT e.*").
		Joins("JOIN shows s ON s.event_id = e.id").
		Joins("JOIN venues v ON s.venue_id = v.id").
		Where("LOWER(v.city) = ?", strings.ToLower(city)).
		Find(&events).Error

	if err != nil {
		return nil, err
	}
	return events, nil
}
