package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// open the CSV file and return a file handle.
func openCSV(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("airport lookup not found")
	}
	return file, nil
}

// read the header row from the CSV and check for malformed data.
func readCSVHeader(csvFile *os.File) ([]string, error) {
	csvReader := csv.NewReader(csvFile)
	csvReader.TrimLeadingSpace = true
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// check for missing or blank columns.
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		for _, column := range row {
			if column == "" {
				return nil, err
			}
		}
	}

	return header, nil
}

// find indices of the needed columns.
func findColumnIndices(header []string) (int, int, int) {
	iataIndex, icaoIndex, nameIndex := -1, -1, -1
	for i, columnName := range header {
		switch columnName {
		case "iata_code":
			iataIndex = i
		case "icao_code":
			icaoIndex = i
		case "name":
			nameIndex = i
		}
	}
	return iataIndex, icaoIndex, nameIndex
}

// look up a code in the CSV and return the name as string.

func lookupCode(code string, csvFile *os.File, iataIndex, icaoIndex, nameIndex int) string {
	lookupCode := strings.TrimPrefix(strings.TrimPrefix(code, "#"), "#")
	csvFile.Seek(0, io.SeekStart)
	csvReader := csv.NewReader(csvFile)
	csvReader.TrimLeadingSpace = true
	csvReader.Read()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error reading a record:", err)
			break
		}
		if record[iataIndex] == lookupCode || record[icaoIndex] == lookupCode {
			return record[nameIndex]
		}
	}
	return code // return original code if not found
}
