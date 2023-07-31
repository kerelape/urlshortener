package storage

import (
	"context"
	"database/sql"
	"encoding/base32"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kerelape/urlshortener/internal/app"
)

// PostgreSQLDatabase is a database that stores data in a Postgres database.
type PostgreSQLDatabase struct {
	db *sql.DB
}

// NewPostgreSQLDatabase returns a new PostgreSQLDatabase.
func NewPostgreSQLDatabase(db *sql.DB) *PostgreSQLDatabase {
	return &PostgreSQLDatabase{
		db: db,
	}
}

// DialPostgreSQLDatabase connects to a Postgres database by the dsn
// and initializes it.
func DialPostgreSQLDatabase(ctx context.Context, dsn string) (*PostgreSQLDatabase, error) {
	db, openError := sql.Open("pgx", dsn)
	if openError != nil {
		return nil, openError
	}
	_, execError := db.ExecContext(
		ctx,
		`
		CREATE TABLE IF NOT EXISTS urls(
			id SERIAL NOT NULL PRIMARY KEY,
			origin TEXT UNIQUE,
			"user" TEXT,
			deleted BOOLEAN DEFAULT FALSE
		)
		`,
	)
	return NewPostgreSQLDatabase(db), execError
}

// Put stores value and returns its id.
func (database *PostgreSQLDatabase) Put(ctx context.Context, user app.Token, value string) (uint, error) {
	if database.db == nil {
		return 0, ErrDatabaseClosed
	}

	same := database.db.QueryRowContext(ctx, "SELECT id FROM urls WHERE origin = $1", value)
	if same.Err() != nil {
		return 0, same.Err()
	}
	var sameID int64
	if err := same.Scan(&sameID); err == nil {
		return 0, NewDuplicateValueError(uint(sameID))
	}
	row := database.db.QueryRowContext(
		ctx,
		`INSERT INTO urls(origin, "user") VALUES($1, $2) RETURNING id`,
		value,
		base32.StdEncoding.EncodeToString(user[:]),
	)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var id int64
	idError := row.Scan(&id)
	return uint(id), idError
}

// Get returns original value by its id.
func (database *PostgreSQLDatabase) Get(ctx context.Context, id uint) (string, error) {
	if database.db == nil {
		return "", ErrDatabaseClosed
	}

	row := database.db.QueryRowContext(ctx, "SELECT origin, deleted FROM urls WHERE id = $1", int64(id))
	var origin string
	var deleted bool
	if err := row.Scan(&origin, &deleted); err != nil {
		return "", err
	}
	if deleted {
		return "", ErrValueDeleted
	}
	return origin, nil
}

// PutAll stores values and returns their ids.
func (database *PostgreSQLDatabase) PutAll(ctx context.Context, user app.Token, values []string) ([]uint, error) {
	if database.db == nil {
		return nil, ErrDatabaseClosed
	}

	transaction, beginError := database.db.Begin()
	if beginError != nil {
		return nil, beginError
	}
	defer transaction.Rollback()
	statement, prepareError := transaction.PrepareContext(ctx, `INSERT INTO urls(origin, "user") VALUES($1, $2) RETURNING id`)
	if prepareError != nil {
		return nil, prepareError
	}
	defer statement.Close()
	ids := make([]uint, len(values))
	for i, value := range values {
		row := statement.QueryRowContext(ctx, value, base32.StdEncoding.EncodeToString(user[:]))
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

// GetAll returns original values by their ids.
func (database *PostgreSQLDatabase) GetAll(ctx context.Context, ids []uint) ([]string, error) {
	if database.db == nil {
		return nil, ErrDatabaseClosed
	}

	args := make([]string, 0, len(ids))
	for _, id := range ids {
		args = append(args, strconv.Itoa(int(id)))
	}
	rows, queryError := database.db.QueryContext(ctx, "SELECT origin FROM urls WHERE id IN $1", fmt.Sprintf("(%s)", strings.Join(args, ", ")))
	if queryError != nil {
		return nil, queryError
	}
	if rows.Err() != nil {
		return nil, rows.Err()
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

// Delete removes values by their ids.
func (database *PostgreSQLDatabase) Delete(ctx context.Context, user app.Token, ids []uint) error {
	if database.db == nil {
		return ErrDatabaseClosed
	}

	transaction, beginError := database.db.Begin()
	if beginError != nil {
		return beginError
	}
	defer transaction.Rollback()
	statement, prepareError := transaction.PrepareContext(
		ctx,
		`UPDATE urls SET deleted = $1 WHERE "user" = $2 AND id = $3`,
	)
	if prepareError != nil {
		return prepareError
	}
	defer statement.Close()
	for _, id := range ids {
		row := statement.QueryRowContext(ctx, true, base32.StdEncoding.EncodeToString(user[:]), id)
		if row.Err() != nil {
			return row.Err()
		}
	}
	commitError := transaction.Commit()
	if commitError != nil {
		return commitError
	}
	return nil
}

// URLs returns count of shorten urls.
func (database *PostgreSQLDatabase) URLs(ctx context.Context) (int, error) {
	if database.db == nil {
		return -1, ErrDatabaseClosed
	}
	rows, err := database.db.QueryContext(ctx, "SELECT id FROM urls")
	if err != nil {
		return -1, err
	}
	var r int
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return -1, err
		}
		r++
	}
	if err := rows.Close(); err != nil {
		return -1, err
	}
	return r, nil
}

// Users returns count of users registered in the database.
func (database *PostgreSQLDatabase) Users(ctx context.Context) (int, error) {
	if database.db == nil {
		return -1, ErrDatabaseClosed
	}
	rows, err := database.db.QueryContext(ctx, `SELECT "user" FROM urls`)
	if err != nil {
		return -1, err
	}
	var r = make(map[string]struct{})
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return -1, err
		}
		var user string
		if err := rows.Scan(&user); err != nil {
			return -1, err
		}
		r[user] = struct{}{}
	}
	if err := rows.Close(); err != nil {
		return -1, err
	}
	return len(r), nil
}

// Ping returns an error if the database is unavailable.
func (database *PostgreSQLDatabase) Ping(ctx context.Context) error {
	if database.db == nil {
		return ErrDatabaseClosed
	}
	return database.db.PingContext(ctx)
}

// Close closes this database.
func (database *PostgreSQLDatabase) Close(context.Context) error {
	conn := database.db
	database.db = nil
	return conn.Close()
}
