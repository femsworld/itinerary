package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

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

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
