package functions

import (
	"fmt"

	"github.com/yywing/sl/ast"
)

func init() {
	LibFunctions[FunctionHas] = HasFunction
	LibFunctions[FunctionGet] = GetFunction
}

const (
	FunctionHas = "has"
	FunctionGet = "get"
)

var (
	paramA  = ast.NewValueTypeParamType("A")
	paramB  = ast.NewValueTypeParamType("B")
	mapOfAB = ast.NewMapType(paramA, paramB)

	HasFunction = ast.NewBaseFunction(
		FunctionHas,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(FunctionHas, []ast.ValueType{mapOfAB, paramA}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					_, exists := args[0].(*ast.MapValue).Get(args[1])
					return ast.NewBoolValue(exists), nil
				},
			},
		},
	)

	GetFunction = ast.NewBaseFunction(
		FunctionGet,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(FunctionGet, []ast.ValueType{mapOfAB, paramA}, paramB),
				Call: func(args []ast.Value) (ast.Value, error) {
					value, exists := args[0].(*ast.MapValue).Get(args[1])
					if !exists {
						return nil, fmt.Errorf("no such key %s", args[1])
					}
					return value, nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGet, []ast.ValueType{mapOfAB, paramA, paramB}, paramB),
				Call: func(args []ast.Value) (ast.Value, error) {
					_, exists := args[0].(*ast.MapValue).Get(args[1])
					if !exists {
						return args[2], nil
					}
					return args[1], nil
				},
			},
		},
	)
)
