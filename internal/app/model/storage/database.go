package storage

type Database interface {
	Put(value string) (uint, error)
	Get(id uint) (string, error)
}
