// File: opa_to_lcov.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type CoverageReport struct {
	Files map[string]FileCoverage `json:"files"`
}

type FileCoverage struct {
	Covered []Range `json:"covered"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Row int `json:"row"`
}

// ConvertOPACoverageToLCOV converts OPA coverage data to LCOV format and writes it to the provided writer
func ConvertOPACoverageToLCOV(opaCoverageData []byte, writer io.Writer) error {
	// Parse OPA coverage JSON
	var coverageReport CoverageReport
	err := json.Unmarshal(opaCoverageData, &coverageReport)
	if err != nil {
		return fmt.Errorf("error parsing OPA coverage JSON: %v", err)
	}

	// Write LCOV formatted data
	for file, fileCoverage := range coverageReport.Files {
		_, err = writer.Write([]byte(fmt.Sprintf("SF:%s\n", file)))
		if err != nil {
			return fmt.Errorf("error writing to LCOV: %v", err)
		}

		totalLines := 0
		coveredLines := 0

		// Write each covered range
		for _, coveredRange := range fileCoverage.Covered {
			for line := coveredRange.Start.Row; line <= coveredRange.End.Row; line++ {
				totalLines++
				coveredLines++
				_, err = writer.Write([]byte(fmt.Sprintf("DA:%d,1\n", line))) // 1 = covered
				if err != nil {
					return fmt.Errorf("error writing to LCOV: %v", err)
				}
			}
		}

		// Write LF and LH tags
		_, err = writer.Write([]byte(fmt.Sprintf("LF:%d\n", totalLines)))   // Total number of instrumented lines
		_, err = writer.Write([]byte(fmt.Sprintf("LH:%d\n", coveredLines))) // Number of covered lines
		if err != nil {
			return fmt.Errorf("error writing to LCOV: %v", err)
		}

		_, err = writer.Write([]byte("end_of_record\n"))
		if err != nil {
			return fmt.Errorf("error writing to LCOV: %v", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: opa-to-lcov <opa-coverage.json> <output-lcov.info>")
		os.Exit(1)
	}

	opaCoverageFile := os.Args[1]
	outputLCOVFile := os.Args[2]

	// Read OPA coverage JSON file
	opaData, err := os.ReadFile(opaCoverageFile)
	if err != nil {
		fmt.Printf("Error reading OPA coverage file: %v\n", err)
		os.Exit(1)
	}

	// Create and open LCOV output file
	outputFile, err := os.Create(outputLCOVFile)
	if err != nil {
		fmt.Printf("Error creating LCOV output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Convert coverage data
	err = ConvertOPACoverageToLCOV(opaData, outputFile)
	if err != nil {
		fmt.Printf("Error converting coverage: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("LCOV file generated successfully.")
}
