# Itinerary-prettifier

Itinerary-prettifier is a command line tool which can help prettify flight itineraries.

This tool reads a text-based itinerary from a file as input, processes the text to make it user friendly, and writes the result to a new file (output.txt).
Program uses a csv lookup file, provided with the program, to change IATA and ICAO codes to airport names.

- **IATA code:** A single # followed by three letters. For example, #LAX represents "LAX" which is the IATA code for Los Angeles International Airport.
- **ICAO code:** Double # followed by four letters. For example, ##EGLL represents "EGLL" which is the ICAO code for London Heathrow Airport.

Program also prettifies dates and times presented in ISO 8601 standard:
- **Dates:** `D(2007-04-05T12:30−02:00)`. Dates are displayed in the output as DD-Mmm-YYYY.
- **12 Hour times:** `T12(2007-04-05T12:30−02:00)`. Displayed as "12:30PM (-02:00)".
- **24 Hour times:** `T24(2007-04-05T12:30−02:00)`. Displayed as "12:30 (-02:00)".

## Usage

The program can be launched from the command line with three arguments:
- Path to the input file
- Path to the output file
- Path to the airport lookup file

```shell
$ go run . ./input.txt ./output.txt ./airport-lookup.csv

$ go run . -h

### Output

itinerary usage:
go run . ./input.txt ./output.txt ./airport-lookup.csv

