package stats

import "context"

// StatsProvider represents a type that can provide stats needed by the Stats API.
type StatsProvider interface {
	// URLs returns count of shorten urls.
	URLs(ctx context.Context) (int, error)

	// Users returns count of users registered in the service.
	Users(ctx context.Context) (int, error)
}
