// File: opa_to_lcov_test.go
package main

import (
	"bytes"
	"testing"
)

func TestConvertOPACoverageToLCOV(t *testing.T) {
	// Sample OPA coverage JSON input
	opaCoverageJSON := `
	{
	  "files": {
		"main_test.rego": {
		  "covered": [
			{
			  "start": {
				"row": 7
			  },
			  "end": {
				"row": 8
			  }
			},
			{
			  "start": {
				"row": 11
			  },
			  "end": {
				"row": 12
			  }
			}
		  ],
		  "covered_lines": 4,
		  "coverage": 100
		}
	  },
	  "covered_lines": 4,
	  "not_covered_lines": 0,
	  "coverage": 100
	}
	`

	// Expected LCOV output
	expectedLCOV := `SF:main_test.rego
DA:7,1
DA:8,1
DA:11,1
DA:12,1
LF:4
LH:4
end_of_record
`

	var output bytes.Buffer

	err := ConvertOPACoverageToLCOV([]byte(opaCoverageJSON), &output)
	if err != nil {
		t.Fatalf("Error converting coverage: %v", err)
	}

	if output.String() != expectedLCOV {
		t.Errorf("LCOV output does not match expected output.\nExpected:\n%s\nGot:\n%s", expectedLCOV, output.String())
	}
}

func TestConvertOPACoverageToLCOV_EmptyCoverage(t *testing.T) {
	// OPA coverage JSON with no coverage data
	opaCoverageJSON := `
	{
	  "files": {
		"main_test.rego": {
		  "covered": [],
		  "covered_lines": 0,
		  "coverage": 0
		}
	  },
	  "covered_lines": 0,
	  "not_covered_lines": 10,
	  "coverage": 0
	}
	`

	// Expected LCOV output
	expectedLCOV := `SF:main_test.rego
LF:0
LH:0
end_of_record
`

	var output bytes.Buffer

	err := ConvertOPACoverageToLCOV([]byte(opaCoverageJSON), &output)
	if err != nil {
		t.Fatalf("Error converting coverage: %v", err)
	}

	if output.String() != expectedLCOV {
		t.Errorf("LCOV output does not match expected output for empty coverage.\nExpected:\n%s\nGot:\n%s", expectedLCOV, output.String())
	}
}

func TestConvertOPACoverageToLCOV_InvalidJSON(t *testing.T) {
	// Invalid JSON input
	opaCoverageJSON := `invalid json`

	var output bytes.Buffer

	err := ConvertOPACoverageToLCOV([]byte(opaCoverageJSON), &output)
	if err == nil {
		t.Fatalf("Expected error for invalid JSON input, got none")
	}
}

func TestConvertOPACoverageToLCOV_MultipleFiles(t *testing.T) {
	// OPA coverage JSON with multiple files
	opaCoverageJSON := `
	{
	  "files": {
		"file1.rego": {
		  "covered": [
			{
			  "start": {
				"row": 1
			  },
			  "end": {
				"row": 2
			  }
			}
		  ],
		  "covered_lines": 2,
		  "coverage": 100
		},
		"file2.rego": {
		  "covered": [
			{
			  "start": {
				"row": 3
			  },
			  "end": {
				"row": 4
			  }
			}
		  ],
		  "covered_lines": 2,
		  "coverage": 100
		}
	  },
	  "covered_lines": 4,
	  "not_covered_lines": 0,
	  "coverage": 100
	}
	`

	// Expected LCOV output
	expectedLCOV := `SF:file1.rego
DA:1,1
DA:2,1
LF:2
LH:2
end_of_record
SF:file2.rego
DA:3,1
DA:4,1
LF:2
LH:2
end_of_record
`

	var output bytes.Buffer

	err := ConvertOPACoverageToLCOV([]byte(opaCoverageJSON), &output)
	if err != nil {
		t.Fatalf("Error converting coverage: %v", err)
	}

	if output.String() != expectedLCOV {
		t.Errorf("LCOV output does not match expected output for multiple files.\nExpected:\n%s\nGot:\n%s", expectedLCOV, output.String())
	}
}
