package model

import (
	"context"
	"errors"
	"strings"

	"github.com/kerelape/urlshortener/internal/app"
	"github.com/kerelape/urlshortener/internal/app/model/storage"
)

// AlphabetShortener is a shortener that uses an Alphabet to encode database ids.
type AlphabetShortener struct {
	database storage.Database
	alphabet Alphabet
}

// NewAlphabetShortener returns a new AlphabetShortener.
func NewAlphabetShortener(database storage.Database, alphabet Alphabet) *AlphabetShortener {
	return &AlphabetShortener{
		database: database,
		alphabet: alphabet,
	}
}

func (shortener *AlphabetShortener) encode(number uint) string {
	if number == 0 {
		return string(shortener.alphabet.Rune(0))
	}
	var cypher []rune
	base := shortener.alphabet.Size()
	for i := number; i > 0; i /= base {
		cypher = append([]rune{shortener.alphabet.Rune(i % base)}, cypher...)
	}
	return string(cypher)
}

func (shortener *AlphabetShortener) decode(encoded string) uint {
	result := uint(0)
	cypher := []rune(encoded)
	lookup := shortener.alphabet.String()
	base := shortener.alphabet.Size()
	for i := 0; i < len(cypher); i++ {
		result = (result * base) + uint(strings.IndexRune(lookup, cypher[i]))
	}
	return result
}

// Shorten shortens the given origin string.
func (shortener *AlphabetShortener) Shorten(ctx context.Context, user app.Token, origin string) (string, error) {
	number, putError := shortener.database.Put(ctx, user, origin)
	if putError != nil {
		var duplicate storage.DuplicateValueError
		if errors.As(putError, &duplicate) {
			return "", NewDuplicateURLError(shortener.encode(duplicate.Origin))
		}
		return "", putError
	}
	return shortener.encode(number), nil
}

// Reveal returns the original string by the shortened.
func (shortener *AlphabetShortener) Reveal(ctx context.Context, shortened string) (string, error) {
	return shortener.database.Get(ctx, shortener.decode(shortened))
}

// ShortenAll shortens a slice of strings and returns
// a slice of short string in the same order.
func (shortener *AlphabetShortener) ShortenAll(ctx context.Context, user app.Token, origins []string) ([]string, error) {
	ids, putError := shortener.database.PutAll(ctx, user, origins)
	if putError != nil {
		return nil, putError
	}
	result := make([]string, 0, len(origins))
	for _, id := range ids {
		result = append(result, shortener.encode(id))
	}
	return result, nil
}

// RevealAll returns a slice of original strings in the same order
// as in shortened.
func (shortener *AlphabetShortener) RevealAll(ctx context.Context, shortened []string) ([]string, error) {
	ids := make([]uint, len(shortened))
	for _, id := range shortened {
		ids = append(ids, shortener.decode(id))
	}
	values, getError := shortener.database.GetAll(ctx, ids)
	if getError != nil {
		return nil, getError
	}
	return values, nil
}

// Delete deletes a string from the shortener.
func (shortener *AlphabetShortener) Delete(ctx context.Context, user app.Token, shortened []string) error {
	ids := make([]uint, len(shortened))
	for _, id := range shortened {
		ids = append(ids, shortener.decode(id))
	}
	return shortener.database.Delete(ctx, user, ids)
}
