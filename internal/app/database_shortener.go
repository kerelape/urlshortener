package app

import "strconv"

// Encoding/decoding taken from https://gist.github.com/dgritsko/9554733
type DatabaseShortener struct {
	Database Database
}

func NewDatabaseShortener(database Database) *DatabaseShortener {
	var shortener = new(DatabaseShortener)
	shortener.Database = database
	return shortener
}

func (self *DatabaseShortener) Shorten(origin string) string {
	var id = self.Database.Put(origin)
	var shortened = strconv.Itoa(int(id))
	return shortened
}

func (self *DatabaseShortener) Reveal(shortened string) (string, error) {
	var id, err = strconv.Atoi(shortened)
	if err != nil {
		return "", err
	}
	return self.Database.Get(uint(id))
}
