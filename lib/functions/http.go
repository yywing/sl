package functions

import (
	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/lib/types"
)

func init() {
	ast.StringFunction.Combine(URLStringFunction)
}

var (
	URLStringFunction = ast.NewBaseFunction(
		ast.String,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.String, []ast.ValueType{types.URLType}, ast.StringType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewStringValue(args[0].(*types.URL).URL), nil
				},
			},
		},
	)
)
