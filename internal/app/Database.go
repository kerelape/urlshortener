package app

type Database interface {
	Put(value string) uint
	Get(id uint) (string, error)
}
