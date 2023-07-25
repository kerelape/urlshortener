package storage

import (
	"context"
	"errors"

	"github.com/kerelape/urlshortener/internal/app"
)

const deletedValue = "__DELETED"

// FakeDatabase is a database that stores values in RAM.
type FakeDatabase struct {
	Values []string
}

// Return new FakeDatabase.
func NewFakeDatabase() *FakeDatabase {
	return &FakeDatabase{}
}

// Put stores value and returns its id.
func (database *FakeDatabase) Put(ctx context.Context, _ app.Token, value string) (uint, error) {
	if database.Values == nil {
		return 0, ErrDatabaseClosed
	}

	for i, v := range database.Values {
		if v == value {
			return 0, NewDuplicateValueError(uint(i))
		}
	}
	database.Values = append(database.Values, value)
	return uint(len(database.Values) - 1), nil
}

// Get returns original value by its id.
func (database *FakeDatabase) Get(ctx context.Context, id uint) (string, error) {
	if database.Values == nil {
		return "", ErrDatabaseClosed
	}

	if id >= uint(len(database.Values)) {
		return "", errors.New("element does not exist")
	}
	value := database.Values[id]
	if value == deletedValue {
		return "", ErrValueDeleted
	}
	return value, nil
}

// PutAll stores values and returns their ids.
func (database *FakeDatabase) PutAll(ctx context.Context, user app.Token, values []string) ([]uint, error) {
	if database.Values == nil {
		return nil, ErrDatabaseClosed
	}

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

// GetAll returns original values by their ids.
func (database *FakeDatabase) GetAll(ctx context.Context, ids []uint) ([]string, error) {
	if database.Values == nil {
		return nil, ErrDatabaseClosed
	}

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

// Delete removes values by their ids.
func (database *FakeDatabase) Delete(ctx context.Context, _ app.Token, ids []uint) error {
	if database.Values == nil {
		return ErrDatabaseClosed
	}

	for _, i := range ids {
		database.Values[i] = deletedValue
	}
	return nil
}

// Ping always returns an error.
func (database *FakeDatabase) Ping(ctx context.Context) error {
	if database.Values == nil {
		return ErrDatabaseClosed
	}
	return errors.New("FakeDatabase")
}

func (database *FakeDatabase) Close(context.Context) error {
	database.Values = nil
	return nil
}
