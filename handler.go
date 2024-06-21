package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func Decode(content string) ([]string, error) {
	m := []string{}

	err := json.Unmarshal([]byte(content), &m)
	if err != nil {
		fmt.Errorf("unmarshal: %w", err)
	}

	log.Printf("Unmarshaled: %v", m)

	return nil, nil
}
