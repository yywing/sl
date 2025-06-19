package sl

import "github.com/yywing/sl/ast"

type VariablesType map[string]ast.ValueType

type Variables map[string]ast.Value

// Get variable type definition
func (v Variables) Type() VariablesType {
	t := make(VariablesType)
	for name, value := range v {
		t[name] = value.Type()
	}
	return t
}
