package functions

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
	"github.com/yywing/sl/native"
)

func init() {
	LibFunctions["contains"] = native.MustNewNativeFunction("contains", Contains)
	LibFunctions["startsWith"] = native.MustNewNativeFunction("startsWith", StartsWith)
	LibFunctions["endsWith"] = native.MustNewNativeFunction("endsWith", EndsWith)
	LibFunctions["matches"] = native.MustNewNativeFunction("matches", Matches)
	LibFunctions["charAt"] = native.MustNewNativeFunction("charAt", CharAt)
	LibFunctions["indexOf"] = native.MustNewNativeFunction("indexOf", IndexOf).WithDefaultArg(int64(0))
	LibFunctions["lastIndexOf"] = native.MustNewNativeFunction("lastIndexOf", LastIndexOf).WithDefaultArg(int64(-1))
	LibFunctions["lowerAscii"] = native.MustNewNativeFunction("lowerAscii", LowerASCII)
	LibFunctions["replace"] = native.MustNewNativeFunction("replace", Replace).WithDefaultArg(int64(-1))
	LibFunctions["split"] = native.MustNewNativeFunction("split", Split).WithDefaultArg(int64(-1))
	LibFunctions["substring"] = native.MustNewNativeFunction("substring", Substring).WithDefaultArg(int64(-1))
	LibFunctions["trim"] = native.MustNewNativeFunction("trim", Trim)
	LibFunctions["upperAscii"] = native.MustNewNativeFunction("upperAscii", UpperASCII)
	// TODO
	// LibFunctions["format"] = native.MustNewNativeFunction("format", Format)
	LibFunctions["quote"] = native.MustNewNativeFunction("quote", Quote)
	LibFunctions["join"] = native.MustNewNativeFunction("join", Join).WithDefaultArg("")
	LibFunctions["reverse"] = native.MustNewNativeFunction("reverse", Reverse)
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func Matches(s, pattern string) (bool, error) {
	re, err := regexp2.Compile(pattern, regexp2.RE2)
	if err != nil {
		return false, fmt.Errorf("regexp %v compile failed: %v\n", pattern, err)
	}

	ret, err := re.MatchString(s)
	if err != nil {
		return false, fmt.Errorf("regexp %v match %v failed: %v\n", pattern, s, err)
	}

	return ret, nil
}

func CharAt(str string, ind int64) (string, error) {
	i := int(ind)
	runes := []rune(str)
	if i < 0 || i > len(runes) {
		return "", fmt.Errorf("index out of range: %d", ind)
	}
	if i == len(runes) {
		return "", nil
	}
	return string(runes[i]), nil
}

func IndexOf(str, substr string, offset int64) (int64, error) {
	if substr == "" {
		return offset, nil
	}
	off := int(offset)
	runes := []rune(str)
	subrunes := []rune(substr)
	if off < 0 {
		return -1, fmt.Errorf("index out of range: %d", off)
	}
	// If the offset exceeds the length, return -1 rather than error.
	if off >= len(runes) {
		return -1, nil
	}
	for i := off; i < len(runes)-(len(subrunes)-1); i++ {
		found := true
		for j := 0; j < len(subrunes); j++ {
			if runes[i+j] != subrunes[j] {
				found = false
				break
			}
		}
		if found {
			return int64(i), nil
		}
	}
	return -1, nil
}

func LastIndexOf(str, substr string, offset int64) (int64, error) {
	if substr == "" {
		if offset < 0 {
			return int64(len(str)), nil
		}
		return offset, nil
	}

	off := int(offset)
	runes := []rune(str)
	subrunes := []rune(substr)
	if off < 0 {
		off = len(runes) - 1
	}

	// If the offset is far greater than the length return -1
	if off >= len(runes) {
		return -1, nil
	}
	if off > len(runes)-len(subrunes) {
		off = len(runes) - len(subrunes)
	}
	for i := off; i >= 0; i-- {
		found := true
		for j := 0; j < len(subrunes); j++ {
			if runes[i+j] != subrunes[j] {
				found = false
				break
			}
		}
		if found {
			return int64(i), nil
		}
	}
	return -1, nil
}

func LowerASCII(str string) (string, error) {
	runes := []rune(str)
	for i, r := range runes {
		if r <= unicode.MaxASCII {
			r = unicode.ToLower(r)
			runes[i] = r
		}
	}
	return string(runes), nil
}

func Replace(str, old, new string, n int64) (string, error) {
	return strings.Replace(str, old, new, int(n)), nil
}

func Split(str, sep string, n int64) ([]string, error) {
	return strings.SplitN(str, sep, int(n)), nil
}

func Substring(str string, start, end int64) (string, error) {
	runes := []rune(str)
	l := len(runes)
	if int(end) < 0 {
		end = int64(l)
	}
	if start > end {
		return "", fmt.Errorf("invalid substring range. start: %d, end: %d", start, end)
	}
	if int(start) < 0 || int(start) > l {
		return "", fmt.Errorf("index out of range: %d", start)
	}
	if int(end) > l {
		return "", fmt.Errorf("index out of range: %d", end)
	}
	return string(runes[int(start):int(end)]), nil
}

func Trim(str string) (string, error) {
	return strings.TrimSpace(str), nil
}

func UpperASCII(str string) (string, error) {
	runes := []rune(str)
	for i, r := range runes {
		if r <= unicode.MaxASCII {
			r = unicode.ToUpper(r)
			runes[i] = r
		}
	}
	return string(runes), nil
}

func Reverse(str string) (string, error) {
	chars := []rune(str)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars), nil
}

func Join(strs []string, separator string) (string, error) {
	return strings.Join(strs, separator), nil
}

// quote implements a string quoting function. The string will be wrapped in
// double quotes, and all valid CEL escape sequences will be escaped to show up
// literally if printed. If the input contains any invalid UTF-8, the invalid runes
// will be replaced with utf8.RuneError.
func Quote(s string) (string, error) {
	var quotedStrBuilder strings.Builder
	for _, c := range sanitize(s) {
		switch c {
		case '\a':
			quotedStrBuilder.WriteString("\\a")
		case '\b':
			quotedStrBuilder.WriteString("\\b")
		case '\f':
			quotedStrBuilder.WriteString("\\f")
		case '\n':
			quotedStrBuilder.WriteString("\\n")
		case '\r':
			quotedStrBuilder.WriteString("\\r")
		case '\t':
			quotedStrBuilder.WriteString("\\t")
		case '\v':
			quotedStrBuilder.WriteString("\\v")
		case '\\':
			quotedStrBuilder.WriteString("\\\\")
		case '"':
			quotedStrBuilder.WriteString("\\\"")
		default:
			quotedStrBuilder.WriteRune(c)
		}
	}
	escapedStr := quotedStrBuilder.String()
	return "\"" + escapedStr + "\"", nil
}

// sanitize replaces all invalid runes in the given string with utf8.RuneError.
func sanitize(s string) string {
	var sanitizedStringBuilder strings.Builder
	for _, r := range s {
		if !utf8.ValidRune(r) {
			sanitizedStringBuilder.WriteRune(utf8.RuneError)
		} else {
			sanitizedStringBuilder.WriteRune(r)
		}
	}
	return sanitizedStringBuilder.String()
}
