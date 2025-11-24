package filesys

import (
	"os"
	"slices"
	"strings"
)

func GetMigrationFiles(
	targetMigration string, appliedMigrations []string,
) ([]os.DirEntry, error) {
	var migrationEntries []os.DirEntry
	entries, err := os.ReadDir(".")
	if err != nil {
		return migrationEntries, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, "_up.sql") {
			continue
		}

		name, _ = strings.CutSuffix(name, "_up.sql")
		shouldAppend := name <= targetMigration &&
			!slices.Contains(appliedMigrations, name)
		if shouldAppend {
			migrationEntries = append(migrationEntries, entry)
		}
	}

	return migrationEntries, nil
}

func GetRollbackFiles(
	targetMigration string, appliedMigrations []string,
) ([]os.DirEntry, error) {
	var rollbackEntries []os.DirEntry
	entries, err := os.ReadDir(".")
	if err != nil {
		return rollbackEntries, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, "_down.sql") {
			continue
		}

		name, _ = strings.CutSuffix(name, "_down.sql")
		shouldAppend := name >= targetMigration &&
			slices.Contains(appliedMigrations, name)
		if shouldAppend {
			rollbackEntries = append(rollbackEntries, entry)
		}
	}

	return rollbackEntries, nil
}
