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

func (shortener *DatabaseShortener) Shorten(origin string) string {
	var id = shortener.Database.Put(origin)
	var shortened = strconv.Itoa(int(id))
	return shortened
}

func (shortener *DatabaseShortener) Reveal(shortened string) (string, error) {
	var id, err = strconv.Atoi(shortened)
	if err != nil {
		return "", err
	}
	return shortener.Database.Get(uint(id))
}
