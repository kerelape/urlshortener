package model

import (
	"context"

	"github.com/kerelape/urlshortener/internal/app"
)

type FakeHistory struct {
	records map[app.Token][]HistoryNode
}

func NewFakeHistory() FakeHistory {
	return FakeHistory{make(map[app.Token][]HistoryNode)}
}

func (fh FakeHistory) Record(ctx context.Context, user app.Token, node HistoryNode) error {
	records, ok := fh.records[user]
	if !ok {
		fh.records[user] = make([]HistoryNode, 0)
		return fh.Record(ctx, user, node)
	}
	fh.records[user] = append(records, node)
	return nil
}

func (fh FakeHistory) GetRecordsByUser(ctx context.Context, user app.Token) ([]HistoryNode, error) {
	records, ok := fh.records[user]
	if !ok {
		return []HistoryNode{}, nil
	}
	return records, nil
}
