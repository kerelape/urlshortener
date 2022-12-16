package app

import "strconv"

// Encoding/decoding taken from https://gist.github.com/dgritsko/9554733
type DatabaseShortener struct {
	Database Database
	Cypher   Alphabet
}

func NewDatabaseShortener(database Database, cypher Alphabet) *DatabaseShortener {
	var shortener = new(DatabaseShortener)
	shortener.Database = database
	shortener.Cypher = cypher
	return shortener
}

func (self *DatabaseShortener) Shorten(origin string) string {
	var id = self.Database.Put(origin)
	if id == 0 {
		return string(self.Cypher.Rune(0))
	}
	var shortened = strconv.Itoa(int(id))
	println(shortened, origin)
	return shortened
}

func (self *DatabaseShortener) Reveal(shortened string) (string, error) {
	var id, err = strconv.Atoi(shortened)
	if err != nil {
		return "", err
	}
	return self.Database.Get(uint(id))
}
