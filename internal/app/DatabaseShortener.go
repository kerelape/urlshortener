package app

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
	var shortened string
	var base = self.Cypher.Size()
	for id > 0 {
		shortened = shortened + string(self.Cypher.Rune(id%base))
		id = id / base
	}
	return shortened
}

func (self *DatabaseShortener) Reveal(shortened string) (string, error) {
	var id uint
	var base = self.Cypher.Size()
	for i := 0; i < len(shortened); i++ {
		id = id*base + uint(self.Cypher.Rune(uint(shortened[i])))
	}
	return self.Database.Get(id)
}
