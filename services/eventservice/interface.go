package eventservice

import (
	"context"
	"eventro2/models"
)

type EventServiceI interface {
	ModerateEvents(ctx context.Context)
	ViewBlockedEvents(ctx context.Context)
	CreateNewEvent() models.Event
}
