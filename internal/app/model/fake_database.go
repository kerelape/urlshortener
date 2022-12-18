package model

import "errors"

type FakeDatabase struct {
	Values []string
}

// Return new FakeDatabase.
func NewFakeDatabase() *FakeDatabase {
	return new(FakeDatabase)
}

func (database *FakeDatabase) Put(value string) uint {
	database.Values = append(database.Values, value)
	return uint(len(database.Values) - 1)
}

func (database *FakeDatabase) Get(id uint) (string, error) {
	if id >= uint(len(database.Values)) {
		return "", errors.New("element does not exist")
	}
	return database.Values[id], nil
}
