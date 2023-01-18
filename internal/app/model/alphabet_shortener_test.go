package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlphabetShortener(t *testing.T) {
	var shortener = NewAlphabetShortener(NewFakeDatabase(), NewASCIIAlphabet(97, 122))
	var short, shortenError = shortener.Shorten("Hello, World!")
	var real, err = shortener.Reveal(short)
	require.Nil(t, err)
	require.Nil(t, shortenError)
	require.Equal(t, "Hello, World!", real)
}
