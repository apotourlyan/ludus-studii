package dbutil

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/apotourlyan/ludus-studii/pkg/panicutil"
)

func Initialize(connectionString string) *sql.DB {
	panicutil.RequireNotEmptyOrWhitespace(connectionString, "connection string")

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully")
	return db
}
