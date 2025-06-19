# sl

Subset of CEL (Common Expression Language)

## difference with cel

Test in [test case](./test/). And all `skipTests` is not support.

## usage

```golang
package main

import (
	"fmt"

	"github.com/yywing/sl"
)

func main() {
	env := sl.NewStdEnv()
	ast, err := sl.Parse("1 + 2")
	if err != nil {
		panic(err)
	}

	program := sl.NewProgram(ast, nil)

	// check
	_, err = env.Check(program)
	if err != nil {
		panic(err)
	}

	// run
	result, err := env.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
```

## doc

```bash
./scripts/gen_doc.sh
```

### function

[fucntion doc](./docs/functions.md)
