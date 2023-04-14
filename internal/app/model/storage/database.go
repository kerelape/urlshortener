package storage

import (
	"context"
	"fmt"
)

type Database interface {
	Put(ctx context.Context, value string) (uint, error)
	Get(ctx context.Context, id uint) (string, error)
	PutAll(ctx context.Context, values []string) ([]uint, error)
	GetAll(ctx context.Context, ids []uint) ([]string, error)
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
