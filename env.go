package sl

import (
	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/lib"
)

// Env 表示执行环境，包含变量和函数
type Env struct {
	variables map[string]ast.Value    // 变量映射
	functions map[string]ast.Function // 函数映射
}

// SetVariable 设置变量
func (e *Env) SetVariable(name string, value ast.Value) {
	e.variables[name] = value
}

// GetVariable 获取变量
func (e *Env) GetVariable(name string) (ast.Value, bool) {
	if value, exists := e.variables[name]; exists {
		return value, true
	}
	return nil, false
}

// SetFunction 设置函数
func (e *Env) SetFunction(name string, fn ast.Function) {
	e.functions[name] = fn
}

// GetFunction 获取函数
func (e *Env) GetFunction(name string) (ast.Function, bool) {
	if fn, exists := e.functions[name]; exists {
		return fn, true
	}
	return nil, false
}

// Variables 返回所有变量名
func (e *Env) Variables() []string {
	var names []string
	for name := range e.variables {
		names = append(names, name)
	}
	return names
}

// Functions 返回所有函数名
func (e *Env) Functions() []string {
	var names []string
	for name := range e.functions {
		names = append(names, name)
	}
	return names
}

func newEnv() *Env {
	return &Env{
		variables: make(map[string]ast.Value),
		functions: make(map[string]ast.Function),
	}
}

func NewBuiltinEnv() *Env {
	env := newEnv()

	for name, fn := range ast.BuiltinFunctions {
		env.SetFunction(name, fn)
	}

	return env
}

func NewStdEnv() *Env {
	env := NewBuiltinEnv()

	for name, fn := range lib.LibFunctions {
		env.SetFunction(name, fn)
	}

	return env
}
