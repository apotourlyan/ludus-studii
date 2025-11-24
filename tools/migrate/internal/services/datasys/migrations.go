package datasys

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strings"
)

//go:embed scripts/create_migrations_table.sql
var createTableScript string

//go:embed scripts/get_migrations.sql
var getMigrationsScript string

//go:embed scripts/delete_migration.sql
var deleteMigrationScript string

//go:embed scripts/insert_migration.sql
var insertMigrationScript string

func CreateMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(createTableScript)
	return err
}

func GetMigrations(db *sql.DB) ([]string, error) {
	var migrations []string
	rows, err := db.Query(getMigrationsScript)
	if err != nil {
		return migrations, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return migrations, err
		}
		migrations = append(migrations, name)
	}

	return migrations, rows.Err()
}

func ApplyMigrations(db *sql.DB, entries []os.DirEntry) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		filename := entry.Name()
		migration, _ := strings.CutSuffix(filename, "_up.sql")

		// Read the migration SQL file
		content, err := os.ReadFile(filename)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Failed to read %s\n", filename)
			return err
		}

		// Execute the migration SQL
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			fmt.Printf("Failed to execute %s\n", filename)
			return err
		}

		// Insert migration record into schema_migrations
		if _, err := tx.Exec(insertMigrationScript, migration); err != nil {
			tx.Rollback()
			fmt.Printf("Failed to insert migration record %s\n", migration)
			return err
		}

		fmt.Printf("Applied: %s\n", migration)
	}

	return tx.Commit()
}

func RollbackMigrations(db *sql.DB, entries []os.DirEntry) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		filename := entry.Name()
		migration, _ := strings.CutSuffix(filename, "_down.sql")

		// Read the rollback SQL file
		content, err := os.ReadFile(filename)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Failed to read %s\n", filename)
			return err
		}

		// Execute the rollback SQL
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			fmt.Printf("Failed to execute %s\n", filename)
			return err
		}

		// Remove migration record from schema_migrations
		if _, err := tx.Exec(deleteMigrationScript, migration); err != nil {
			tx.Rollback()
			fmt.Printf("Failed to remove migration record %s\n", migration)
			return err
		}

		fmt.Printf("Rolled back: %s\n", migration)
	}

	return tx.Commit()
}
