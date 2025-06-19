package sl

import (
	"fmt"

	"github.com/yywing/sl/ast"
)

type Program struct {
	ast.ASTNode
	variablesType VariablesType
}

func NewProgram(node ast.ASTNode, variablesType VariablesType) *Program {
	return &Program{
		ASTNode:       node,
		variablesType: variablesType,
	}
}

// SetVariable sets a variable
func (e *Program) SetVariable(name string, value ast.ValueType) {
	if e.variablesType == nil {
		e.variablesType = make(VariablesType)
	}
	e.variablesType[name] = value
}

// GetVariable gets a variable
func (e *Program) GetVariable(name string) (ast.ValueType, bool) {
	if value, exists := e.variablesType[name]; exists {
		return value, true
	}
	return nil, false
}

// Variables returns all variable names
func (e *Program) Variables() []string {
	var names []string
	for name := range e.variablesType {
		names = append(names, name)
	}
	return names
}

func (e *Program) CheckVariables(variables Variables) error {
	for k, v := range e.variablesType {
		value, exists := variables[k]
		if !exists {
			return fmt.Errorf("variable %s is not defined", k)
		}
		if !ast.TypeEquals(v, value.Type()) {
			return fmt.Errorf("variable %s is not compatible with %s", k, v)
		}
	}
	return nil
}
