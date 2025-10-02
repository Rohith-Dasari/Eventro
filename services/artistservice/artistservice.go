package artistservice

import (
	"context"
	"errors"
	"eventro2/models"
	artistrepository "eventro2/repository/artists_repository"

	"github.com/google/uuid"
)

type Artistservice struct {
	ArtistRepo artistrepository.ArtistRepository
}

func NewArtistService(artistRepo artistrepository.ArtistRepository) Artistservice {
	return Artistservice{
		ArtistRepo: artistRepo,
	}
}

func (as *Artistservice) CreateArtist(ctx context.Context, name, bio string) (models.Artist, error) {

	if len(bio) < 12 {
		return models.Artist{}, errors.New("bio must be at least 12 characters long")
	}

	artist := models.Artist{
		ID:   uuid.New().String(),
		Name: name,
		Bio:  bio,
	}

	if err := as.ArtistRepo.Create(&artist); err != nil {
		return models.Artist{}, err
	}

	return artist, nil
}

func (as *Artistservice) DeleteArtist(ctx context.Context, id string) error {
	if err := as.ArtistRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (as *Artistservice) GetArtists(ctx context.Context, name string) ([]models.Artist, error) {
	if name != "" {
		return as.ArtistRepo.SearchByName(name)
	}
	return as.ArtistRepo.List()
}
func (as *Artistservice) GetArtistByID(ctx context.Context, id string) (*models.Artist, error) {
	return as.ArtistRepo.GetByID(id)
}
