package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// Parse command-line arguments
	if len(os.Args) != 4 || os.Args[1] == "-h" {
		fmt.Println("itinerary usage: go run . ../input.txt ../output.txt ../airport-lookup.csv")
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
	outputExists := fileExists(outputFile)
	if outputExists {
		fmt.Println("Output file already exists. Exiting.")
		os.Exit(1)
	}

	err = ioutil.WriteFile(outputFile, []byte(processedItinerary), 0644)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}
}
