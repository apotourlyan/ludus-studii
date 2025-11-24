package utilities

import (
	"fmt"

	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/types"
)

func PrintGeneralHelp() {
	fmt.Println("Usage: migrate <command>")
	fmt.Println("Commands:")
	fmt.Println("  init              Create migrations table")
	fmt.Println("  create <name>     Create new migration files")
	fmt.Println("  up <name>         Run migration up")
	fmt.Println("  down <name>       Run migration down")
}

func PrintCommandHelp(c types.Command) {
	switch c {
	case types.CommandInit:
		fmt.Println("Usage: migrate init <connection-string>")
		fmt.Println(`       migrate init "postgres://user:pass@localhost:5432/dbname?option=value"`)
		fmt.Println()
		fmt.Println("Initializes the migrations system by creating the schema_migrations table.")

	case types.CommandCreate:
		fmt.Println("Usage: migrate create <migration-name>")
		fmt.Println(`       migrate create "create_users_table"`)
		fmt.Println()
		fmt.Println("Creates timestamped up and down migration files")

	case types.CommandUp:
		fmt.Println("Usage: migrate up <connection-string> [migration-name]")
		fmt.Println(`       migrate up "postgres://user:pass@localhost:5432/dbname"`)
		fmt.Println(`       migrate up "postgres://user:pass@localhost:5432/dbname" "20250123150405_create_users_table"`)
		fmt.Println()
		fmt.Println("Run migrations:")
		fmt.Println("  - With migration-name: runs all migrations up to and including the specified one")
		fmt.Println("  - Without migration-name: runs all pending migrations")

	case types.CommandDown:
		fmt.Println("Usage: migrate down <connection-string> <migration-name>")
		fmt.Println(`       migrate down "postgres://user:pass@localhost:5432/dbname" "20250123150405_create_users_table"`)
		fmt.Println()
		fmt.Println("Rollback migrations down to and including the specified migration")
	default:
		PrintGeneralHelp()
	}
}
