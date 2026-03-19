package main

import (
	"log/slog"
	"os"
)

func main() {
	// 1. JSON Handler for stdout
	jsonH := slog.NewJSONHandler(os.Stdout, nil)

	// 2. Text Handler for stderr
	textH := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	// CORRECTED: Use the constructor NewMultiHandler
	multi := slog.NewMultiHandler(jsonH, textH)
	logger := slog.New(multi)

	// This goes to stdout only (Info level)
	logger.Info("App starting", "version", "1.26")
	// This goes to stdout AND stderr (Error level)
	logger.Error("Database down", "db", "primary")
}
