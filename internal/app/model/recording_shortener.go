package model

import (
	"context"

	"github.com/kerelape/urlshortener/internal/app"
)

type ContextKey string

const ContextUserTokenKey = ContextKey("urlshortener.user.token")

type RecordingShortener struct {
	origin  Shortener
	history History
}

func NewRecordingShortener(origin Shortener, history History) RecordingShortener {
	return RecordingShortener{origin, history}
}

func (rs RecordingShortener) Shorten(ctx context.Context, origin string) (string, error) {
	short, err := rs.origin.Shorten(ctx, origin)
	if err != nil {
		return "", err
	}
	user, ok := ctx.Value(ContextUserTokenKey).(app.Token)
	if !ok {
		panic("context has no user")
	}
	err = rs.history.Record(ctx, user, HistoryNode{origin, short})
	if err != nil {
		return "", err
	}
	return short, nil
}

func (rs RecordingShortener) Reveal(ctx context.Context, short string) (string, error) {
	return rs.origin.Reveal(ctx, short)
}

func (rs RecordingShortener) ShortenAll(ctx context.Context, origins []string) ([]string, error) {
	shortened, err := rs.origin.ShortenAll(ctx, origins)
	if err != nil {
		return nil, err
	}
	user, ok := ctx.Value(ContextUserTokenKey).(app.Token)
	if !ok {
		panic("context has no user")
	}
	for i := range origins {
		err = rs.history.Record(ctx, user, HistoryNode{origins[i], shortened[i]})
		if err != nil {
			return nil, err
		}
	}
	return shortened, nil
}

func (rs RecordingShortener) RevealAll(ctx context.Context, shoretened []string) ([]string, error) {
	return rs.origin.RevealAll(ctx, shoretened)
}
