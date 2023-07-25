package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/kerelape/urlshortener/internal/app"
)

var (
	// ErrValueDeleted is returned when the values has been removes from the database.
	ErrValueDeleted = errors.New("deleted")

	// ErrDatabaseClosed is returned when the database has been closed.
	ErrDatabaseClosed = errors.New("database is closed")
)

// Database is a storage of strings.
type Database interface {
	// Put stores the given string value and returns its id.
	Put(ctx context.Context, user app.Token, value string) (uint, error)

	// Get returns the stored string by id.
	Get(ctx context.Context, id uint) (string, error)

	// PutAll stores many strings and returns their ids in the same order.
	PutAll(ctx context.Context, user app.Token, values []string) ([]uint, error)

	// GetAll returns the stored strings by the ids.
	GetAll(ctx context.Context, ids []uint) ([]string, error)

	// Delete removes a string by its id from the database.
	Delete(ctx context.Context, user app.Token, ids []uint) error

	// Ping returns an error if the database is unavailable.
	Ping(ctx context.Context) error

	// Close closes the databse.
	Close(ctx context.Context) error
}

// DuplicateValueError is returned when trying to put a string
// into the database that has already been added.
type DuplicateValueError struct {
	Origin uint
}

// NewDuplicateValueError returns a new DuplicateValueError.
func NewDuplicateValueError(origin uint) DuplicateValueError {
	return DuplicateValueError{
		Origin: origin,
	}
}

// Error returns description of the error.
func (e DuplicateValueError) Error() string {
	return fmt.Sprintf("duplicate value of ID: %d", e.Origin)
}
