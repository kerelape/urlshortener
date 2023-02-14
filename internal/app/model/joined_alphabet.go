package model

type JoinedAlphabet struct {
	Base Alphabet
	Tail Alphabet
}

func NewJoinedAlphabet(base Alphabet, tail Alphabet) *JoinedAlphabet {
	return &JoinedAlphabet{
		Base: base,
		Tail: tail,
	}
}

func (alphabet *JoinedAlphabet) Size() uint {
	return alphabet.Base.Size() + alphabet.Tail.Size()
}

func (alphabet *JoinedAlphabet) Rune(id uint) rune {
	var startSize = alphabet.Base.Size()
	if id >= startSize {
		return alphabet.Tail.Rune(id - startSize)
	}
	return alphabet.Base.Rune(id)
}

func (alphabet *JoinedAlphabet) String() string {
	return alphabet.Base.String() + alphabet.Tail.String()
}
