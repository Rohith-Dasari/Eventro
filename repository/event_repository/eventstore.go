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
func (r *EventRepositoryPG) GetFilteredEvents(filter models.EventFilter) ([]models.EventResponse, error) {
	var events []models.EventResponse

	query := r.db.Table("events").
		Select(`
        events.id,
        events.name,
        events.description,
        events.duration,
        events.category,
        events.is_blocked,
        COALESCE(array_agg(DISTINCT ea.artist_id) FILTER (WHERE ea.artist_id IS NOT NULL), '{}') as artist_ids,
        COALESCE(array_agg(DISTINCT a.name) FILTER (WHERE a.name IS NOT NULL), '{}') as artist_names
    `).
		Joins("LEFT JOIN event_artists ea ON ea.event_id = events.id").
		Joins("LEFT JOIN artists a ON ea.artist_id = a.id").
		Group("events.id")

	if filter.Name != "" {
		query = query.Where("LOWER(events.name) LIKE ?", "%"+strings.ToLower(filter.Name)+"%")
	}
	if filter.Category != "" {
		query = query.Where("events.category = ?", filter.Category)
	}
	if filter.Location != "" {
		query = query.Joins("JOIN shows s ON s.event_id = events.id").
			Joins("JOIN venues v ON s.venue_id = v.id").
			Where("LOWER(v.city) = ?", strings.ToLower(filter.Location))
	}
	if filter.IsBlocked != nil {
		query = query.Where("events.is_blocked = ?", *filter.IsBlocked)
	}
	if filter.ArtistName != "" {
		query = query.Where("LOWER(a.name) LIKE ?", "%"+strings.ToLower(filter.ArtistName)+"%")
	}
	if filter.EventID != "" {
		query = query.Where("events.id = ?", filter.EventID)
	}

	if err := query.Scan(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
