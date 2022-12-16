package app

import "errors"

type FakeDatabase struct {
	Values []string
}

// Return new FakeDatabase.
func NewFakeDatabase() *FakeDatabase {
	return new(FakeDatabase)
}

func (self *FakeDatabase) Put(value string) uint {
	self.Values = append(self.Values, value)
	return uint(len(self.Values) - 1)
}

func (self *FakeDatabase) Get(id uint) (string, error) {
	if id >= uint(len(self.Values)) {
		return "", errors.New("element does not exist")
	}
	return self.Values[id], nil
}
