package main

import (
	"log/slog"
	"testing"
)

func TestExtractSATDs(t *testing.T) {
	NewLogger()

	// input
	workspacePath := "."
	ignorePath := "vendor"

	// process & output
	satds, err := extractSATDs(workspacePath, ignorePath)
	if err != nil {
		t.Fail()
	}

	// verify
	if len(satds) == 0 {
		t.Fail()
	}

	logger.Debug("satds", slog.Any("satds", satds))
}
