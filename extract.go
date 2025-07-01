package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
)

const satdPattern = `((\/\/)|(\*))\s?TODO(\([0-9A-Za-z].*\))?:?\s?[0-9A-Za-z_\s]*\s?(->)?\s?[0-9A-Za-z]*\-?[0-9A-Za-z]*\s?(=>)?\s?\$?\$?\$?\$?\$?(\w|$)`

var (
	ErrNotFound = errors.New("not SATDs found")
	regex       = regexp.MustCompile(satdPattern)
)

// extractSATDs walks through the workspace directory and extracts all SATDs from the Go files.
// It returns a slice of TechnicalDebt and an error if any occurs during the process.
// The function uses a recursive approach to traverse the directory structure and
// identifies SATDs based on the defined pattern. The extracted SATDs are then parsed
// and returned as a slice of TechnicalDebt.
func extractSATDs(workspaceDir, ignorePath string) ([]TechnicalDebt, error) {
	satds := make([]TechnicalDebt, 0)

	// TODO(alexandreliberato): detect which files have satds using .Walk with sync and
	// put the content in a map and parse using go routines, maybe concurrently?
	err := filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if info.IsDir() && info.Name() == ignorePath {
			logger.Error("skipping directory")

			return filepath.SkipDir
		}

		content, errDetectAndParse := detect(path, info, err)
		if errDetectAndParse != nil {
			return fmt.Errorf("detecting and parsing: %w", errDetectAndParse)
		}

		// TODO(alexandreliberato): use goroutines to parse the files concurrently
		newSatds, err := ParseRegex(content, path, info.Name())
		if err != nil {
			logger.Error("parsing", slog.Any("error", err))

			return fmt.Errorf("parsing: %w", err)
		}

		satds = append(satds, newSatds...)

		return nil
	})
	if err != nil {
		logger.Error("reading files", slog.Any("error", err))

		panic(err)
	}

	if len(satds) == 0 {
		logger.Info("no SATDs found")

		return nil, ErrNotFound
	}

	return satds, nil
}

func detect(path string, info os.FileInfo, err error) (string, error) {
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "", nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("reading file", slog.Any("error", err))

		return "", err
	}

	strContent := string(content)

	match := regex.MatchString(strContent)
	if !match {
		return "", nil
	}

	return strContent, nil
}
