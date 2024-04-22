package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// process the input file and return the output string.
func processInputFile(inputFile *os.File, csvFile *os.File, iataIndex, icaoIndex, nameIndex int) (string, error) {
	scanner := bufio.NewScanner(inputFile)
	iataRegex, _ := regexp.Compile(`#[A-Z]{3}`)
	icaoRegex, _ := regexp.Compile(`##[A-Z]{4}`)
	dateRegex, _ := regexp.Compile(`([DT])(\d{2})?\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(Z|[\+\-]\d{2}:\d{2}))\)`)
	var output []string
	var lastLineWasBlank bool

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "\v", "\n")
		line = strings.ReplaceAll(line, "\f", "\n")
		line = strings.ReplaceAll(line, "\r", "\n")

		line = strings.TrimSpace(line)

		if line == "" {
			if !lastLineWasBlank {
				output = append(output, "\n")
				lastLineWasBlank = true
			}
		} else {
			lastLineWasBlank = false
			matches := findAllMatches(line, iataRegex, icaoRegex, dateRegex)
			output = append(output, processMatches(matches, line, csvFile, iataIndex, icaoIndex, nameIndex))
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading from input file: %v", err)
	}

	return strings.Join(output, ""), nil
}

// find all IATA and ICAO code and date matches.
func findAllMatches(line string, iataRegex, icaoRegex, dateRegex *regexp.Regexp) []Match {
	var matches []Match

	iataMatches := iataRegex.FindAllStringIndex(line, -1)
	for _, match := range iataMatches {
		matches = append(matches, Match{Index: match[0], Value: line[match[0]:match[1]], Type: "iata"})
	}

	icaoMatches := icaoRegex.FindAllStringIndex(line, -1)
	for _, match := range icaoMatches {
		matches = append(matches, Match{Index: match[0], Value: line[match[0]:match[1]], Type: "icao"})
	}

	dateMatches := dateRegex.FindAllStringIndex(line, -1)
	for _, match := range dateMatches {
		matches = append(matches, Match{Index: match[0], Value: line[match[0]:match[1]], Type: "date"})
	}

	return matches
}

// process matches and return the result string.
func processMatches(matches []Match, line string, csvFile *os.File, iataIndex, icaoIndex, nameIndex int) string {
	// sort matches by index
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Index < matches[j].Index
	})

	for _, match := range matches {
		var replacement string
		var err error
		switch match.Type {
		case "iata", "icao":
			replacement = lookupCode(match.Value, csvFile, iataIndex, icaoIndex, nameIndex)
		case "date":
			replacement, err = processLine(match.Value) // Get replacement or error
			if err != nil {
				// If there's an error, leave the original value
				replacement = match.Value
			}
		}

		// Replace the original value with the processed one in the line
		line = strings.Replace(line, match.Value, replacement, 1)
	}

	return line + "\n"
}
