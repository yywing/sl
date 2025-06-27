package types_test

import (
	"net/http"
	"testing"

	"github.com/yywing/sl"
	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/lib/types"
)

type TestCase struct {
	variables sl.Variables
	expr      string
	want      ast.Value
	wantErr   bool
}

func RunTestCase(t *testing.T, testCase TestCase) {
	env := sl.NewStdEnv()

	ast, err := sl.Parse(testCase.expr)
	if err != nil {
		t.Fatal(err)
	}

	program := sl.NewProgram(ast, testCase.variables.Type())

	_, err = env.Check(program)

	if testCase.wantErr {
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		return
	}

	if err != nil {
		t.Fatal(err)
	}

	result, err := env.Run(program, testCase.variables)
	if err != nil {
		t.Fatal(err)
	}

	if !testCase.want.Equal(result) {
		t.Fatalf("want %v, got %v", testCase.want, result)
	}
}

func TestHTTPRequestType(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	value, err := types.NewHTTPRequestFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	var testCases = []TestCase{
		{
			variables: sl.Variables{
				"request": value,
			},
			expr: "request.url.scheme",
			want: ast.NewStringValue("https"),
		},
		{
			variables: sl.Variables{
				"request": value,
			},
			expr: "string(request.url)",
			want: ast.NewStringValue("https://example.com"),
		},
		{
			variables: sl.Variables{
				"request": value,
			},
			expr:    "request.xxx",
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		RunTestCase(t, testCase)
	}
}
