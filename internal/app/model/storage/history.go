package storage

import "github.com/kerelape/urlshortener/internal/app"

type History interface {
	Record(user app.Token, node *HistoryNode) error
	GetRecordsByUser(user app.Token) ([]*HistoryNode, error)
}

type HistoryNode struct {
	OriginalURL string
	ShortURL    string
}
