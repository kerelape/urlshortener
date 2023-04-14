package storage

import (
	"context"
	"errors"
)

type FakeDatabase struct {
	Values []string
}

// Return new FakeDatabase.
func NewFakeDatabase() *FakeDatabase {
	return &FakeDatabase{}
}

func (database *FakeDatabase) Put(ctx context.Context, value string) (uint, error) {
	for i, u := range database.Values {
		if u == value {
			return 0, NewDuplicateValueError(uint(i))
		}
	}
	database.Values = append(database.Values, value)
	return uint(len(database.Values) - 1), nil
}

func (database *FakeDatabase) Get(ctx context.Context, id uint) (string, error) {
	if id >= uint(len(database.Values)) {
		return "", errors.New("element does not exist")
	}
	return database.Values[id], nil
}

func (database *FakeDatabase) PutAll(ctx context.Context, values []string) ([]uint, error) {
	result := make([]uint, len(values))
	for i := 0; i < len(values); i++ {
		id, putError := database.Put(ctx, values[i])
		if putError != nil {
			return nil, putError
		}
		result[i] = id
	}
	return result, nil
}

func (database *FakeDatabase) GetAll(ctx context.Context, ids []uint) ([]string, error) {
	result := make([]string, len(ids))
	for i := 0; i < len(ids); i++ {
		value, getError := database.Get(ctx, ids[i])
		if getError != nil {
			return nil, getError
		}
		result[i] = value
	}
	return result, nil
}

func (database *FakeDatabase) Ping(ctx context.Context) error {
	return errors.New("FakeDatabase")
}
