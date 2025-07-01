package main

import (
	"log/slog"
	"os"
	"path/filepath"
)

type config struct {
	workspacePath string
	ignorePath    string
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

	satds, err := extractSATDs(conf.workspacePath, conf.ignorePath)
	if err != nil {
		panic(err)
	}

	writeToCSV(logger, satds, conf.outputPath)
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

	conf.outputPath = os.Getenv("OUTPUT_PATH")
	conf.ignorePath = os.Getenv("IGNORE_PATH")
	if conf.outputPath == "" {
		conf.outputPath = defaultFileName
	}

	// TODO(alexandreliberato): Get CSV header from environment variable

	// If the outputPath is not absolute, make it relative to the workspace
	if !filepath.IsAbs(conf.outputPath) {
		conf.outputPath = filepath.Join(conf.workspacePath, conf.outputPath)
	}

	if !filepath.IsAbs(conf.ignorePath) {
		conf.ignorePath = filepath.Join(conf.workspacePath, conf.ignorePath)
	}

	return conf, nil
}
