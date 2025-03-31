package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

const (
	patternAuthor      = `TODO\(([0-9A-Za-z\.]*)\)`
	patternDescription = `(TODO|\)|:)(\s)?([0-9A-Za-z\(\)'"_\-\s]*)(\s)?`
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

// ParseRegex parses the content of a file and extracts technical debts based on the defined regex patterns.
func ParseRegex(content string, path string, fileName string) ([]TechnicalDebt, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	debts := make([]TechnicalDebt, 0)

	var n int

	for scanner.Scan() {
		n++

		line := scanner.Text()

		if !regex.MatchString(line) {
			continue
		}

		var strAuthor string

		author := regexAuthor.FindAllStringSubmatch(line, -1)

		// author exists
		if len(author) > 0 {
			subAuthor := author[0]

			strAuthor = strings.TrimSpace(subAuthor[1])
		}

		var strDescription string

		description := regexDescription.FindAllStringSubmatch(line, 3)

		// verify if description exists and is not empty
		if len(description) > 0 {
			var subDesc []string

			if len(description) > 2 {
				subDesc = description[2]
			} else {
				subDesc = description[0]

				if len(description) == 2 {
					subDesc = description[1]
				}
			}

			if len(subDesc) > 2 {
				strDescription = strings.TrimSpace(subDesc[0])
			}

			// clean
			strDescription = strings.TrimLeft(strDescription, ":")
			strDescription = strings.TrimRight(strDescription, "-")
			strDescription = strings.TrimSpace(strDescription)
		}

		// verify if type and cost exist
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

		// create the TechnicalDebt struct
		file := fmt.Sprintf("%s/%s", path, fileName)

		td := TechnicalDebt{
			Author:      strAuthor,
			Description: strDescription,
			Type:        strType,
			Cost:        strCost,
			File:        file,
			Line:        n,
		}

		debts = append(debts, td)
	}

	return debts, nil
}
