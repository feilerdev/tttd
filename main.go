package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	workspacePath string
	ignorePaths   []string
	outputPath    string
}

func main() {
	NewLogger()

	execute()
}

func execute() {

	conf, err := loadConf()
	if err != nil {
		logger.Error("loading conf", slog.Any("err", err))
	}

	satds, err := extractSATDs(conf.workspacePath, conf.ignorePaths)
	if err != nil {
		// No SATDs is a valid outcome: emit an empty (headers-only) report
		// instead of failing the workflow.
		if !errors.Is(err, ErrNotFound) {
			panic(err)
		}

		logger.Info("no SATDs found; writing empty report")

		satds = nil
	}

	if err := writeToCSV(logger, satds, conf.outputPath); err != nil {
		logger.Error("writing CSV", slog.Any("error", err))

		os.Exit(1)
	}
}

func loadConf() (config, error) {
	const defaultFileName = "satds.csv"

	// Get paths from environment variables
	var (
		err  error
		conf config
	)

	conf.workspacePath = os.Getenv("WORKSPACE_PATH")

	if conf.workspacePath == "" {
		// use current path
		conf.workspacePath, err = os.Getwd()
		if err != nil {
			logger.Error("getting root directory", slog.Any("error", err))

			panic(err)
		}
	}

	conf.outputPath = fmt.Sprintf("%s/%s", os.Getenv("OUTPUT_PATH"), defaultFileName)

	// IGNORE_PATH accepts one or more comma-separated paths (e.g. "vendor,.github").
	// Each is matched as a path segment against the (relative) paths produced by
	// filepath.Walk, so keep them as given rather than absolutizing them.
	conf.ignorePaths = parseIgnorePaths(os.Getenv("IGNORE_PATH"))

	// TODO(alexandreliberato): Get CSV header from environment variable

	// If the outputPath is not absolute, make it relative to the workspace
	if !filepath.IsAbs(conf.outputPath) {
		conf.outputPath = filepath.Join(conf.workspacePath, conf.outputPath)
	}

	return conf, nil
}

// parseIgnorePaths splits a comma-separated IGNORE_PATH value into individual
// paths, trimming surrounding whitespace and dropping empty entries.
func parseIgnorePaths(raw string) []string {
	parts := strings.Split(raw, ",")

	paths := make([]string, 0, len(parts))

	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			paths = append(paths, p)
		}
	}

	return paths
}
