package storage

import (
	"context"

	"github.com/kerelape/urlshortener/internal/app"
)

// VirtualHistory is a History that stores records in RAM.
type VirtualHistory struct {
	history map[app.Token]([]HistoryNode)
}

// NewVirtualHistory returns a new VirtualHistory.
func NewVirtualHistory() *VirtualHistory {
	return &VirtualHistory{
		history: map[app.Token]([]HistoryNode){},
	}
}

// Record records a user's action.
func (history *VirtualHistory) Record(_ context.Context, user app.Token, node HistoryNode) error {
	records := append(history.history[user], node)
	history.history[user] = records
	return nil
}

// GetRecordsByUser returns all actions by the user.
func (history *VirtualHistory) GetRecordsByUser(_ context.Context, user app.Token) ([]HistoryNode, error) {
	records := history.history[user]
	return records, nil
}
