package main

import (
	"os"

	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/commands"
	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/types"
	"github.com/apotourlyan/ludus-studii/tools/migrate/internal/utilities"
)

func main() {
	if len(os.Args) < 2 {
		utilities.PrintGeneralHelp()
		return
	}

	c := types.Command(os.Args[1])
	switch c {
	case types.CommandInit:
		commands.Initialize(os.Args[2:])
	case types.CommandCreate:
		commands.Create(os.Args[2:])
	case types.CommandUp:
		commands.MigrateUp(os.Args[2:])
	case types.CommandDown:
		commands.MigrateDown(os.Args[2:])
	default:
		utilities.PrintGeneralHelp()
	}
}
