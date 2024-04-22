package main

import (
	"flag"
	"fmt"
	"os"
)

type Match struct {
	Index int    // sort the matches
	Value string // matched text
	Type  string
}

func main() {
	helpFlag := flag.Bool("h", false, "Display help")
	flag.Parse()

	// Check if no arguments were passed, or the help flag is set
	if *helpFlag || len(flag.Args()) == 0 {
		fmt.Println("itinerary usage:")
		fmt.Println("go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// Get the expected file paths from command line arguments
	args := flag.Args()
	if len(args) < 3 {
		// fmt.Println("Missing required arguments.")
		fmt.Println("itinerary usage:")
		fmt.Println("go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	inputFilePath := "./input.txt"
	outputFilePath := "./output.txt"
	csvFilePath := "./airport-lookup.csv"

	// Check if output.txt file exists
	if fileExists(outputFilePath) {
		fmt.Println("Output file already exists. Exiting.")
		return
	}

	csvFile, err := openCSV(csvFilePath)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer csvFile.Close()

	header, err := readCSVHeader(csvFile)
	if err != nil {
		fmt.Println("Airport lookup malformed.", err)
		return
	}

	iataIndex, icaoIndex, nameIndex := findColumnIndices(header)
	if iataIndex == -1 || icaoIndex == -1 || nameIndex == -1 {
		fmt.Println("Airport lookup malformed.")
		return
	}

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Input not found.")
		return
	}
	defer inputFile.Close()

	output, err := processInputFile(inputFile, csvFile, iataIndex, icaoIndex, nameIndex)
	if err != nil {
		fmt.Println("error:")
		return
	}

	if err := writeOutput(outputFilePath, output); err != nil {
		fmt.Println("error:", err)
	}
}

// write the output string to a file.
func writeOutput(filename, output string) error {
	if output == "" {
		return nil
	}
	return os.WriteFile(filename, []byte(output), 0644)
}
