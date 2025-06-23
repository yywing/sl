package functions

import (
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/native"
)

func init() {
	LibFunctions[FunctionJSONPath] = ast.NewBaseFunction(
		FunctionJSONPath,
		native.MustNewNativeFunction(FunctionJSONPath, JSONPath).Definitions(),
	)
}

const (
	FunctionJSONPath = "jsonPath"
)

func JSONPath(json, path string) ([]any, error) {
	obj, err := oj.ParseString(json)
	if err != nil {
		return nil, err
	}

	x, err := jp.ParseString(path)
	if err != nil {
		return nil, err
	}
	return x.Get(obj), nil
}
