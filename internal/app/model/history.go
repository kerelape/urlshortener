package model

import (
	"context"
	"errors"

	"github.com/kerelape/urlshortener/internal/app"
)

var ErrUserNotFound = errors.New("user not found")

type History interface {
	Record(ctx context.Context, user app.Token, node HistoryNode) error
	GetRecordsByUser(ctx context.Context, user app.Token) ([]HistoryNode, error)
}

type HistoryNode struct {
	From string
	To   string
}
