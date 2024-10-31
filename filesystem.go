package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func writeToCSV(l *slog.Logger, satds []TechnicalDebt, path string) error {
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

	// It's based on Jira default fields
	records := [][]string{
		{"summary", "description", "reporter", "issue type"},
	}

	for _, satd := range satds {
		desc := fmt.Sprintf("%s | file:%s, line: %d", satd.Description, satd.File, satd.Line)
		// TODO(alexandre.liberato): add cost
		row := []string{satd.Description, desc, satd.Author, satd.Type}

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

func GetFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("reading file", slog.Any("err", err))

		return "", fmt.Errorf("reading file: %w", err)
	}

	return string(content), nil
}
