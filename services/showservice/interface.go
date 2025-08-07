package showservice

import "context"

type ShowServiceInterface interface {
	ShowViewer
	ShowModerator
	ShowCreator
}

type ShowViewer interface {
	BrowseShowsByEvent(ctx context.Context, eventID string)
	DisplayShow(ctx context.Context, showID string)
	BrowseHostShows(ctx context.Context)
}

type ShowModerator interface {
	ModerateShow(ctx context.Context)
	ViewBlockedShows(ctx context.Context)
}

type ShowCreator interface {
	CreateShow(ctx context.Context)
	RemoveHostShow(ctx context.Context)
}
