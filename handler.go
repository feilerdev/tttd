package main

import (
	"encoding/json"
	"fmt"
)

func Decode(content string) ([]string, error) {
	m := []string{}

	err := json.Unmarshal([]byte(content), &m)
	if err != nil {
		logger.Error("unmarshaling json", err)

		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return nil, nil
}
