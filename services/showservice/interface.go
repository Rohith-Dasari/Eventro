package showservice

import (
	"context"
	"eventro2/models"
	"time"
)

//go:generate mockgen -destination=../../mocks/show_service_mock.go -package=mocks -source=interface.go
type ShowServiceInterface interface {
	UpdateShow(ctx context.Context, showID string, userID string, update models.UpdateShowData) (models.ShowResponse, error)
	BrowseShows(ctx context.Context, filter models.ShowFilter) ([]models.ShowResponse, error)
	CreateShow(ctx context.Context, eventID string, venueID string, hostID string, price float64, showDate time.Time, showTime string) (models.ShowResponse, error)
	DeleteShow(ctx context.Context, showID string) error
}
