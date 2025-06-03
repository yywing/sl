package test

import (
	"fmt"
	"slices"
	"testing"
)

func TestConversions(t *testing.T) {

	tests := []string{
		"testdata/conversions.textproto",
	}
	skipTests := []string{
		// TODO: timestamp not support
		"conversions/int/timestamp",
		"conversions/identity/timestamp",
		// TODO: duration not support
		"conversions/identity/duration",

		// TODO: type function not support
		"conversions/type/bool_denotation",
		"conversions/type/int_denotation",
		"conversions/type/uint_denotation",
		"conversions/type/double_denotation",
		"conversions/type/null_type_denotation",
		"conversions/type/string_denotation",
		"conversions/type/bytes_denotation",
		"conversions/type/list_denotation",
		"conversions/type/map_denotation",
		"conversions/type/type",
		"conversions/type/type_denotation",
		"conversions/type/type_type",

		// feature:dyn not support
		"conversions/dyn/dyn_heterogeneous_list",
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
