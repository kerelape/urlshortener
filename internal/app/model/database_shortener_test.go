package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabaseShortener(t *testing.T) {
	var shortener = NewDatabaseShortener(NewFakeDatabase())
	var short = shortener.Shorten("Hello, World!")
	var real, err = shortener.Reveal(short)
	require.Nil(t, err)
	require.Equal(t, "Hello, World!", real)
}
