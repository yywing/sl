package test

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestStringExt(t *testing.T) {

	tests := []string{
		"testdata/string_ext.textproto",
	}
	skipTests := []string{
		// TODO: format not support
		"string_ext/format/no-op",
		"string_ext/format/mid-string substitution",
		"string_ext/format/percent escaping",
		"string_ext/format/substution inside escaped percent signs",
		"string_ext/format/substitution with one escaped percent sign on the right",
		"string_ext/format/substitution with one escaped percent sign on the left",
		"string_ext/format/multiple substitutions",
		"string_ext/format/percent sign escape sequence support",
		"string_ext/format/fixed point formatting clause",
		"string_ext/format/binary formatting clause",
		"string_ext/format/uint support for binary formatting",
		"string_ext/format/bool support for binary formatting",
		"string_ext/format/octal formatting clause",
		"string_ext/format/uint support for octal formatting clause",
		"string_ext/format/lowercase hexadecimal formatting clause",
		"string_ext/format/uppercase hexadecimal formatting clause",
		"string_ext/format/unsigned support for hexadecimal formatting clause",
		"string_ext/format/string support with hexadecimal formatting clause",
		"string_ext/format/string support with uppercase hexadecimal formatting clause",
		"string_ext/format/byte support with hexadecimal formatting clause",
		"string_ext/format/byte support with uppercase hexadecimal formatting clause",
		"string_ext/format/scientific notation formatting clause",
		"string_ext/format/default precision for fixed-point clause",
		"string_ext/format/default precision for scientific notation",
		"string_ext/format/unicode output for scientific notation",
		"string_ext/format/NaN support for fixed-point",
		"string_ext/format/positive infinity support for fixed-point",
		"string_ext/format/negative infinity support for fixed-point",
		"string_ext/format/uint support for decimal clause",
		"string_ext/format/null support for string",
		"string_ext/format/int support for string",
		"string_ext/format/bytes support for string",
		"string_ext/format/type() support for string",
		"string_ext/format/timestamp support for string",
		"string_ext/format/duration support for string",
		"string_ext/format/list support for string",
		"string_ext/format/map support for string",
		"string_ext/format/map support (all key types)",
		"string_ext/format/boolean support for %s",
		"string_ext/format/dyntype support for string formatting clause",
		"string_ext/format/dyntype support for numbers with string formatting clause",
		"string_ext/format/dyntype support for integer formatting clause",
		"string_ext/format/dyntype support for integer formatting clause (unsigned)",
		"string_ext/format/dyntype support for hex formatting clause",
		"string_ext/format/dyntype support for hex formatting clause (uppercase)",
		"string_ext/format/dyntype support for unsigned hex formatting clause",
		"string_ext/format/dyntype support for fixed-point formatting clause",
		"string_ext/format/dyntype support for scientific notation",
		"string_ext/format/dyntype NaN/infinity support for fixed-point",
		"string_ext/format/dyntype support for timestamp",
		"string_ext/format/dyntype support for duration",
		"string_ext/format/dyntype support for lists",
		"string_ext/format/dyntype support for maps",
		"string_ext/format/message field support",
		"string_ext/format/string substitution in a string variable",
		"string_ext/format/multiple substitutions in a string variable",
		"string_ext/format/substution inside escaped percent signs in a string variable",
		"string_ext/format/fixed point formatting clause in a string variable",
		"string_ext/format/binary formatting clause in a string variable",
		"string_ext/format/scientific notation formatting clause in a string variable",
		"string_ext/format/default precision for fixed-point clause in a string variable",

		// feature: follow cel-go
		"string_ext/value_errors/indexof_out_of_range",
		"string_ext/value_errors/lastindexof_out_of_range",

		// feature: lastIndexOf support ngeative index
		"string_ext/value_errors/lastindexof_negative_index",

		// feature: strings.quote not support, use quote instead
		"string_ext/quote/multiline",
		"string_ext/quote/escaped",
		"string_ext/quote/backspace",
		"string_ext/quote/form_feed",
		"string_ext/quote/carriage_return",
		"string_ext/quote/horizontal_tab",
		"string_ext/quote/vertical_tab",
		"string_ext/quote/double_slash",
		"string_ext/quote/two_escape_sequences",
		"string_ext/quote/verbatim",
		"string_ext/quote/ends_with",
		"string_ext/quote/starts_with",
		"string_ext/quote/printable_unicode",
		"string_ext/quote/mid_string_quote",
		"string_ext/quote/single_quote_with_double_quote",
		"string_ext/quote/size_unicode_char",
		"string_ext/quote/size_unicode_string",
		"string_ext/quote/unicode",
		"string_ext/quote/unicode_code_points",
		"string_ext/quote/unicode_2",
		"string_ext/quote/empty_quote",
		"string_ext/format_errors/multiline",
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

func TestStringExtQuote(t *testing.T) {

	tests := []string{
		"testdata/string_ext.textproto",
	}

	files := LoadTestFile(tests)
	for _, file := range files {
		for _, section := range file.GetSection() {
			for _, testCase := range section.GetTest() {
				name := fmt.Sprintf("%s/%s/%s", file.GetName(), section.GetName(), testCase.GetName())

				if !strings.HasPrefix(name, "string_ext/quote/") && name != "string_ext/format_errors/multiline" {
					continue
				}

				testCase.Expr = strings.ReplaceAll(testCase.Expr, "strings.quote", "quote")

				if err := RunTestCase(testCase); err != nil {
					t.Errorf("RunTestCase(%q) error: %v", name, err)
				}
			}
		}
	}

}
