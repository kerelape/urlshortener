package app

type AsciiAlphabet struct {
	Min uint8
	Max uint8
}

func NewAsciiAlphabet(min uint8, max uint8) *AsciiAlphabet {
	var alphabet = new(AsciiAlphabet)
	alphabet.Min = min
	alphabet.Max = max
	return alphabet
}

func (self *AsciiAlphabet) Size() uint {
	return uint(self.Max-self.Min) + 1
}

func (self *AsciiAlphabet) Rune(id uint) rune {
	if id >= self.Size() {
		panic("Out of alphabet")
	}
	return rune(id + uint(self.Min))
}
