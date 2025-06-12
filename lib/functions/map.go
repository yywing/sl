package functions

import (
	"fmt"

	"github.com/yywing/sl/ast"
)

func init() {
	LibFunctions["has"] = HasFunction
	LibFunctions["get"] = GetFunction
}

var (
	paramA  = ast.NewValueTypeParamType("A")
	paramB  = ast.NewValueTypeParamType("B")
	mapOfAB = ast.NewMapType(paramA, paramB)

	HasFunction = ast.NewBaseFunction("has", []ast.FunctionType{
		*ast.NewFunctionType("has", []ast.ValueType{mapOfAB, paramA}, ast.BoolType),
	}, func(args []ast.Value) (ast.Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("has expects 2 arguments, got %d", len(args))
		}

		mapValue, ok := args[0].(*ast.MapValue)
		if !ok {
			return nil, fmt.Errorf("has expects map argument, got %s", args[0].Type())
		}

		_, exists := mapValue.Get(args[1])
		return ast.NewBoolValue(exists), nil
	})

	GetFunction = ast.NewBaseFunction("get", []ast.FunctionType{
		*ast.NewFunctionType("get", []ast.ValueType{mapOfAB, paramA}, paramB),
		*ast.NewFunctionType("get", []ast.ValueType{mapOfAB, paramA, paramB}, paramB),
	}, func(args []ast.Value) (ast.Value, error) {
		switch len(args) {
		case 2:
			mapValue, ok := args[0].(*ast.MapValue)
			if !ok {
				return nil, fmt.Errorf("get expects map argument, got %s", args[0].Type())
			}

			value, exists := mapValue.Get(args[1])
			if !exists {
				return nil, fmt.Errorf("no such key %s", args[1])
			} else {
				return value, nil
			}

		case 3:
			mapValue, ok := args[0].(*ast.MapValue)
			if !ok {
				return nil, fmt.Errorf("get expects map argument, got %s", args[0].Type())
			}
			_, exists := mapValue.Get(args[1])
			if !exists {
				return args[2], nil
			}
			return args[1], nil
		default:
			return nil, fmt.Errorf("get expects 2 or 3 arguments, got %d", len(args))
		}
	})
)
