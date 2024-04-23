package main

import (
	"flag"
	"fmt"
	"os"
)

// ANSI escape codes for coloring and formatting text
const (
	ansiReset  = "\033[0m" // Reset all formatting
	boldBlue   = "\033[1;34m"
	boldYellow = "\033[1;33m"
)

type Match struct {
	Index int    // sort the matches
	Value string // matched text
	Type  string
}

func main() {
	helpFlag := flag.Bool("h", false, "Display help")
	flag.Parse()

	// Check if no arguments were passed or if the help flag is set
	if *helpFlag || len(flag.Args()) == 0 {
		fmt.Println("itinerary usage:")
		fmt.Println("go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	args := flag.Args()
	if len(args) < 3 {
		// Missing required arguments
		fmt.Println("itinerary usage:")
		fmt.Println("go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	inputFilePath := args[0]
	outputFilePath := args[1]
	csvFilePath := args[2]

	// Check if output.txt file exists
	if fileExists(outputFilePath) {
		fmt.Println("Output file already exists. Exiting.")
		return
	}

	// Open the CSV file and find indices
	csvFile, err := openCSV(csvFilePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer csvFile.Close()

	header, err := readCSVHeader(csvFile)
	if err != nil {
		fmt.Println("Error reading CSV header:", err)
		return
	}

	// Inserted code snippet here
	iataIndex, icaoIndex, nameIndex, cityIndex := findColumnIndices(header)

	if iataIndex == -1 || icaoIndex == -1 || nameIndex == -1 || cityIndex == -1 {
		fmt.Println("Airport lookup malformed. Header:", header)
		return
	}

	// Open input file and process it
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Input not found:", err)
		return
	}
	defer inputFile.Close()

	output, err := processInputFile(inputFile, csvFile, iataIndex, icaoIndex, nameIndex, cityIndex) // Updated to include cityIndex
	if err != nil {
		fmt.Println("Error processing input file:", err)
		return
	}

	// Write output to file
	if err := writeOutput(outputFilePath, output); err != nil {
		fmt.Println("Error writing to output file:", err)
	}

	// Also print the output to the terminal
	fmt.Print(output) // output contains the formatted text
}

func writeOutput(filename, output string) error {
	if output == "" {
		return nil
	}
	return os.WriteFile(filename, []byte(output), 0644)
}
