package storage

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/kerelape/urlshortener/internal/app"
)

// ErrTooLargeValue is returned when the provided stirng is to long.
var ErrTooLargeValue = errors.New("too large value")

// FileDatabase is a database that stores values in a file.
type FileDatabase struct {
	file      *os.File
	rw        sync.Mutex
	buffer    []byte
	chunkSize int
}

// NewFileDatabase returns a new FileDatabase.
func NewFileDatabase(file *os.File, chunkSize int) *FileDatabase {
	return &FileDatabase{
		file:      file,
		buffer:    make([]byte, chunkSize),
		chunkSize: chunkSize,
	}
}

// OpenFileDatabase opens a file and returned a new FileDatabase associated with the file.
func OpenFileDatabase(name string, create bool, permission fs.FileMode, chunkSize int) (*FileDatabase, error) {
	flag := os.O_RDWR
	if create {
		flag |= os.O_CREATE
	}
	file, openError := os.OpenFile(name, flag, permission)
	if openError != nil {
		return nil, openError
	}
	return NewFileDatabase(file, chunkSize), nil
}

// Put stores value and returns its id.
func (database *FileDatabase) Put(ctx context.Context, user app.Token, value string) (uint, error) {
	if database.file == nil {
		return 0, ErrDatabaseClosed
	}

	database.rw.Lock()
	stat, statError := database.file.Stat()
	database.rw.Unlock()
	if statError != nil {
		return 0, statError
	}
	if len(value) > database.chunkSize {
		return 0, ErrTooLargeValue
	}
	id := stat.Size() / int64(database.chunkSize)
	if id != 0 {
		for i := uint(0); i < uint(id); i++ {
			sameURL, err := database.Get(ctx, i)
			if err != nil && !errors.Is(err, ErrValueDeleted) {
				return 0, err
			}
			if sameURL == value {
				return 0, NewDuplicateValueError(i)
			}
		}
	}
	buffer := append([]byte(value), make([]byte, database.chunkSize-len(value))...)
	database.rw.Lock()
	_, writeError := database.file.WriteAt(buffer, id*int64(database.chunkSize))
	database.rw.Unlock()
	return uint(id), writeError
}

// Get returns original value by its id.
func (database *FileDatabase) Get(ctx context.Context, id uint) (string, error) {
	if database.file == nil {
		return "", ErrDatabaseClosed
	}

	database.rw.Lock()
	defer database.rw.Unlock()
	buffer := make([]byte, database.chunkSize)
	_, readError := database.file.ReadAt(buffer, int64(int(id)*database.chunkSize))
	if readError != nil {
		return "", readError
	}
	value, readStringError := bytes.NewBuffer(buffer).ReadString(0x00)
	if strings.HasPrefix(value, deletedValue) {
		return "", ErrValueDeleted
	}
	return value[:len(value)-1], readStringError
}

// PutAll stores values and returns their ids.
func (database *FileDatabase) PutAll(ctx context.Context, user app.Token, values []string) ([]uint, error) {
	if database.file == nil {
		return nil, ErrDatabaseClosed
	}

	result := make([]uint, len(values))
	for i := 0; i < len(values); i++ {
		id, putError := database.Put(ctx, user, values[i])
		if putError != nil {
			var duplicate DuplicateValueError
			if errors.As(putError, &duplicate) {
				id = duplicate.Origin
			} else {
				return nil, putError
			}
		}
		result[i] = id
	}
	return result, nil
}

// GetAll returns original values by their ids.
func (database *FileDatabase) GetAll(ctx context.Context, ids []uint) ([]string, error) {
	if database.file == nil {
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
func (database *FileDatabase) Delete(ctx context.Context, _ app.Token, ids []uint) error {
	if database.file == nil {
		return ErrDatabaseClosed
	}

	database.rw.Lock()
	defer database.rw.Unlock()
	for _, i := range ids {
		_, err := database.file.WriteAt([]byte(deletedValue), int64(i)*int64(database.chunkSize))
		if err != nil {
			return err
		}
	}
	return nil
}

// URLs returns count of URL stored in this file database.
func (database *FileDatabase) URLs(ctx context.Context) (int, error) {
	if database.file == nil {
		return -1, ErrDatabaseClosed
	}
	stat, err := database.file.Stat()
	if err != nil {
		return -1, err
	}
	return int(stat.Size() / int64(database.chunkSize)), nil
}

// Users always return -1 and an error indicating that the database does not
// support users.
func (database *FileDatabase) Users(ctx context.Context) (int, error) {
	if database.file == nil {
		return -1, ErrDatabaseClosed
	}
	return -1, errors.New("FileDatabase doesn't support users")
}

// Ping always returns an error.
func (database *FileDatabase) Ping(ctx context.Context) error {
	return errors.New("FileDatabase")
}

// Close closes the file.
func (database *FileDatabase) Close(ctx context.Context) error {
	file := database.file
	database.file = nil
	return file.Close()
}
