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

// extractSATDs walks through the workspace directory and extracts all SATDs from the Go files.
// It returns a slice of TechnicalDebt and an error if any occurs during the process.
// The function uses a recursive approach to traverse the directory structure and
// identifies SATDs based on the defined pattern. The extracted SATDs are then parsed
// and returned as a slice of TechnicalDebt.
func extractSATDs(workspaceDir string, ignorePaths []string) ([]TechnicalDebt, error) {
	satds := make([]TechnicalDebt, 0)

	// TODO(alexandreliberato): detect which files have satds using .Walk with sync and
	// put the content in a map and parse using go routines, maybe concurrently?
	err := filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ignore(path, ignorePaths) {
			if info.IsDir() {
				logger.Debug("skipping ignored directory", slog.String("path", path))

				return filepath.SkipDir
			}

			return nil
		}

		if info.IsDir() {
			return nil
		}

		content, errDetectAndParse := detect(path, info, err)
		if errDetectAndParse != nil {
			return fmt.Errorf("detecting and parsing: %w", errDetectAndParse)
		}

		// ParseRegex joins the directory and file name, so pass the file's
		// directory (not the full path) to avoid duplicating the file name.
		// TODO(alexandreliberato): use goroutines to parse the files concurrently
		newSatds, err := ParseRegex(content, filepath.Dir(path), info.Name())
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

// ignore reports whether path p is, or lives under, any of the ignored paths.
// It matches on full path segments so it works for relative top-level paths
// (e.g. "vendor/x.go"), nested paths ("pkg/vendor/x.go") and absolute paths.
func ignore(p string, ignorePaths []string) bool {
	p = filepath.ToSlash(p)

	for _, ip := range ignorePaths {
		if ip == "" {
			continue
		}

		ip = filepath.ToSlash(ip)

		if p == ip ||
			strings.HasPrefix(p, ip+"/") ||
			strings.Contains(p, "/"+ip+"/") {
			return true
		}
	}

	return false
}
