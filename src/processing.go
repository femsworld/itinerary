package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

// process the input file and return the output string.
func processInputFile(inputFile *os.File, csvFile *os.File, iataIndex, icaoIndex, nameIndex, cityIndex int) (formattedOutput, unformattedOutput string, err error) {
	// Declare and initialize the regex patterns
	iataRegex, _ := regexp.Compile(`#([A-Z]{3})`)
	icaoRegex, _ := regexp.Compile(`##([A-Z]{4})`)
	dateRegex, _ := regexp.Compile(`([DT])(\d{2})?\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(Z|[\+\-]\d{2}:\d{2}))\)`) // New regex pattern for date
	starRegex, _ := regexp.Compile(`\*#([A-Z]{3})`)

	// Read and process input
	scanner := bufio.NewScanner(inputFile)
	var formattedLines, unformattedLines []string
	var lastLineWasBlank bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := findAllMatches(line, iataRegex, icaoRegex, dateRegex, starRegex) // Include all four regex patterns

		if line == "" { // Handle blank lines
			if !lastLineWasBlank {
				formattedLines = append(formattedLines, "\n")
				unformattedLines = append(unformattedLines, "\n")
				lastLineWasBlank = true
			}
		} else {
			lastLineWasBlank = false
			// Process formatted and unformatted outputs
			formattedLines = append(formattedLines, processMatches(matches, line, csvFile, iataIndex, icaoIndex, nameIndex, cityIndex, true))
			unformattedLines = append(unformattedLines, processMatches(matches, line, csvFile, iataIndex, icaoIndex, nameIndex, cityIndex, false))
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("error reading from input file: %v", err)
	}

	// Join the outputs into single strings
	formattedOutput = strings.Join(formattedLines, "")
	unformattedOutput = strings.Join(unformattedLines, "")

	return formattedOutput, unformattedOutput, nil
}

// find all IATA and ICAO code and date matches.
func findAllMatches(line string, iataRegex, icaoRegex, dateRegex, starRegex *regexp.Regexp) []Match {
	var matches []Match

	// IATA matches
	iataMatches := iataRegex.FindAllStringIndex(line, -1)
	for _, match := range iataMatches {
		matches = append(matches, Match{Index: match[0], Value: line[match[0]:match[1]], Type: "iata"})
	}

	// ICAO matches
	icaoMatches := icaoRegex.FindAllStringIndex(line, -1)
	for _, match := range icaoMatches {
		matches = append(matches, Match{Index: match[0], Value: line[match[0]:match[1]], Type: "icao"})
	}

	// Starred matches
	starMatches := starRegex.FindAllStringIndex(line, -1)
	for _, match := range starMatches {
		matches = append(matches, Match{Index: match[0], Value: line[match[0]:match[1]], Type: "starred"})
	}

	// Date matches
	dateMatches := dateRegex.FindAllStringIndex(line, -1)
	for _, match := range dateMatches {
		matches = append(matches, Match{Index: match[0], Value: line[match[0]:match[1]], Type: "date"})
	}

	return matches
}

// process matches and return the result string.
func processMatches(matches []Match, line string, csvFile *os.File, iataIndex, icaoIndex, nameIndex, cityIndex int, formatOutput bool) string {
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Index < matches[j].Index
	})

	for _, match := range matches {
		var replacement string
		var err error
		var code string
		if strings.HasPrefix(match.Value, "*#") {
			// It's a city name lookup
			code = match.Value[2:] // Remove "*#"
			replacement = lookupCityName(code, csvFile, iataIndex, icaoIndex, cityIndex)
			if formatOutput {
				replacement = boldBlue + replacement + ansiReset // format for terminal
			}
		} else {
			switch match.Type {
			case "iata", "icao":
				replacement = lookupCode(match.Value, csvFile, iataIndex, icaoIndex, nameIndex)
				if formatOutput {
					replacement = boldYellow + replacement + ansiReset // format for terminal
				}
			case "date":
				replacement, err = processLine(match.Value)
				if err != nil {
					replacement = match.Value
				}
			}
		}

		line = strings.Replace(line, match.Value, replacement, 1) // replaces in the line
	}

	return line + "\n"
}

// lookup city name based on IATA/ICAO code
func lookupCityName(code string, csvFile *os.File, iataIndex, icaoIndex, cityIndex int) string {
	// Reset file position to the beginning to start reading from the start.
	csvFile.Seek(0, io.SeekStart)

	reader := csv.NewReader(csvFile)
	cityName := code // Default to code if not found

	// Iterate through the CSV file to find the matching IATA or ICAO code
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return cityName // Return the original code if there's an error
		}

		// If IATA or ICAO matches the code, retrieve the corresponding municipality
		if strings.TrimSpace(record[iataIndex]) == code || strings.TrimSpace(record[icaoIndex]) == code {
			if cityIndex < len(record) {
				cityName = record[cityIndex] // Get the corresponding municipality
			}
			break
		}
	}

	return cityName
}
