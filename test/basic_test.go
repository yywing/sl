package test

import (
	"fmt"
	"slices"
	"testing"
)

func TestBasic(t *testing.T) {

	tests := []string{
		"testdata/basic.textproto",
	}
	skipTests := []string{}

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
