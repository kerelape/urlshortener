package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/kerelape/urlshortener/internal/app"
)

var ErrValueDeleted = errors.New("deleted")

type Database interface {
	Put(ctx context.Context, user app.Token, value string) (uint, error)
	Get(ctx context.Context, id uint) (string, error)
	PutAll(ctx context.Context, user app.Token, values []string) ([]uint, error)
	GetAll(ctx context.Context, ids []uint) ([]string, error)
	Delete(ctx context.Context, user app.Token, ids []uint) error
	Ping(ctx context.Context) error
}

type DuplicateValueError struct {
	Origin uint
}

func NewDuplicateValueError(origin uint) DuplicateValueError {
	return DuplicateValueError{
		Origin: origin,
	}
}

func (e DuplicateValueError) Error() string {
	return fmt.Sprintf("duplicate value of ID: %d", e.Origin)
}
