package storage

import (
	"context"

	"github.com/kerelape/urlshortener/internal/app"
)

// History is user's history.
type History interface {
	// Record records a user's action.
	Record(ctx context.Context, user app.Token, node HistoryNode) error

	// GetRecordsByUser returns all actions by the user.
	GetRecordsByUser(ctx context.Context, user app.Token) ([]HistoryNode, error)
}

// HistoryNode is a user's action.
type HistoryNode struct {
	// Original URL is the original URL.
	OriginalURL string

	// ShortURL is the shortened URL.
	ShortURL string
}
