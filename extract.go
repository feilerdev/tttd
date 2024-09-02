package main

import (
	"errors"
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

	err := filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		return detectAndParse(path, info, err, satds)
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

func detectAndParse(path string, info os.FileInfo, err error, satds []TechnicalDebt) error {
	if err != nil {
		return err
	}

	if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("reading file", slog.Any("error", err))

		panic(err)
	}

	strContent := string(content)

	match := regex.MatchString(strContent)
	if !match {
		return nil
	}

	fileSatds, err := ParseRegex(strContent, info.Name())
	if err != nil {
		logger.Error("parsing", slog.Any("error", err))

		panic(err)
	}

	satds = append(satds, fileSatds...)

	return nil
}
