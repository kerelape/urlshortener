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

func (alphabet *AsciiAlphabet) Size() uint {
	return uint(alphabet.Max-alphabet.Min) + 1
}

func (alphabet *AsciiAlphabet) Rune(id uint) rune {
	if id >= alphabet.Size() {
		panic("Out of alphabet")
	}
	return rune(id + uint(alphabet.Min))
}
