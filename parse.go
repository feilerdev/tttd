package main

import (
	"bufio"
	"fmt"
	"strings"
)

type TechnicalDebt struct {
	Type        string
	Description string
	File        string
	Line        int
}

// Parse extracts SATDs from a string content parsing it using a pre-agreed token.
// TODO: td-design > add 'token' param, so it the consumer can customize the content.
func Parse(content string) ([]*TechnicalDebt, error) {
	const (
		satdToken = "TODO"
		satdSep   = ":"
		satdPos   = 1
	)

	scanner := bufio.NewScanner(strings.NewReader(content))

	debts := make([]*TechnicalDebt, 0)

	var i int

	for scanner.Scan() {
		i++

		line := scanner.Text()

		if strings.Contains(line, satdToken) {
			satd := strings.Split(line, satdSep)[satdPos]

			td, err := extract(satd)
			if err != nil {
				return nil, fmt.Errorf("extracting content: %w", err)
			}

			if td == nil {
				continue
			}

			td.Line = i

			debts = append(debts, td)
		}
	}

	return debts, nil
}

func extract(satd string) (*TechnicalDebt, error) {
	const (
		tdSep   = ">"
		typePos = 0
		descPos = 1
	)

	tokens := strings.Split(satd, tdSep)

	if len(tokens) == 1 {
		return nil, nil
	}

	return &TechnicalDebt{
		Type:        strings.TrimSpace(tokens[typePos]),
		Description: strings.TrimSpace(tokens[descPos]),
	}, nil
}
