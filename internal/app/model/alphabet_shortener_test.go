package model

import (
	"context"
	"testing"

	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
	"github.com/stretchr/testify/require"
)

func TestAlphabetShortener(t *testing.T) {
	shortener := NewAlphabetShortener(storage.NewFakeDatabase(), NewASCIIAlphabet(97, 122))
	short, shortenError := shortener.Shorten(context.Background(), app.NewToken(), "Hello, World!")
	real, err := shortener.Reveal(context.Background(), short)
	require.Nil(t, err)
	require.Nil(t, shortenError)
	require.Equal(t, "Hello, World!", real)
}
