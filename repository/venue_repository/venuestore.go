package venuerepository

import (
	"eventro2/models"

	"gorm.io/gorm"
)

type VenueRepositoryPG struct {
	db *gorm.DB
}

func NewVenueRepositoryPG(db *gorm.DB) *VenueRepositoryPG {
	return &VenueRepositoryPG{db: db}
}

// Create venue
func (r *VenueRepositoryPG) Create(venue *models.Venue) error {
	return r.db.Create(venue).Error
}

// Get venue by ID
func (r *VenueRepositoryPG) GetByID(id string) (*models.Venue, error) {
	var venue models.Venue
	if err := r.db.First(&venue, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &venue, nil
}

// List all venues
func (r *VenueRepositoryPG) List() ([]models.Venue, error) {
	var venues []models.Venue
	if err := r.db.Find(&venues).Error; err != nil {
		return nil, err
	}
	return venues, nil
}

// List venues by host
func (r *VenueRepositoryPG) ListByHost(hostID string) ([]models.Venue, error) {
	var venues []models.Venue
	if err := r.db.Where("host_id = ?", hostID).Find(&venues).Error; err != nil {
		return nil, err
	}
	return venues, nil
}

// List venues by city
func (r *VenueRepositoryPG) ListByCity(city string) ([]models.Venue, error) {
	var venues []models.Venue
	if err := r.db.Where("city ILIKE ?", city).Find(&venues).Error; err != nil {
		return nil, err
	}
	return venues, nil
}

// Update venue
func (r *VenueRepositoryPG) Update(venue *models.Venue) error {
	return r.db.Save(venue).Error
}

// Delete venue
func (r *VenueRepositoryPG) Delete(id string) error {
	return r.db.Delete(&models.Venue{}, "id = ?", id).Error
}
