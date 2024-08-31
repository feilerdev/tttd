package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type TechnicalDebt struct {
	Author      string
	Description string
	Type        string
	File        string
	Line        int
	sad         string
}

const (
	patternAuthor      = `TODO(\([0-9A-Za-z].*\))`
	patternDescription = `((TODO)|(:)|(\))\s?)[0-9A-Za-z_\s]*\s?`
	patternType        = `(->)\s?[0-9A-Za-z_\s-]*\s`
	patternCost        = `(=>)\s?[$]*`
)

var (
	regexAuthor      = regexp.MustCompile(patternAuthor)
	regexDescription = regexp.MustCompile(patternDescription)
	regexType        = regexp.MustCompile(patternType)
	regexCost        = regexp.MustCompile(patternCost)
)

// Parse extracts SATDs from a string content parsing it using a pre-agreed token.
func Parse(content string) ([]TechnicalDebt, error) {
	const (
		satdToken = "TODO"
		satdSep   = ":"
		satdPos   = 1
	)

	scanner := bufio.NewScanner(strings.NewReader(content))

	debts := make([]TechnicalDebt, 0)

	var i int

	for scanner.Scan() {
		i++

		line := scanner.Text()

		if strings.Contains(line, satdToken) {
			tokens := strings.Split(line, satdSep)

			if len(tokens) <= satdPos {
				continue
			}

			satd := tokens[satdPos]

			td, err := extract(satd)
			if err != nil {
				logger.Error("parsing", err)

				return nil, fmt.Errorf("extracting content: %w", err)
			}

			if td.File == "" {
				continue
			}

			td.Line = i

			debts = append(debts, td)
		}
	}

	return debts, nil
}

func ParseRegex(content string, file string) ([]TechnicalDebt, error) {

	scanner := bufio.NewScanner(strings.NewReader(content))

	debts := make([]TechnicalDebt, 0)

	var n int

	for scanner.Scan() {
		n++

		line := scanner.Text()

		if !regex.MatchString(line) {
			continue
		}

		author := regexAuthor.FindString(line)
		desc := regexAuthor.FindString(line)
		tdType := regexAuthor.FindString(line)

		td := TechnicalDebt{
			Author:      author,
			Description: desc,
			Type:        tdType,
			File:        file,
			Line:        n,
		}

		debts = append(debts, td)
	}

	return debts, nil
}

func extract(satd string) (TechnicalDebt, error) {
	const (
		tdSep   = ">"
		typePos = 0
		descPos = 1
	)

	tokens := strings.Split(satd, tdSep)

	if len(tokens) == 1 {
		return TechnicalDebt{}, nil
	}

	return TechnicalDebt{
		Type:        strings.TrimSpace(tokens[typePos]),
		Description: strings.TrimSpace(tokens[descPos]),
	}, nil
}
