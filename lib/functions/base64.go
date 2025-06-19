package functions

import (
	"encoding/base64"

	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/native"
)

func init() {
	LibFunctions["base64Encode"] = ast.NewBaseFunction(
		"base64Encode",
		append(
			native.MustNewNativeFunction("base64Encode", Base64Encode).Definitions(),
			native.MustNewNativeFunction("base64EncodeBytes", Base64EncodeBytes).Definitions()...,
		),
	)
	LibFunctions["base64Decode"] = ast.NewBaseFunction(
		"base64Decode",
		append(
			native.MustNewNativeFunction("base64Decode", Base64Decode).Definitions(),
			native.MustNewNativeFunction("base64DecodeBytes", Base64DecodeBytes).Definitions()...,
		),
	)
}

func Base64Encode(str string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(str)), nil
}

func Base64EncodeBytes(bytes []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func Base64Decode(str string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		if _, tryAltEncoding := err.(base64.CorruptInputError); tryAltEncoding {
			return base64.RawStdEncoding.DecodeString(str)
		}
		return nil, err
	}
	return decoded, nil
}

func Base64DecodeBytes(bytes []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(bytes))
	if err != nil {
		if _, tryAltEncoding := err.(base64.CorruptInputError); tryAltEncoding {
			return base64.RawStdEncoding.DecodeString(string(bytes))
		}
		return nil, err
	}
	return decoded, nil
}
