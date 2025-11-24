package commands

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/services/datasys"
	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/services/filesys"
	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/types"
	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/utilities"
)

func MigrateUp(args []string) {
	if len(args) < 1 {
		utilities.PrintCommandHelp(types.CommandUp)
		return
	}

	connectionString := args[0]
	var target string
	if len(args) >= 2 {
		target = args[1] // Specific migration
	} else {
		// Trick to run all migrations that are not applied
		target = "99999999999999"
	}

	db, err := utilities.ConnectToDb(connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	migrations, err := datasys.GetMigrations(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	entries, err := filesys.GetMigrationFiles(target, migrations)
	if err != nil {
		fmt.Println(err)
		return
	}

	slices.SortFunc(entries, func(a os.DirEntry, b os.DirEntry) int {
		// Ascending order (oldest first)
		return strings.Compare(a.Name(), b.Name())
	})

	err = datasys.ApplyMigrations(db, entries)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Migration completed successfully")
}
