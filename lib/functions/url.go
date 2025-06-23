package functions

import (
	"net/url"

	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/native"
)

func init() {
	LibFunctions[FunctionURLDecode] = ast.NewBaseFunction(
		FunctionURLDecode,
		native.MustNewNativeFunction(FunctionURLDecode, URLDecode).Definitions(),
	)
	LibFunctions[FunctionURLEncode] = ast.NewBaseFunction(
		FunctionURLEncode,
		native.MustNewNativeFunction(FunctionURLEncode, URLEncode).Definitions(),
	)
}

const (
	FunctionURLDecode = "urlDecode"
	FunctionURLEncode = "urlEncode"
)

func URLDecode(d string) (string, error) {
	return url.QueryUnescape(d)
}

func URLEncode(d string) string {
	return url.QueryEscape(d)
}
