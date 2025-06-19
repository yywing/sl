package sl

import (
	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/lib/functions"
)

// Env 表示执行环境，包含变量和函数
type Env struct {
	functions map[string]ast.Function // 函数映射
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

// Functions 返回所有函数名
func (e *Env) Functions() []string {
	var names []string
	for name := range e.functions {
		names = append(names, name)
	}
	return names
}

func (e *Env) Check(p *Program) (ast.ValueType, error) {
	checker := NewChecker(e, p)
	return checker.Check()
}

func (e *Env) Run(p *Program, variables Variables) (ast.Value, error) {
	if err := p.CheckVariables(variables); err != nil {
		return nil, err
	}

	runner := NewRunner(e, p, variables)
	return runner.Eval()
}

func newEnv() *Env {
	return &Env{
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

	for name, fn := range functions.LibFunctions {
		env.SetFunction(name, fn)
	}

	return env
}
