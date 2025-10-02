package showrepository

import (
	"eventro2/models"

	"gorm.io/gorm"
)

type ShowRepositoryPG struct {
	db *gorm.DB
}

func NewShowRepositoryPG(db *gorm.DB) *ShowRepositoryPG {
	return &ShowRepositoryPG{db: db}
}

// Create a new show
func (r *ShowRepositoryPG) Create(show *models.Show) error {
	return r.db.Create(show).Error
}

// Get show by ID
func (r *ShowRepositoryPG) GetByID(id string) (*models.Show, error) {
	var show models.Show
	if err := r.db.First(&show, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &show, nil
}

// List all shows
func (r *ShowRepositoryPG) List() ([]models.Show, error) {
	var shows []models.Show
	if err := r.db.Find(&shows).Error; err != nil {
		return nil, err
	}
	return shows, nil
}

// List shows for a given Event
func (r *ShowRepositoryPG) ListByEvent(eventID string) ([]models.Show, error) {
	var shows []models.Show
	if err := r.db.Where("event_id = ?", eventID).Find(&shows).Error; err != nil {
		return nil, err
	}
	return shows, nil
}

func (r *ShowRepositoryPG) Update(show *models.Show) error {
	return r.db.Save(show).Error
}

func (r *ShowRepositoryPG) Delete(id string) error {
	return r.db.Delete(&models.Show{}, "id = ?", id).Error
}

func (r *ShowRepositoryPG) Find(filter models.ShowFilter) ([]models.Show, error) {
	var shows []models.Show
	query := r.db.Model(&models.Show{})

	if filter.ShowID != "" {
		query = query.Where("id = ?", filter.ShowID)
	}
	if filter.EventID != "" {
		query = query.Where("event_id = ?", filter.EventID)
	}
	if filter.HostID != "" {
		query = query.Where("host_id = ?", filter.HostID)
	}
	if filter.VenueID != "" {
		query = query.Where("venue_id = ?", filter.VenueID)
	}

	if err := query.Find(&shows).Error; err != nil {
		return nil, err
	}
	return shows, nil
}
