package app

type JoinedAlphabet struct {
	Start Alphabet
	End   Alphabet
}

func NewJoinedAlphabet(start Alphabet, end Alphabet) *JoinedAlphabet {
	var alphabet = new(JoinedAlphabet)
	alphabet.Start = start
	alphabet.End = end
	return alphabet
}

func (alphabet *JoinedAlphabet) Size() uint {
	return alphabet.Start.Size() + alphabet.End.Size()
}

func (alphabet *JoinedAlphabet) Rune(id uint) rune {
	var startSize = alphabet.Start.Size()
	if id >= startSize {
		return alphabet.End.Rune(id - startSize)
	}
	return alphabet.Start.Rune(id)
}
