package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const satdPattern = `((\/\/)|(\*))\s?TODO(\([0-9A-Za-z].*\))?:?\s?[0-9A-Za-z_\s]*\s?(->)?\s?[0-9A-Za-z]*\-?[0-9A-Za-z]*\s?(=>)?\s?\$?\$?\$?\$?\$?(\w|$)`

var regex = regexp.MustCompile(satdPattern)

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

	// TODO: Get CSV header from environment variable

	// If the outputPath is not absolute, make it relative to the workspace
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(workspaceDir, outputPath)
	}

	logger.Info("output path", slog.String("path", outputPath))

	satds := make([]TechnicalDebt, 0)

	err = filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		return detectAndParse(path, info, err, satds)
	})
	if err != nil {
		logger.Error("reading files", slog.Any("error", err))

		panic(err)
	}

	if len(satds) == 0 {
		logger.Info("no valid SATDs in file")
	}

	debug(logger, satds)

	writeToCSV(logger, satds, outputPath)
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

func debug(l *slog.Logger, satds []TechnicalDebt) {
	for _, satd := range satds {
		l.Info(fmt.Sprintf("SATD Description: %s", satd.Description))
		l.Info(fmt.Sprintf("SATD Type: %s", satd.Type))
		l.Info("SATD Line:", slog.Int("line", satd.Line))
	}
}
