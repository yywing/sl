package sl

import (
	"testing"

	"github.com/yywing/sl/ast"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  ast.ASTNode
	}{
		{
			input: "1 + 2",
			want:  ast.NewFunctionCall(ast.NewIdent(ast.Add, false), []ast.ASTNode{ast.NewIdent("1", false), ast.NewIdent("2", false)}),
		},
		{
			input: "a(1)",
			want:  ast.NewFunctionCall(ast.NewIdent("a", false), []ast.ASTNode{ast.NewIdent("1", false)}),
		},
		{
			input: "a.b(1)",
			want:  ast.NewFunctionCall(ast.NewMemberAccess(ast.NewIdent("a", false), "b", false), []ast.ASTNode{ast.NewIdent("1", false)}),
		},
	}

	for _, test := range tests {
		got, err := Parse(test.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", test.input, err)
		}
		if got.String() != test.want.String() {
			t.Errorf("Parse(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}
