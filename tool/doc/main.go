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
	flag.StringVar(&output, "output", "functions.md", "输出markdown文件路径")
	flag.Parse()

	env := sl.NewStdEnv()

	// 生成markdown内容
	content := generateMarkdown(env)

	// 写入文件
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

	// 创建表格头
	builder.WriteString("| Name | Input | Output |\n")
	builder.WriteString("|------|-------|--------|\n")

	// 获取所有函数名并排序
	functionNames := env.Functions()
	sort.Strings(functionNames)

	for _, name := range functionNames {
		function, _ := env.GetFunction(name)

		// 获取函数的所有类型签名
		types := function.Types()

		for i, fnType := range types {
			var functionName string
			if i == 0 {
				// 转义函数名中的管道符号
				functionName = "`" + strings.ReplaceAll(name, "|", "\\|") + "`"
			}

			// 格式化输入参数
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

			// 格式化输出类型
			outputStr := fmt.Sprintf("`%s`", fnType.ReturnType().String())

			// 添加表格行
			builder.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				functionName, inputStr, outputStr))
		}
	}

	return builder.String()
}
