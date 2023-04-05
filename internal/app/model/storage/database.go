package storage

import "fmt"

type Database interface {
	Put(value string) (uint, error)
	Get(id uint) (string, error)
	PutAll(values []string) ([]uint, error)
	GetAll(ids []uint) ([]string, error)
	Ping() error
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
