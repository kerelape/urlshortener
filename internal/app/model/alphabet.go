package model

import "fmt"

// Alphabet is an alphabet of runes.
type Alphabet interface {
	fmt.Stringer

	// Size returns the amount of runes in this alphabet.
	Size() uint

	// Rune returns a rune by index.
	Rune(id uint) rune
}
