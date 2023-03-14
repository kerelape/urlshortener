package storage

import "errors"

type FakeDatabase struct {
	Values []string
}

// Return new FakeDatabase.
func NewFakeDatabase() *FakeDatabase {
	return &FakeDatabase{}
}

func (database *FakeDatabase) Put(value string) (uint, error) {
	database.Values = append(database.Values, value)
	return uint(len(database.Values) - 1), nil
}

func (database *FakeDatabase) Get(id uint) (string, error) {
	if id >= uint(len(database.Values)) {
		return "", errors.New("element does not exist")
	}
	return database.Values[id], nil
}

func (database *FakeDatabase) PutAll(values []string) ([]uint, error) {
	result := make([]uint, len(values))
	for i := 0; i < len(values); i++ {
		id, putError := database.Put(values[i])
		if putError != nil {
			return nil, putError
		}
		result[i] = id
	}
	return result, nil
}

func (database *FakeDatabase) GetAll(ids []uint) ([]string, error) {
	result := make([]string, len(ids))
	for i := 0; i < len(ids); i++ {
		value, getError := database.Get(ids[i])
		if getError != nil {
			return nil, getError
		}
		result[i] = value
	}
	return result, nil
}

func (database *FakeDatabase) Ping() error {
	return errors.New("FakeDatabase")
}
