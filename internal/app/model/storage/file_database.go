package storage

import (
	"bytes"
	"errors"
	"os"
	"sync"
)

const FileDatabaseChunkSize = 2048

var ErrTooLargeValue = errors.New("too large value")

type FileDatabase struct {
	file   *os.File
	rw     sync.Mutex
	buffer []byte
}

func NewFileDatabase(file *os.File) *FileDatabase {
	return &FileDatabase{file: file, buffer: make([]byte, FileDatabaseChunkSize)}
}

func (database *FileDatabase) Put(value string) (uint, error) {
	database.rw.Lock()
	defer database.rw.Unlock()
	var stat, statError = database.file.Stat()
	if statError != nil {
		return 0, statError
	}
	if len(value) > FileDatabaseChunkSize {
		return 0, ErrTooLargeValue
	}
	var id = stat.Size() / FileDatabaseChunkSize
	var buffer = append([]byte(value), make([]byte, FileDatabaseChunkSize-len(value))...)
	var _, writeError = database.file.WriteAt(buffer, int64(id*FileDatabaseChunkSize))
	return uint(id), writeError
}

func (database *FileDatabase) Get(id uint) (string, error) {
	database.rw.Lock()
	defer database.rw.Unlock()
	var buffer = make([]byte, FileDatabaseChunkSize)
	var _, readError = database.file.ReadAt(buffer, int64(id)*FileDatabaseChunkSize)
	if readError != nil {
		return "", readError
	}
	var value, readStringError = bytes.NewBuffer(buffer).ReadString(0x00)
	return value[:len(value)-1], readStringError
}
