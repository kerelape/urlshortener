package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

func DialPostgreSQLDatabase(ctx context.Context, dsn string) (*PostgreSQLDatabase, error) {
	db, openError := sql.Open("pgx", dsn)
	if openError != nil {
		return nil, openError
	}
	_, execError := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS urls(id SERIAL NOT NULL PRIMARY KEY, origin TEXT UNIQUE)")
	return NewPostgreSQLDatabase(db), execError
}

func (database *PostgreSQLDatabase) Put(ctx context.Context, value string) (uint, error) {
	same := database.db.QueryRowContext(ctx, "SELECT id FROM urls WHERE origin = $1", value)
	if same.Err() != nil {
		return 0, same.Err()
	}
	var sameID int64
	if err := same.Scan(&sameID); err == nil {
		return 0, NewDuplicateValueError(uint(sameID))
	}
	row := database.db.QueryRowContext(ctx, "INSERT INTO urls(origin) VALUES($1) RETURNING id", value)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var id int64
	idError := row.Scan(&id)
	return uint(id), idError
}

func (database *PostgreSQLDatabase) Get(ctx context.Context, id uint) (string, error) {
	row := database.db.QueryRowContext(ctx, "SELECT origin FROM urls WHERE id = $1", int64(id))
	var origin string
	scanError := row.Scan(&origin)
	return origin, scanError
}

func (database *PostgreSQLDatabase) PutAll(ctx context.Context, values []string) ([]uint, error) {
	transaction, beginError := database.db.Begin()
	if beginError != nil {
		return nil, beginError
	}
	defer transaction.Rollback()
	statement, prepareError := transaction.PrepareContext(ctx, "INSERT INTO urls(origin) VALUES($1) RETURNING id")
	if prepareError != nil {
		return nil, prepareError
	}
	defer statement.Close()
	ids := make([]uint, len(values))
	for i, value := range values {
		row := statement.QueryRowContext(ctx, value)
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

func (database *PostgreSQLDatabase) GetAll(ctx context.Context, ids []uint) ([]string, error) {
	args := make([]string, 0, len(ids))
	for _, id := range ids {
		args = append(args, strconv.Itoa(int(id)))
	}
	rows, queryError := database.db.QueryContext(ctx, "SELECT origin FROM urls WHERE id IN $1", fmt.Sprintf("(%s)", strings.Join(args, ", ")))
	if queryError != nil {
		return nil, queryError
	}
	defer rows.Close()
	result := make([]string, 0, len(ids))
	for rows.Next() {
		var origin string
		if err := rows.Scan(&origin); err != nil {
			return nil, err
		}
		result = append(result, origin)
	}
	return result, nil
}

func (database *PostgreSQLDatabase) Ping(ctx context.Context) error {
	return database.db.PingContext(ctx)
}
