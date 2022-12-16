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

func (self *JoinedAlphabet) Size() uint {
	return self.Start.Size() + self.End.Size()
}

func (self *JoinedAlphabet) Rune(id uint) rune {
	var startSize = self.Start.Size()
	if id >= startSize {
		return self.End.Rune(id - startSize)
	}
	return self.Start.Rune(id)
}
