package main

import (
	"fmt"

	"github.com/yywing/sl"
	"github.com/yywing/sl/ast"
)

func main() {
	env := sl.NewStdEnv()
	node, err := sl.Parse("a + 2")
	if err != nil {
		panic(err)
	}

	program := sl.NewProgram(node, map[string]ast.ValueType{
		"a": ast.IntType,
	})

	// check
	_, err = env.Check(program)
	if err != nil {
		panic(err)
	}

	// run
	result, err := env.Run(program, map[string]ast.Value{
		"a": ast.NewIntValue(1),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
