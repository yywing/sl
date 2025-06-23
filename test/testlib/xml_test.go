package lib

import (
	"fmt"
	"slices"
	"testing"

	"github.com/yywing/sl/test"
)

func TestXML(t *testing.T) {
	tests := []string{
		"testdata/xml.textproto",
	}
	skipTests := []string{}

	files := test.LoadTestFile(tests)
	for _, file := range files {
		for _, section := range file.GetSection() {
			for _, testCase := range section.GetTest() {
				name := fmt.Sprintf("%s/%s/%s", file.GetName(), section.GetName(), testCase.GetName())

				if slices.Contains(skipTests, name) {
					continue
				}

				if err := test.RunTestCase(testCase); err != nil {
					t.Errorf("RunTestCase(%q) error: %v", name, err)
				}
			}
		}
	}

}
