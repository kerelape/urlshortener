package model

func NewBase62Alphabet() Alphabet {
	return NewJoinedAlphabet(
		NewASCIIAlphabet(48, 57),
		NewJoinedAlphabet(
			NewASCIIAlphabet(65, 90),
			NewASCIIAlphabet(97, 122),
		),
	)
}
