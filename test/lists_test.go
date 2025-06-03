package test

import (
	"fmt"
	"slices"
	"testing"
)

func TestLists(t *testing.T) {

	tests := []string{
		"testdata/lists.textproto",
	}
	skipTests := []string{
		// feature: dyn not support
		"lists/index/zero_based_double",
		"lists/index/zero_based_uint",
		"lists/index/index_out_of_bounds_or_true",
		"lists/index/index_out_of_bounds_and_false",
		"lists/index/bad_index_type_or_true",
		"lists/index/bad_index_type_and_false",
		"lists/in/double_in_ints",
		"lists/in/uint_in_ints",
		"lists/in/int_in_doubles",
		"lists/in/uint_in_doubles",
		"lists/in/int_in_uints",
		"lists/in/double_in_uints",
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
