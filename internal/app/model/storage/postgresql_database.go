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
	_, execError := db.Exec("CREATE TABLE IF NOT EXISTS urls(id SERIAL NOT NULL PRIMARY KEY, origin TEXT)")
	return NewPostgreSQLDatabase(db), execError
}

func (database *PostgreSQLDatabase) Put(value string) (uint, error) {
	same := database.db.QueryRow("SELECT id FROM urls WHERE origin = $1", value)
	if same.Err() != nil {
		return 0, same.Err()
	}
	var sameID int64
	if err := same.Scan(&sameID); err == nil {
		return 0, NewDuplicateValueError(uint(sameID))
	}
	row := database.db.QueryRow("INSERT INTO urls(origin) VALUES($1) RETURNING id", value)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var id int64
	idError := row.Scan(&id)
	return uint(id), idError
}

func (database *PostgreSQLDatabase) Get(id uint) (string, error) {
	row := database.db.QueryRow("SELECT origin FROM urls WHERE id = $1", int64(id))
	var origin string
	scanError := row.Scan(&origin)
	return origin, scanError
}

func (database *PostgreSQLDatabase) PutAll(values []string) ([]uint, error) {
	transaction, beginError := database.db.Begin()
	if beginError != nil {
		return nil, beginError
	}
	defer transaction.Rollback()
	statement, prepareError := transaction.Prepare("INSERT INTO urls(origin) VALUES($1) RETURNING id")
	if prepareError != nil {
		return nil, prepareError
	}
	defer statement.Close()
	ids := make([]uint, len(values))
	for i, value := range values {
		row := statement.QueryRow(value)
		if row.Err() != nil {
			return nil, row.Err()
		}
		row.Scan(&ids[i])
	}
	commitError := transaction.Commit()
	if commitError != nil {
		return nil, commitError
	}
	return ids, nil
}

func (database *PostgreSQLDatabase) GetAll(ids []uint) ([]string, error) {
	transaction, beginError := database.db.Begin()
	if beginError != nil {
		return nil, beginError
	}
	defer transaction.Rollback()
	statement, prepareError := transaction.Prepare("SELECT origin FROM urls WHERE id = $1")
	if prepareError != nil {
		return nil, prepareError
	}
	defer statement.Close()
	values := make([]string, len(ids))
	for i, id := range ids {
		row := statement.QueryRow(int64(id))
		if row.Err() != nil {
			return nil, row.Err()
		}
		scanError := row.Scan(&values[i])
		if scanError != nil {
			return nil, scanError
		}
	}
	commitError := transaction.Commit()
	if commitError != nil {
		return nil, commitError
	}
	return values, nil
}

func (database *PostgreSQLDatabase) Ping() error {
	return database.db.Ping()
}
