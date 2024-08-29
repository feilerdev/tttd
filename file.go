package main

import (
	"fmt"
	"os"
)

func GetFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("reading file", err)

		return "", fmt.Errorf("reading file: %w", err)
	}

	return string(content), nil
}
