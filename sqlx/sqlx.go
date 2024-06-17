package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

const DefaultDriver = "sqlite"

type DB struct {
	*sql.DB
}

func Dial(ctx context.Context, driver, dsn string) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database driver %q by %q: %w", driver, dsn, err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &DB{DB: db}, nil
}

func IsPrimaryKeyViolation(err error) bool {
	var serr *sqlite.Error
	if errors.As(err, &serr) {
		return serr.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY
	}
	return false
}
