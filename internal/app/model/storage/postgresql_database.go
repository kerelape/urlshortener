package storage

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQLDatabase struct {
	db *sql.DB
}

func NewPostgreSQLDatabase(db *sql.DB) *PostgreSQLDatabase {
	return &PostgreSQLDatabase{
		db: db,
	}
}

func DialPostgreSQLDatabase(dsn string) (*PostgreSQLDatabase, error) {
	db, openError := sql.Open("pgx", dsn)
	if openError != nil {
		return nil, openError
	}
	_, execError := db.Exec("CREATE TABLE IF NOT EXISTS urls(id int primary key auto_incerement, origin text)")
	return NewPostgreSQLDatabase(db), execError
}

func (database *PostgreSQLDatabase) Put(value string) (uint, error) {
	result, execError := database.db.Exec("INSERT INTO urls(origin) VALUES($1) RETURNING id", value)
	if execError != nil {
		return 0, execError
	}
	id, idError := result.LastInsertId()
	return uint(id), idError
}

func (database *PostgreSQLDatabase) Get(id uint) (string, error) {
	row := database.db.QueryRow("SELECT origin FROM urls WHERE id = $1", int64(id))
	var origin string
	scanError := row.Scan(&origin)
	return origin, scanError
}

func (database *PostgreSQLDatabase) Ping() error {
	return database.Ping()
}
