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

	// fileListPath := os.Args[1]

	// The workspace directory is automatically set as the working directory
	workspaceDir, _ := os.Getwd()

	logger.Info(workspaceDir)

	// Get the output path from the environment variable
	outputPath := os.Getenv("OUTPUT_PATH")

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
			logger.Error("reading file: %w", err)
			panic(err)
		}

		fileSatds, err := Parse(string(content))
		if err != nil {
			logger.Error("parsing: %w", err)
			panic(err)
		}

		satds = append(satds, fileSatds...)

		return nil
	})
	if err != nil {
		logger.Error("reading files: %", err)
		panic(err)
	}

	if len(satds) == 0 {
		logger.Info("No valid SATDs in file")
	}

	debug(logger, satds)

	export(logger, satds, outputPath)
}

func export(l *slog.Logger, satds []*TechnicalDebt, path string) error {
	logger.Info(path + "/satds.csv")
	csvFile, err := os.Create(path + "/satds.csv")
	if err != nil {
		l.Error("failed creating file: %s", err)

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
		l.Error("error writing csv:", err)

		return err
	}

	return nil
}

func debug(l *slog.Logger, satds []*TechnicalDebt) {
	for _, satd := range satds {
		l.Info(satd.Type)
		l.Info(satd.Description)
		l.Info("line", slog.Int("line", satd.Line))
	}
}
