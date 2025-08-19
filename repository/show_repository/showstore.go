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
