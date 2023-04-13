package storage

import (
	"context"

	"github.com/kerelape/urlshortener/internal/app"
)

type VirtualHistory struct {
	history map[app.Token]([]HistoryNode)
}

func NewVirtualHistory() *VirtualHistory {
	return &VirtualHistory{
		history: map[app.Token]([]HistoryNode){},
	}
}

func (history *VirtualHistory) Record(_ context.Context, user app.Token, node HistoryNode) error {
	records := append(history.history[user], node)
	history.history[user] = records
	return nil
}

func (history *VirtualHistory) GetRecordsByUser(_ context.Context, user app.Token) ([]HistoryNode, error) {
	records := history.history[user]
	return records, nil
}
