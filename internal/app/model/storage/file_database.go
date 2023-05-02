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

var ErrTooLargeValue = errors.New("too large value")

type FileDatabase struct {
	file      *os.File
	rw        sync.Mutex
	buffer    []byte
	chunkSize int
}

func NewFileDatabase(file *os.File, chunkSize int) *FileDatabase {
	return &FileDatabase{
		file:      file,
		buffer:    make([]byte, chunkSize),
		chunkSize: chunkSize,
	}
}

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

func (database *FileDatabase) Put(ctx context.Context, user app.Token, value string) (uint, error) {
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
			if err != nil {
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

func (database *FileDatabase) Get(ctx context.Context, id uint) (string, error) {
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

func (database *FileDatabase) PutAll(ctx context.Context, user app.Token, values []string) ([]uint, error) {
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

func (database *FileDatabase) GetAll(ctx context.Context, ids []uint) ([]string, error) {
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

func (database *FileDatabase) Delete(ctx context.Context, _ app.Token, ids []uint) error {
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

func (database *FileDatabase) Ping(ctx context.Context) error {
	return errors.New("FileDatabase")
}
