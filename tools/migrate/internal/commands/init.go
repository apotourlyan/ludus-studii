package commands

import (
	"fmt"

	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/services/datasys"
	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/types"
	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/utilities"
)

func Initialize(args []string) {
	if len(args) == 0 {
		utilities.PrintCommandHelp(types.CommandInit)
		return
	}

	connectionString := args[0]

	db, err := utilities.ConnectToDb(connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = datasys.CreateMigrationsTable(db)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("✓ Migrations table created successfully")
}
