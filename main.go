package main

import (
	"github.com/g4stly/gofast/commands"
	"github.com/g4stly/gofast/common"
	"os"
)

func main() {
	if len(common.Args) < 1 {
		common.Fatal("see usage")
	}

	commandWord := common.Args[0]
	command, ok := commands.Lookup[commandWord]
	if !ok {
		common.Fatal("Invalid Command: %v", commandWord)
	}
	os.Exit(command.Exec(common.Args[1:]))
}
