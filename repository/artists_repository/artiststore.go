package artistrepository

import (
	"eventro2/models"
	"strings"

	"gorm.io/gorm"
)

type ArtistRepositoryPG struct {
	db *gorm.DB
}

func NewArtistRepositoryPG(db *gorm.DB) *ArtistRepositoryPG {
	return &ArtistRepositoryPG{db: db}
}

func (r *ArtistRepositoryPG) Create(artist *models.Artist) error {
	return r.db.Create(artist).Error
}

func (r *ArtistRepositoryPG) GetByID(id string) (*models.Artist, error) {
	var artist models.Artist
	if err := r.db.First(&artist, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

func (r *ArtistRepositoryPG) List() ([]models.Artist, error) {
	var artists []models.Artist
	if err := r.db.Find(&artists).Error; err != nil {
		return nil, err
	}
	return artists, nil
}

// get by event id
func (r *ArtistRepositoryPG) ListByEvent(eventID string) ([]models.Artist, error) {
	var artists []models.Artist
	if err := r.db.Joins("JOIN event_artists ea ON ea.artist_id = artists.id").
		Where("ea.event_id = ?", eventID).
		Find(&artists).Error; err != nil {
		return nil, err
	}
	return artists, nil
}

// update artist-need to be implemented in admin
func (r *ArtistRepositoryPG) Update(artist *models.Artist) error {
	return r.db.Save(artist).Error
}

// delete-admin artist moderation
func (r *ArtistRepositoryPG) Delete(id string) error {
	return r.db.Delete(&models.Artist{}, "id = ?", id).Error
}

func (r *ArtistRepositoryPG) GetEventsByArtistID(artistID string) ([]models.Event, error) {
	var events []models.Event

	err := r.db.
		Joins("JOIN event_artists ea ON ea.event_id = events.id").
		Where("ea.artist_id = ?", artistID).
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *ArtistRepositoryPG) SearchByName(name string) ([]models.Artist, error) {
	var artists []models.Artist
	likePattern := "%" + strings.ToLower(name) + "%"
	err := r.db.
		Where("LOWER(name) LIKE ?", likePattern).
		Find(&artists).Error
	if err != nil {
		return nil, err
	}
	return artists, nil
}
