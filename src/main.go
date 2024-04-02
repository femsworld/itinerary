package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	// Parse command-line arguments
	if len(os.Args) != 4 || os.Args[1] == "-h" {
		// fmt.Println("Usage: go run . input.txt output.txt airport-lookup.csv")
		fmt.Println("itinerary usage: go run . ./input.txt ./output.txt ../airport-lookup.csv")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	airportLookupFile := os.Args[3]

	// Read input file
	itineraryText, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Input not found")
		os.Exit(1)
	}

	// Read airport lookup CSV
	airportLookup, err := readAirportLookup(airportLookupFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Process itinerary
	processedItinerary := processItinerary(string(itineraryText), airportLookup)

	// Write output file
	err = ioutil.WriteFile(outputFile, []byte(processedItinerary), 0644)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}
}

func readAirportLookup(filename string) (map[string]string, error) {
	airportLookup := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Airport lookup not found")
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	_, err = reader.Read() // Skip header row
	if err != nil {
		return nil, fmt.Errorf("Error reading airport lookup file: %v", err)
	}

	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		if len(line) < 5 {
			return nil, fmt.Errorf("Airport lookup malformed: Insufficient data in line")
		}
		airportName := line[0]
		icaoCode := line[3]
		iataCode := line[4]
		if icaoCode != "" {
			airportLookup["##"+strings.TrimSpace(icaoCode)] = airportName
		}
		if iataCode != "" {
			airportLookup["#"+strings.TrimSpace(iataCode)] = airportName
		}
	}

	return airportLookup, nil
}

// func processItinerary(itinerary string, airportLookup map[string]string) string {
// 	lines := strings.Split(itinerary, "\n")
// 	var processedLines []string
// 	for _, line := range lines {
// 		trimmedLine := strings.TrimSpace(line)
// 		if trimmedLine == "" {
// 			processedLines = append(processedLines, "")
// 			continue
// 		}

// 		parts := strings.Split(trimmedLine, "(")
// 		if len(parts) < 2 {
// 			processedLines = append(processedLines, trimmedLine)
// 			continue
// 		}

// 		text := parts[0]
// 		dateTime := strings.TrimRight(parts[1], ")")
// 		if strings.HasPrefix(dateTime, "D") {
// 			date := strings.TrimPrefix(dateTime, "D")
// 			parsedDate, err := time.Parse(time.RFC3339, date)
// 			if err != nil {
// 				processedLines = append(processedLines, trimmedLine)
// 				continue
// 			}
// 			formattedDate := parsedDate.Format("02 Jan 2006")
// 			processedLines = append(processedLines, fmt.Sprintf("%s (%s)", text, formattedDate))
// 		} else if strings.HasPrefix(dateTime, "T12") {
// 			timeWithOffset := strings.TrimPrefix(dateTime, "T12")
// 			parsedTime, err := time.Parse(time.RFC3339, timeWithOffset)
// 			if err != nil {
// 				processedLines = append(processedLines, trimmedLine)
// 				continue
// 			}
// 			formattedTime := parsedTime.Format("03:04PM (-07:00)")
// 			processedLines = append(processedLines, fmt.Sprintf("%s (%s)", text, formattedTime))
// 		} else if strings.HasPrefix(dateTime, "T24") {
// 			timeWithOffset := strings.TrimPrefix(dateTime, "T24")
// 			parsedTime, err := time.Parse(time.RFC3339, timeWithOffset)
// 			if err != nil {
// 				processedLines = append(processedLines, trimmedLine)
// 				continue
// 			}
// 			formattedTime := parsedTime.Format("15:04 (-07:00)")
// 			processedLines = append(processedLines, fmt.Sprintf("%s (%s)", text, formattedTime))
// 		} else {
// 			processedLines = append(processedLines, trimmedLine)
// 		}
// 	}
// 	return strings.Join(processedLines, "\n")
// }

func processItinerary(itinerary string, airportLookup map[string]string) string {
	lines := strings.Split(itinerary, "\n")
	var processedLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			processedLines = append(processedLines, "")
			continue
		}

		parts := strings.Split(trimmedLine, "(")
		if len(parts) < 2 {
			processedLines = append(processedLines, trimmedLine)
			continue
		}

		text := parts[0]
		dateTime := strings.TrimRight(parts[1], ")")

		if strings.HasPrefix(dateTime, "D") {
			date := strings.TrimPrefix(dateTime, "D")
			parsedDate, err := time.Parse("2006-01-02T15:04:05Z", date)
			if err != nil {
				processedLines = append(processedLines, trimmedLine)
				continue
			}
			formattedDate := parsedDate.Format("02 Jan 2006")
			processedLines = append(processedLines, fmt.Sprintf("%s (%s)", text, formattedDate))
		} else if strings.HasPrefix(dateTime, "T12") || strings.HasPrefix(dateTime, "T24") {
			timeWithOffset := strings.TrimPrefix(dateTime, "T12")
			timeWithOffset = strings.TrimPrefix(timeWithOffset, "T24")
			parsedTime, err := time.Parse("2006-01-02T15:04:05-07:00", timeWithOffset)
			if err != nil {
				processedLines = append(processedLines, trimmedLine)
				continue
			}
			formattedTime := parsedTime.Format("03:04PM (-07:00)")
			if strings.HasPrefix(dateTime, "T24") {
				formattedTime = parsedTime.Format("15:04 (-07:00)")
			}
			processedLines = append(processedLines, fmt.Sprintf("%s (%s)", text, formattedTime))
		} else {
			processedLines = append(processedLines, trimmedLine)
		}
	}
	return strings.Join(processedLines, "\n")
}
