package app

import (
	"encoding/hex"
	"math/rand"
)

type Token [8]byte

func NewToken() Token {
	var token Token
	rand.Read(token[:])
	return token
}

func TokenFromString(origin string) (Token, error) {
	bytes, decodeError := hex.DecodeString(origin)
	var token Token
	copy(token[:], bytes)
	return token, decodeError
}
