package main

import (
	"fmt"
	"regexp"
	"time"
)

func processLine(line string) (string, error) {
	// regex pattern to match the format
	pattern := `([DT])(\d{2})?\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(Z|[\+\-]\d{2}:\d{2}))\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 5 {
		return "", fmt.Errorf("invalid line format")
	}

	// extract components from the matches
	prefix := matches[1]
	timeString := matches[3]
	offset := matches[4]

	// parse time
	layout := "2006-01-02T15:04Z07:00"
	timeValue, err := time.Parse(layout, timeString)
	if err != nil {
		return "", err
	}

	// format time based on the prefix
	var output string
	switch prefix {
	case "D":
		output = timeValue.Format("02 Jan 2006")
	case "T":
		if matches[2] == "12" {
			output = timeValue.Format("03:04PM")
		} else if matches[2] == "24" {
			output = timeValue.Format("15:04")
		} else {
			return "", fmt.Errorf("unknown time format")
		}
	}

	if offset == "Z" {
		offset = "(+00:00)"
	} else {
		offset = "(" + offset + ")"
	}

	// exclude offset for dates
	if prefix == "D" {
		return output, nil
	}
	return output + " " + offset, nil
}
