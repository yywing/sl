package test

import (
	"fmt"
	"slices"
	"testing"
)

func TestLogic(t *testing.T) {

	tests := []string{
		"testdata/logic.textproto",
	}
	skipTests := []string{
		// feature: not return error
		"logic/OR/error_right",
		"logic/OR/error_left",

		// wrong type not support
		"logic/AND/short_circuit_type_left",
		"logic/AND/short_circuit_type_right",
		"logic/AND/short_circuit_error_left",
		"logic/AND/short_circuit_error_right",
		"logic/OR/short_circuit_type_left",
		"logic/OR/short_circuit_type_right",
	}

	files := LoadTestFile(tests)
	for _, file := range files {
		for _, section := range file.GetSection() {
			for _, testCase := range section.GetTest() {
				name := fmt.Sprintf("%s/%s/%s", file.GetName(), section.GetName(), testCase.GetName())

				if slices.Contains(skipTests, name) {
					continue
				}

				if err := RunTestCase(testCase); err != nil {
					t.Errorf("RunTestCase(%q) error: %v", name, err)
				}
			}
		}
	}

}
