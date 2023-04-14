package storage

import (
	"context"

	"github.com/kerelape/urlshortener/internal/app"
)

type History interface {
	Record(ctx context.Context, user app.Token, node HistoryNode) error
	GetRecordsByUser(ctx context.Context, user app.Token) ([]HistoryNode, error)
}

type HistoryNode struct {
	OriginalURL string
	ShortURL    string
}
