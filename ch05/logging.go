package main

import (
	"io"
	"log/slog"
	"log/syslog"
	"os"
)

// multiWriter wraps multiple io.Writers
type multiWriter struct {
	writers []io.Writer
}

func (m *multiWriter) Write(p []byte) (int, error) {
	var lastErr error
	for _, w := range m.writers {
		if _, err := w.Write(p); err != nil {
			// Capture the error but continue to other writers
			lastErr = err
		}
	}
	// Return len(p) to indicate data was processed,
	// even if some sinks failed.
	return len(p), lastErr
}

func main() {
	// Open a JSON log file
	file, err := os.OpenFile("/tmp/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Syslog writer
	syslogWriter, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "myapp")
	if err != nil {
		panic(err)
	}
	defer syslogWriter.Close()

	// Combine stdout, file, and syslog writers
	w := &multiWriter{
		writers: []io.Writer{
			file,
			os.Stdout,
			syslogWriter,
		},
	}

	// Create JSON handler with empty HandlerOptions
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{})

	// Create logger
	logger := slog.New(handler)

	// Example logs
	logger.Info("Application started")
	logger.Debug("Processing request",
		"id", 12345,
		"user", "alice",
		"attempt", 1,
	)

	err = doSomething()
	if err != nil {
		logger.Error("Failed to process data",
			"error", err,
		)
	}

	logger.Info("Application finished successfully")
}

// Simulate a function that may return an error
func doSomething() error {
	// For demonstration, return nil or a simple error
	return nil
}
