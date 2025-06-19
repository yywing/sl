package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/yywing/sl"
)

func main() {
	var output string
	flag.StringVar(&output, "output", "functions.md", "Output markdown file path")
	flag.Parse()

	env := sl.NewStdEnv()

	// Generate markdown content
	content := generateMarkdown(env)

	// Write to file
	err := os.WriteFile(output, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("success: %s\n", output)
}

func generateMarkdown(env *sl.Env) string {
	var builder strings.Builder

	builder.WriteString("# Functions\n\n")

	// Create table header
	builder.WriteString("| Name | Input | Output |\n")
	builder.WriteString("|------|-------|--------|\n")

	// Get all function names and sort them
	functionNames := env.Functions()
	sort.Strings(functionNames)

	for _, name := range functionNames {
		function, _ := env.GetFunction(name)

		// Get all type signatures for the function
		types := function.Types()

		for i, fnType := range types {
			var functionName string
			if i == 0 {
				// Escape pipe symbols in function names
				functionName = "`" + strings.ReplaceAll(name, "|", "\\|") + "`"
			}

			// Format input parameters
			paramTypes := fnType.ParamTypes()
			var inputStr string
			if len(paramTypes) == 0 {
				inputStr = "-"
			} else {
				paramStrings := make([]string, len(paramTypes))
				for j, paramType := range paramTypes {
					paramStrings[j] = fmt.Sprintf("`%s`", paramType.String())
				}
				inputStr = strings.Join(paramStrings, ", ")
			}

			// Format output type
			outputStr := fmt.Sprintf("`%s`", fnType.ReturnType().String())

			// Add table row
			builder.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				functionName, inputStr, outputStr))
		}
	}

	return builder.String()
}
