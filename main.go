package main

import (
	"log/slog"
	"os"

	"github.com/akaspin/meisterwerk/cmd"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
		slog.Default().Error("app", "error", err)
		os.Exit(2)
	}
}
