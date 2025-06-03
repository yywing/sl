package test

import (
	"fmt"
	"slices"
	"testing"
)

func TestFields(t *testing.T) {

	tests := []string{
		"testdata/fields.textproto",
	}
	skipTests := []string{
		// TODO：namespace not supported
		"fields/qualified_identifier_resolution/qualified_identifier_resolution_unchecked",
		"fields/qualified_identifier_resolution/ident_with_longest_prefix_check",

		// feature: support more key type
		"fields/qualified_identifier_resolution/map_key_float",
		"fields/qualified_identifier_resolution/map_key_null",

		// feature: unsupported mixed key types
		"fields/map_fields/map_key_mixed_type",
		"fields/map_fields/map_key_mixed_numbers_double_key",
		"fields/map_fields/map_key_mixed_numbers_uint_key",
		"fields/map_fields/map_key_mixed_numbers_int_key",
		"fields/in/mixed_numbers_and_keys_present",
		"fields/in/mixed_numbers_and_keys_absent",

		// feature: dyn not supported
		"fields/map_fields/map_no_such_key_or_true",
		"fields/map_fields/map_no_such_key_and_false",
		"fields/map_fields/map_bad_key_type_or_true",
		"fields/map_fields/map_bad_key_type_and_false",
		"fields/map_fields/map_field_select_no_such_key_or_true",
		"fields/map_fields/map_field_select_no_such_key_and_false",

		// feature: map has not supported, use new has instead
		"fields/map_has/has",
		"fields/map_has/has_not",
		"fields/map_has/has_empty",
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
