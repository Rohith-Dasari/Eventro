package artistservice

import (
	"context"
	"eventro2/models"
)

//go:generate mockgen -destination=../../mocks/artist_service_mock.go -package=mocks -source=interface.go
type ArtistServiceI interface {
	CreateArtist(ctx context.Context, name, bio string) (models.Artist, error)
	DeleteArtist(ctx context.Context, id string) error
	GetArtists(ctx context.Context, name string) ([]models.Artist, error)
	GetArtistByID(ctx context.Context, id string) (*models.Artist, error)
}
