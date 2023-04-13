package model

import (
	"context"
	"errors"
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

func (shortener *AlphabetShortener) encode(number uint) string {
	if number == 0 {
		return string(shortener.Alphabet.Rune(0))
	}
	var cypher []rune
	base := shortener.Alphabet.Size()
	for i := number; i > 0; i /= base {
		cypher = append([]rune{shortener.Alphabet.Rune(i % base)}, cypher...)
	}
	return string(cypher)
}

func (shortener *AlphabetShortener) decode(encoded string) uint {
	result := uint(0)
	cypher := []rune(encoded)
	lookup := shortener.Alphabet.String()
	base := shortener.Alphabet.Size()
	for i := 0; i < len(cypher); i++ {
		result = (result * base) + uint(strings.IndexRune(lookup, cypher[i]))
	}
	return result
}

func (shortener *AlphabetShortener) Shorten(ctx context.Context, origin string) (string, error) {
	number, putError := shortener.Database.Put(ctx, origin)
	if putError != nil {
		var duplicate storage.DuplicateValueError
		if errors.As(putError, &duplicate) {
			return "", NewDuplicateURLError(shortener.encode(duplicate.Origin))
		}
		return "", putError
	}
	return shortener.encode(number), nil
}

func (shortener *AlphabetShortener) Reveal(ctx context.Context, shortened string) (string, error) {
	return shortener.Database.Get(ctx, shortener.decode(shortened))
}

func (shortener *AlphabetShortener) ShortenAll(ctx context.Context, origins []string) ([]string, error) {
	ids, putError := shortener.Database.PutAll(ctx, origins)
	if putError != nil {
		return nil, putError
	}
	result := make([]string, 0, len(origins))
	for _, id := range ids {
		result = append(result, shortener.encode(id))
	}
	return result, nil
}

func (shortener *AlphabetShortener) RevealAll(ctx context.Context, shortened []string) ([]string, error) {
	ids := make([]uint, len(shortened))
	for _, id := range shortened {
		ids = append(ids, shortener.decode(id))
	}
	values, getError := shortener.Database.GetAll(ctx, ids)
	if getError != nil {
		return nil, getError
	}
	return values, nil
}
