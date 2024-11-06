package main

import (
	"bufio"
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

		var strAuthor string

		author := regexAuthor.FindAllStringSubmatch(line, -1)

		// author exists
		if len(author) > 0 {
			subAuthor := author[0]

			strAuthor = strings.TrimSpace(subAuthor[1])
		}

		var strDesc string

		desc := regexDescription.FindAllStringSubmatch(line, 3)

		// description exists
		if len(desc) > 0 {
			var subDesc []string

			if len(desc) > 2 {
				subDesc = desc[2]
			} else {
				subDesc = desc[0]

				if len(desc) == 2 {
					subDesc = desc[1]
				}
			}

			if len(subDesc) > 2 {
				strDesc = strings.TrimSpace(subDesc[0])
			}

			// clean
			strDesc = strings.TrimLeft(strDesc, ":")
			strDesc = strings.TrimRight(strDesc, "-")
			strDesc = strings.TrimSpace(strDesc)
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
