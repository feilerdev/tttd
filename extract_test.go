package main

import (
	"log/slog"
	"testing"
)

func TestExtractSATDs(t *testing.T) {
	NewLogger()

	// input
	workspace := "./tests/"

	// process & output
	satds, err := extractSATDs(workspace)
	if err != nil {
		t.Fail()
	}

	// verify
	if len(satds) == 0 {
		t.Fail()
	}

	logger.Debug("satds", slog.Any("satds", satds))
}
