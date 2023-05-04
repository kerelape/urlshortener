package storage

import (
	"context"
	"errors"

	"github.com/kerelape/urlshortener/internal/app"
)

const deletedValue = "__DELETED"

type FakeDatabase struct {
	Values []string
}

// Return new FakeDatabase.
func NewFakeDatabase() *FakeDatabase {
	return &FakeDatabase{}
}

func (database *FakeDatabase) Put(ctx context.Context, _ app.Token, value string) (uint, error) {
	for i, v := range database.Values {
		if v == value {
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
	value := database.Values[id]
	if value == deletedValue {
		return "", ErrValueDeleted
	}
	return value, nil
}

func (database *FakeDatabase) PutAll(ctx context.Context, user app.Token, values []string) ([]uint, error) {
	result := make([]uint, len(values))
	for i := 0; i < len(values); i++ {
		id, putError := database.Put(ctx, user, values[i])
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

func (database *FakeDatabase) Delete(ctx context.Context, _ app.Token, ids []uint) error {
	for _, i := range ids {
		database.Values[i] = deletedValue
	}
	return nil
}

func (database *FakeDatabase) Ping(ctx context.Context) error {
	return errors.New("FakeDatabase")
}
