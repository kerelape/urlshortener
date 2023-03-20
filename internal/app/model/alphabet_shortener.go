package model

import (
	"strings"

	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

type AlphabetShortener struct {
	Database storage.Database
	Alphabet Alphabet
}

func NewAlphabetShortener(database storage.Database, alphabet Alphabet) *AlphabetShortener {
	return &AlphabetShortener{
		Database: database,
		Alphabet: alphabet,
	}
}

func (shortener *AlphabetShortener) Shorten(origin string) (string, error) {
	number, putError := shortener.Database.Put(origin)
	if putError != nil {
		return "", putError
	}
	if number == 0 {
		return string(shortener.Alphabet.Rune(0)), nil
	}
	var cypher []rune
	base := shortener.Alphabet.Size()
	for i := number; i > 0; i /= base {
		cypher = append([]rune{shortener.Alphabet.Rune(i % base)}, cypher...)
	}
	return string(cypher), nil
}

func (shortener *AlphabetShortener) Reveal(shortened string) (string, error) {
	var decoded uint
	encoded := []rune(shortened)
	lookup := shortener.Alphabet.String()
	base := shortener.Alphabet.Size()
	for i := 0; i < len(encoded); i++ {
		decoded = (decoded * base) + uint(strings.IndexRune(lookup, encoded[i]))
	}
	return shortener.Database.Get(decoded)
}

func (shortener *AlphabetShortener) ShortenAll(origins []string) ([]string, error) {
	result := make([]string, len(origins))
	for i, origin := range origins {
		short, shortenError := shortener.Shorten(origin)
		if shortenError != nil {
			return nil, shortenError
		}
		result[i] = short
	}
	return result, nil
}

func (shortener *AlphabetShortener) RevealAll(shortened []string) ([]string, error) {
	result := make([]string, len(shortened))
	for i, shortened := range shortened {
		reveal, revealError := shortener.Reveal(shortened)
		if revealError != nil {
			return nil, revealError
		}
		result[i] = reveal
	}
	return result, nil
}
