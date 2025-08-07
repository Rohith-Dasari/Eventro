package searchevents

import "context"

type SearchServiceInterface interface {
	Search(ctx context.Context)
	SearchByEventName()
	SearchByCategory()
	SearchByLocation()
}
