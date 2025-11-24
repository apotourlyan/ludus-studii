package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/types"
	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/utilities"
)

func Create(args []string) {
	if len(args) == 0 {
		utilities.PrintCommandHelp(types.CommandCreate)
		return
	}

	migrationName := args[0]
	timestamp := time.Now().Format("20060102150405")

	// Create up migration
	upFile := fmt.Sprintf("%s_%s_up.sql", timestamp, migrationName)
	err := os.WriteFile(upFile, []byte("-- Migration up\n"), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("✓ Created %s\n", upFile)

	// Create down migration
	downFile := fmt.Sprintf("%s_%s_down.sql", timestamp, migrationName)
	err = os.WriteFile(downFile, []byte("-- Migration down\n"), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("✓ Created %s\n", downFile)
}
