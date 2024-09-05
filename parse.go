package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

const (
	patternAuthor      = `TODO\(([0-9A-Za-z\.]*)\)`
	patternDescription = `(TODO|\)|:)(\s)?([0-9A-Za-z_\-\s]*)(\s)?`
	patternType        = `->\s?([0-9A-Za-z_\s-]*)`
	patternCost        = `=>\s?([$]*)`
)

var (
	// TODO(alexandre.liberato): add date
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

		fmt.Printf("\n\n---->> %s\n", line)

		var strAuthor string

		author := regexAuthor.FindAllStringSubmatch(line, -1)
		// author exists
		if len(author) > 0 {
			subAuthor := author[0]

			strAuthor = strings.TrimSpace(subAuthor[1])
		}

		var strDesc string

		fmt.Printf("---->> description\n")
		desc := regexDescription.FindAllStringSubmatch(line, 3)
		// description exists
		if len(desc) > 0 {
			var subDesc []string

			fmt.Printf("---->> description[0]: %s\n", desc[0])
			if len(desc) > 1 {
				fmt.Printf("---->> description[1]: %s\n", desc[1])
			}

			if len(desc) > 2 {
				subDesc = desc[2]
			}

			if len(subDesc) > 2 {
				strDesc = strings.TrimSpace(subDesc[2])
			}
		}

		var strType string
		tdType := regexType.FindAllStringSubmatch(line, -1)
		if len(tdType) > 0 {
			subType := tdType[0]
			strType = strings.TrimSpace(subType[1])
		}

		var strCost string
		cost := regexCost.FindAllStringSubmatch(line, -1)
		if len(cost) > 0 {
			subCost := cost[0]
			strCost = strings.TrimSpace(subCost[1])
		}

		td := TechnicalDebt{
			Author:      strAuthor,
			Description: strDesc,
			Type:        strType,
			Cost:        strCost,
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
