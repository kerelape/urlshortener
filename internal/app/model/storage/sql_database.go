package storage

import "database/sql"

type SQLDatabase struct {
	db *sql.DB
}

func NewSQLDatabase(db *sql.DB) *SQLDatabase {
	return &SQLDatabase{
		db: db,
	}
}

func (database *SQLDatabase) Put(value string) (uint, error) {
	result, execError := database.db.Exec("INSERT INTO urls(origin) VALUES($1) RETURNING id", value)
	if execError != nil {
		return 0, execError
	}
	id, idError := result.LastInsertId()
	return uint(id), idError
}

func (database *SQLDatabase) Get(id uint) (string, error) {
	row := database.db.QueryRow("SELECT origin FROM urls WHERE id = $1", int64(id))
	var origin string
	scanError := row.Scan(&origin)
	return origin, scanError
}
