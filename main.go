package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	NewLogger()

	// The workspace directory is automatically set as the working directory
	workspaceDir, _ := os.Getwd()

	logger.Info(workspaceDir)

	// Get the output path from the environment variable
	outputPath := os.Getenv("OUTPUT_PATH")
	if outputPath == "" {
		outputPath = "satds.csv"
	}

	outputPath = filepath.Join(workspaceDir, outputPath)

	logger.Info(outputPath)

	satds := make([]*TechnicalDebt, 0)

	err := filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
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

		fileSatds, err := Parse(string(content))
		if err != nil {
			logger.Error("parsing", slog.Any("error", err))

			panic(err)
		}

		satds = append(satds, fileSatds...)

		return nil
	})
	if err != nil {
		logger.Error("reading files", slog.Any("error", err))

		panic(err)
	}

	if len(satds) == 0 {
		logger.Info("No valid SATDs in file")
	}

	debug(logger, satds)

	writeToCSV(logger, satds, outputPath)
}

func writeToCSV(l *slog.Logger, satds []*TechnicalDebt, path string) error {
	l.Info(fmt.Sprintf("Writing CSV to %s", path))

	// Ensure the directory exists
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		l.Error("failed creating directory", slog.Any("error", err))
		return err
	}

	csvFile, err := os.Create(path)
	if err != nil {
		l.Error("failed creating file", slog.Any("error", err))

		return err
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)

	records := [][]string{
		{"summary", "description", "reporter", "issue type"},
	}

	for _, satd := range satds {
		desc := fmt.Sprintf("%s, line: %d", satd.Description, satd.Line)
		row := []string{satd.Description, desc, "tttd", satd.Type}

		records = append(records, row)
	}

	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		l.Error("error writing csv", slog.Any("error", err))

		return err
	}

	l.Info("CSV file written successfully", slog.String("path", path))

	return nil
}

func debug(l *slog.Logger, satds []*TechnicalDebt) {
	for _, satd := range satds {
		l.Info(fmt.Sprintf("SATD Description: %s", satd.Description))
		l.Info(fmt.Sprintf("SATD Type: %s", satd.Type))
		l.Info("SATD Line:", slog.Int("line", satd.Line))
	}
}
