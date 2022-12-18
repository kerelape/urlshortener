package model

type ASCIIAlphabet struct {
	Min uint8
	Max uint8
}

func NewASCIIAlphabet(min uint8, max uint8) *ASCIIAlphabet {
	var alphabet = new(ASCIIAlphabet)
	alphabet.Min = min
	alphabet.Max = max
	return alphabet
}

func (alphabet *ASCIIAlphabet) Size() uint {
	return uint(alphabet.Max-alphabet.Min) + 1
}

func (alphabet *ASCIIAlphabet) Rune(id uint) rune {
	if id >= alphabet.Size() {
		panic("Out of alphabet")
	}
	return rune(id + uint(alphabet.Min))
}

func (alphabet *ASCIIAlphabet) String() string {
	var runes []rune
	for i := uint(0); i < alphabet.Size(); i++ {
		runes = append(runes, alphabet.Rune(i))
	}
	return string(runes)
}
