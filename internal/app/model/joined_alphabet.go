package model

// JoinedAlphabet is an alphabet that consists of two other alphabets.
type JoinedAlphabet struct {
	base Alphabet
	tail Alphabet
}

// NewJoinedAlphabet returns a new Alphabet.
func NewJoinedAlphabet(base Alphabet, tail Alphabet) *JoinedAlphabet {
	return &JoinedAlphabet{
		base: base,
		tail: tail,
	}
}

// Size returns the amount of runes in this alphabet.
func (alphabet *JoinedAlphabet) Size() uint {
	return alphabet.base.Size() + alphabet.tail.Size()
}

// Rune returns a rune by index.
func (alphabet *JoinedAlphabet) Rune(id uint) rune {
	startSize := alphabet.base.Size()
	if id >= startSize {
		return alphabet.tail.Rune(id - startSize)
	}
	return alphabet.base.Rune(id)
}

// String returns this alphabet as a sequence of runes in it.
func (alphabet *JoinedAlphabet) String() string {
	return alphabet.base.String() + alphabet.tail.String()
}
