package storage

type Database interface {
	Put(value string) (uint, error)
	Get(id uint) (string, error)
	PutAll(values []string) ([]uint, error)
	GetAll(ids []uint) ([]string, error)
	Ping() error
}

type ErrDuplicate struct {
	Origin uint
}

func (err *ErrDuplicate) Error() string {
	return "Duplicate"
}
