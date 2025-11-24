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

func MigrateDown(args []string) {
	if len(args) < 2 {
		utilities.PrintCommandHelp(types.CommandDown)
		return
	}

	connectionString := args[0]
	target := args[1] // target migration

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

	entries, err := filesys.GetRollbackFiles(target, migrations)
	if err != nil {
		fmt.Println(err)
		return
	}

	slices.SortFunc(entries, func(a os.DirEntry, b os.DirEntry) int {
		// Descending order (newest first)
		return strings.Compare(b.Name(), a.Name())
	})

	err = datasys.RollbackMigrations(db, entries)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Rollback completed successfully")
}
