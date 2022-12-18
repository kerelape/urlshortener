package model

import "fmt"

type Alphabet interface {
	fmt.Stringer
	Size() uint
	Rune(id uint) rune
}
