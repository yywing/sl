package test

import (
	"fmt"
	"slices"
	"testing"
)

func TestTimestamps(t *testing.T) {

	tests := []string{
		"testdata/timestamps.textproto",
	}
	skipTests := []string{
		// feature: time type name not samw
		"timestamps/timestamp_conversions/toType_timestamp",
		"timestamps/duration_conversions/toType_duration",
		"timestamps/duration_converters/get_milliseconds",
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
