package model

// ASCIIALphabet is an alphabet that contains a sequence of ascii characters
// in the specified range.
type ASCIIAlphabet struct {
	min uint8
	max uint8
}

// NewASCIIAlphabet returns a new ASCIIAlphabet.
func NewASCIIAlphabet(min uint8, max uint8) *ASCIIAlphabet {
	return &ASCIIAlphabet{
		min: min,
		max: max,
	}
}

// Size returns the amount of runes in this alphabet.
func (alphabet *ASCIIAlphabet) Size() uint {
	return uint(alphabet.max-alphabet.min) + 1
}

// Rune returns a rune by index.
func (alphabet *ASCIIAlphabet) Rune(id uint) rune {
	if id >= alphabet.Size() {
		panic("Out of alphabet")
	}
	return rune(id + uint(alphabet.min))
}

// String returns this alphabet as a sequence of runes in it.
func (alphabet *ASCIIAlphabet) String() string {
	var runes []rune
	for i := uint(0); i < alphabet.Size(); i++ {
		runes = append(runes, alphabet.Rune(i))
	}
	return string(runes)
}
