package model

import "strings"

type AlphabetShortener struct {
	Database Database
	Alphabet Alphabet
}

func NewAlphabetShortener(database Database, alphabet Alphabet) *AlphabetShortener {
	var shortener = new(AlphabetShortener)
	shortener.Database = database
	shortener.Alphabet = alphabet
	return shortener
}

func (shortener *AlphabetShortener) Shorten(origin string) string {
	var number = shortener.Database.Put(origin)
	if number == 0 {
		return string(shortener.Alphabet.Rune(0))
	}
	var cypher []rune
	var base = shortener.Alphabet.Size()
	for i := number; i > 0; i /= base {
		cypher = append([]rune{shortener.Alphabet.Rune(i % base)}, cypher...)
	}
	return string(cypher)
}

func (shortener *AlphabetShortener) Reveal(shortened string) (string, error) {
	var decoded uint
	var encoded = []rune(shortened)
	var lookup = shortener.Alphabet.String()
	var base = shortener.Alphabet.Size()
	for i := 0; i < len(encoded); i++ {
		decoded = (decoded * base) + uint(strings.IndexRune(lookup, encoded[i]))
	}
	return shortener.Database.Get(decoded)
}
