package app

type Alphabet interface {
	Size() uint
	Rune(id uint) rune
}
