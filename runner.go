package sl

import (
	"fmt"

	"github.com/yywing/sl/ast"
)

// RuntimeError 表示运行时错误
type RuntimeError struct {
	Message string
	Node    ast.ASTNode
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error: %s", e.Message)
}

// Runner 实现表达式求值
type Runner struct {
	env *Env
}

// NewRunner 创建新的求值器
func NewRunner(env *Env) *Runner {
	return &Runner{env: env}
}

// Eval 求值表达式
func (runner *Runner) Eval(node ast.ASTNode) (ast.Value, error) {
	result, err := node.Accept(runner)
	if err != nil {
		return nil, err
	}
	if value, ok := result.(ast.Value); ok {
		return value, nil
	}
	return nil, fmt.Errorf("internal error: Runner returned non-value")
}

func (runner *Runner) VisitLiteral(node *ast.LiteralNode) (interface{}, error) {
	return node.Value, nil
}

func (runner *Runner) VisitIdent(node *ast.IdentNode) (interface{}, error) {
	if value, exists := runner.env.GetVariable(node.Name); exists {
		return value, nil
	}

	return nil, &RuntimeError{
		Message: fmt.Sprintf("undefined identifier: %s", node.Name),
		Node:    node,
	}
}

func (runner *Runner) VisitMemberAccess(node *ast.MemberAccessNode) (interface{}, error) {
	object, err := runner.Eval(node.Object)
	if err != nil {
		return nil, err
	}

	switch obj := object.(type) {
	case *ast.MapValue:
		if value, exists := obj.Get(ast.NewStringValue(node.Member)); exists {
			return value, nil
		}
		if node.Optional {
			return ast.NewNullValue(), nil
		}
		return nil, &RuntimeError{
			Message: fmt.Sprintf("map does not have member: %s", node.Member),
			Node:    node,
		}
	default:
		// 对于其他类型，简化处理
		if node.Optional {
			return ast.NewNullValue(), nil
		}
		return nil, &RuntimeError{
			Message: fmt.Sprintf("cannot access member %s on type %T", node.Member, object),
			Node:    node,
		}
	}
}

func (runner *Runner) VisitFunctionCall(node *ast.FunctionCallNode) (interface{}, error) {
	var fnName string
	var args []ast.ASTNode
	switch fn := node.Function.(type) {
	case *ast.IdentNode:
		fnName = fn.Name
		args = node.Args
	case *ast.MemberAccessNode:
		fnName = fn.Member
		args = append([]ast.ASTNode{fn.Object}, node.Args...)
	default:
		return nil, &CheckError{
			Message: fmt.Sprintf("function call must be an identifier or member access, got %s", node.Function.String()),
			Node:    node,
		}
	}

	fn, exists := runner.env.GetFunction(fnName)
	if !exists {
		return nil, &RuntimeError{
			Message: fmt.Sprintf("function %s not found", fnName),
			Node:    node,
		}
	}

	// 求值参数
	argValues := make([]ast.Value, len(args))
	for i, arg := range args {
		argValue, err := runner.Eval(arg)
		if err != nil {
			// 特殊处理一下 or
			if fn.Name() == ast.LogicalOr {
				argValue = ast.NewBoolValue(false)
			} else {
				return nil, err
			}
		}
		argValues[i] = argValue
	}

	// 调用函数
	result, err := fn.Call(argValues)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (runner *Runner) VisitIndex(node *ast.IndexNode) (interface{}, error) {
	object, err := runner.Eval(node.Object)
	if err != nil {
		return nil, err
	}

	index, err := runner.Eval(node.Index)
	if err != nil {
		return nil, err
	}

	switch obj := object.(type) {
	case *ast.ListValue:
		var idx int
		switch i := index.(type) {
		case *ast.IntValue:
			idx = int(i.IntValue)
		case *ast.UintValue:
			idx = int(i.UintValue)
		default:
			return nil, &RuntimeError{
				Message: fmt.Sprintf("list index must be integer, got %T", index),
				Node:    node,
			}
		}

		values := obj.ListValue
		if idx < 0 || idx >= len(values) {
			if node.Optional {
				return ast.NewNullValue(), nil
			}
			return nil, &RuntimeError{
				Message: fmt.Sprintf("list index out of range: %d", idx),
				Node:    node,
			}
		}
		return values[idx], nil

	case *ast.MapValue:
		if value, exists := obj.Get(index); exists {
			return value, nil
		}
		if node.Optional {
			return ast.NewNullValue(), nil
		}
		return nil, &RuntimeError{
			Message: fmt.Sprintf("map does not have key: %s", index.String()),
			Node:    node,
		}

	default:
		return nil, &RuntimeError{
			Message: fmt.Sprintf("cannot index type %T", object),
			Node:    node,
		}
	}
}

func (runner *Runner) VisitConditional(node *ast.ConditionalNode) (interface{}, error) {
	condition, err := runner.Eval(node.Condition)
	if err != nil {
		return nil, err
	}

	result, ok := condition.(*ast.BoolValue)
	if !ok {
		return nil, &RuntimeError{
			Message: fmt.Sprintf("condition must be boolean, got %T", condition),
			Node:    node,
		}
	}

	if result.BoolValue {
		return runner.Eval(node.TrueExpr)
	} else {
		return runner.Eval(node.FalseExpr)
	}
}

func (runner *Runner) VisitList(node *ast.ListNode) (interface{}, error) {
	values := make([]ast.Value, len(node.Elements))
	var elementType ast.ValueType = ast.AnyType

	for i, elem := range node.Elements {
		value, err := runner.Eval(elem)
		if err != nil {
			return nil, err
		}
		values[i] = value

		if i == 0 {
			elementType = value.Type()
		}

		if !value.Type().Equals(elementType) {
			elementType = ast.AnyType
		}
	}

	return ast.NewListValue(values, elementType), nil
}

func (runner *Runner) VisitMap(node *ast.MapNode) (interface{}, error) {
	values := make(map[ast.Value]ast.Value)
	var keyType, valueType ast.ValueType = ast.AnyType, ast.AnyType
	keys := []ast.Value{}

	for i, entry := range node.Entries {
		key, err := runner.Eval(entry.Key)
		if err != nil {
			return nil, err
		}

		value, err := runner.Eval(entry.Value)
		if err != nil {
			return nil, err
		}

		for _, j := range keys {
			if j.Equal(key) {
				return nil, &RuntimeError{
					Message: fmt.Sprintf("map has repeated key: %s", key.String()),
					Node:    node,
				}
			}
		}

		keys = append(keys, key)
		values[key] = value

		if i == 0 {
			keyType = key.Type()
			valueType = value.Type()
		}

		if !key.Type().Equals(keyType) {
			keyType = ast.AnyType
		}
		if !value.Type().Equals(valueType) {
			valueType = ast.AnyType
		}
	}

	return ast.NewMapValue(values, keyType, valueType), nil
}

// TODO:
func (runner *Runner) VisitStruct(node *ast.StructNode) (interface{}, error) {
	// // 创建一个映射来表示结构体实例
	// values := make(map[string]ast.Value)

	// // 求值所有字段
	// for _, field := range node.Fields {
	// 	value, err := runner.Eval(field.Value)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	values[field.Name] = value
	// }

	// // 返回一个映射值来表示结构体实例
	// // 在更完整的实现中，可能需要专门的StructValue类型
	// return ast.NewMapValue(values, ast.StringType, ast.StringType), nil
	return nil, fmt.Errorf("not implemented")
}
