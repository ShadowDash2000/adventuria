package main

import (
	"log/slog"
	"os"
)

var commands = map[string]func() error{
	"create_effect": createEffect,
}

func main() {
	if len(os.Args) < 3 {
		slog.Error("Usage: <command> <type> ...")
		os.Exit(1)
	}

	cmdName := os.Args[1] + "_" + os.Args[2]
	if cmd, ok := commands[cmdName]; ok {
		if err := cmd(); err != nil {
			slog.Error("Command failed", "error", err)
			os.Exit(1)
		}
	} else {
		slog.Error("Unknown command", "command", cmdName)
		os.Exit(1)
	}
}
