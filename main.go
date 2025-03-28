package main

import (
	"log/slog"
	"os"
	"path/filepath"
)

type config struct {
	workspaceDir string
	outputPath   string
}

func main() {
	NewLogger()

	execute()
}

func execute() {
	const defaultFileName = "satds.csv"

	// The workspace directory is automatically set as the working directory
	workspaceDir, err := os.Getwd()
	if err != nil {
		logger.Error("getting root directory", slog.Any("error", err))

		panic(err)
	}

	logger.Info(workspaceDir)

	// Get the output path from the environment variable
	outputPath := os.Getenv("INPUT_OUTPUT_PATH")
	if outputPath == "" {
		outputPath = defaultFileName
	}

	// TODO(alexandreliberato): Get CSV header from environment variable

	// If the outputPath is not absolute, make it relative to the workspace
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(workspaceDir, outputPath)
	}

	satds, err := extractSATDs(workspaceDir)
	if err != nil {
		panic(err)
	}

	writeToCSV(logger, satds, outputPath)
}
