package sl

import "github.com/yywing/sl/ast"

type VariablesType map[string]ast.ValueType

type Variables map[string]ast.Value

// 获取变量类型定义
func (v Variables) Type() VariablesType {
	t := make(VariablesType)
	for name, value := range v {
		t[name] = value.Type()
	}
	return t
}
