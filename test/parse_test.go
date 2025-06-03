package test

import (
	"fmt"
	"slices"
	"testing"
)

func TestParse(t *testing.T) {

	tests := []string{
		"testdata/parse.textproto",
	}
	skipTests := []string{
		// TODO: struct not support
		"parse/nest/message_literal",
		"parse/repeat/select",
		"parse/repeat/message_literal",
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
