package db

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func Connect(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	for _, migration := range migrations {
		_, err = db.Exec(migration)
		if err != nil {
			return nil, fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func UUID() string {
	return uuid.Must(uuid.NewV4()).String()
}
