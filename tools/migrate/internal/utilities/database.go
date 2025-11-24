package utilities

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectToDb(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	// Verify connection actually works
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
