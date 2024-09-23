package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const satdPattern = `((\/\/)|(\*))\s?TODO(\([0-9A-Za-z].*\))?:?\s?[0-9A-Za-z_\s]*\s?(->)?\s?[0-9A-Za-z]*\-?[0-9A-Za-z]*\s?(=>)?\s?\$?\$?\$?\$?\$?(\w|$)`

var (
	ErrNotFound = errors.New("not SATDs found")
	regex       = regexp.MustCompile(satdPattern)
)

func extractSATDs(workspaceDir string) ([]TechnicalDebt, error) {
	satds := make([]TechnicalDebt, 0)

	// TODO: detect which files have satds using .Walk (sync)
	// put the content in a map and parse using go routines (concurrently)?
	err := filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		// logger.Info(fmt.Sprintf("%d", len(satds)))
		if info.IsDir() {
			return nil
		}

		newSatds, errDetectAndParse := detectAndParse(path, info, err)
		if errDetectAndParse != nil {
			return fmt.Errorf("detecting and parsing: %w", errDetectAndParse)
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

func detectAndParse(path string, info os.FileInfo, err error) ([]TechnicalDebt, error) {
	if err != nil {
		return nil, err
	}

	satds := make([]TechnicalDebt, 0)

	if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
		return satds, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("reading file", slog.Any("error", err))

		return nil, err
	}

	strContent := string(content)

	match := regex.MatchString(strContent)
	if !match {
		return satds, nil
	}

	satds, err = ParseRegex(strContent, info.Name())
	if err != nil {
		logger.Error("parsing", slog.Any("error", err))

		return nil, err
	}

	return satds, nil
}
