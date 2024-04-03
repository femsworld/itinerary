package main

import (
	"fmt"
	"strings"
	"time"
)

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
			formattedTime, err := parseTime(dateTime)
			if err != nil {
				processedLines = append(processedLines, trimmedLine)
				continue
			}
			processedLines = append(processedLines, fmt.Sprintf("%s (%s)", text, formattedTime))
		} else {
			processedLines = append(processedLines, trimmedLine)
		}
	}
	return strings.Join(processedLines, "\n")
}

func parseTime(dateTime string) (string, error) {
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04Z",
		"2006-01-02T15:04-07:00",
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, dateTime)
		if err == nil {
			if strings.Contains(format, "Z") {
				return parsedTime.Format("03:04PM (-07:00)"), nil
			}
			return parsedTime.Format("03:04PM"), nil
		}
	}

	return "", fmt.Errorf("could not parse time: %s", dateTime)
}
