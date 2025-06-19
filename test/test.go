package test

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	expr "cel.dev/expr"
	testpb "cel.dev/expr/conformance/test"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/yywing/sl"
	"github.com/yywing/sl/ast"
)

func LoadTestFile(paths []string) []*testpb.SimpleTestFile {
	var files []*testpb.SimpleTestFile
	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read file %q: %v", path, err)
		}
		file := &testpb.SimpleTestFile{}
		err = prototext.Unmarshal(b, file)
		if err != nil {
			log.Fatalf("failed to parse file %q: %v", path, err)
		}
		files = append(files, file)
	}
	return files
}

func ExprValueToValue(i *expr.ExprValue) (ast.Value, error) {
	v := i.GetValue()
	var result ast.Value
	switch t := v.GetKind().(type) {
	case *expr.Value_NullValue:
		result = ast.NewNullValue()
	case *expr.Value_BoolValue:
		result = ast.NewBoolValue(t.BoolValue)
	case *expr.Value_Int64Value:
		result = ast.NewIntValue(t.Int64Value)
	case *expr.Value_Uint64Value:
		result = ast.NewUintValue(t.Uint64Value)
	case *expr.Value_DoubleValue:
		result = ast.NewDoubleValue(t.DoubleValue)
	case *expr.Value_StringValue:
		result = ast.NewStringValue(t.StringValue)
	case *expr.Value_BytesValue:
		result = ast.NewBytesValue(t.BytesValue)
	case *expr.Value_TypeValue:
		result = ast.NewTypeValue(t.TypeValue)
	case *expr.Value_ListValue:
		values := make([]ast.Value, len(t.ListValue.Values))
		var err error
		var valueType ast.ValueType
		for i, v := range t.ListValue.Values {
			val := &expr.ExprValue{Kind: &expr.ExprValue_Value{Value: v}}
			values[i], err = ExprValueToValue(val)
			if err != nil {
				return nil, err
			}
			if valueType == nil {
				valueType = values[i].Type()
			} else {
				if !valueType.Equals(values[i].Type()) {
					return nil, fmt.Errorf("list element %d has type %s, expected %s", i, values[i].Type(), valueType)
				}
			}
		}
		if valueType == nil {
			return nil, fmt.Errorf("list is empty")
		}
		result = ast.NewListValue(values, valueType)
	case *expr.Value_MapValue:
		values := make(map[ast.Value]ast.Value)
		var valueType ast.ValueType
		for i, v := range t.MapValue.Entries {
			key, err := ExprValueToValue(&expr.ExprValue{Kind: &expr.ExprValue_Value{Value: v.Key}})
			if err != nil {
				return nil, err
			}
			val := &expr.ExprValue{Kind: &expr.ExprValue_Value{Value: v.Value}}
			value, err := ExprValueToValue(val)
			if err != nil {
				return nil, err
			}
			if key.Type().Kind() != ast.TypeKindString {
				return nil, fmt.Errorf("map key %d has type %s, expected string", i, key.Type())
			}

			if valueType == nil {
				valueType = value.Type()
			} else {
				if !valueType.Equals(value.Type()) {
					return nil, fmt.Errorf("map value %d has type %s, expected %s", i, value.Type(), valueType)
				}
			}
			values[key] = value
		}
		if valueType == nil {
			return nil, fmt.Errorf("map is empty")
		}
		result = ast.NewMapValue(values, ast.StringType, valueType)
	default:
		return nil, fmt.Errorf("unknown type on transform")
	}
	return result, nil
}

func ValueToExprValue(res ast.Value) (*expr.Value, error) {
	switch res.Type().Kind() {
	case ast.TypeKindNull:
		return &expr.Value{Kind: &expr.Value_NullValue{NullValue: 0}}, nil
	case ast.TypeKindBool:
		return &expr.Value{Kind: &expr.Value_BoolValue{BoolValue: res.(*ast.BoolValue).BoolValue}}, nil
	case ast.TypeKindInt:
		return &expr.Value{Kind: &expr.Value_Int64Value{Int64Value: res.(*ast.IntValue).IntValue}}, nil
	case ast.TypeKindUint:
		return &expr.Value{Kind: &expr.Value_Uint64Value{Uint64Value: res.(*ast.UintValue).UintValue}}, nil
	case ast.TypeKindDouble:
		return &expr.Value{Kind: &expr.Value_DoubleValue{DoubleValue: res.(*ast.DoubleValue).DoubleValue}}, nil
	case ast.TypeKindString:
		return &expr.Value{Kind: &expr.Value_StringValue{StringValue: res.(*ast.StringValue).StringValue}}, nil
	case ast.TypeKindBytes:
		return &expr.Value{Kind: &expr.Value_BytesValue{BytesValue: res.(*ast.BytesValue).BytesValue}}, nil
	case ast.TypeKindType:
		return &expr.Value{Kind: &expr.Value_TypeValue{TypeValue: res.(*ast.TypeValue).Value}}, nil
	case ast.TypeKindList:
		values := res.(*ast.ListValue).ListValue
		exprValues := make([]*expr.Value, len(values))
		var err error
		for i, v := range values {
			exprValues[i], err = ValueToExprValue(v)
			if err != nil {
				return nil, err
			}
		}
		return &expr.Value{Kind: &expr.Value_ListValue{ListValue: &expr.ListValue{Values: exprValues}}}, nil
	case ast.TypeKindMap:
		values := res.(*ast.MapValue).MapValue
		exprValues := make([]*expr.MapValue_Entry, 0, len(values))
		for k, v := range values {
			key, err := ValueToExprValue(k)
			if err != nil {
				return nil, err
			}
			value, err := ValueToExprValue(v)
			if err != nil {
				return nil, err
			}
			exprValues = append(exprValues, &expr.MapValue_Entry{
				Key:   key,
				Value: value,
			})
		}
		return &expr.Value{Kind: &expr.Value_MapValue{MapValue: &expr.MapValue{Entries: exprValues}}}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", res.Type())
	}
}

func GetTestVariables(t *testpb.SimpleTest) (sl.Variables, error) {
	vars := make(sl.Variables)
	for k, v := range t.GetBindings() {
		value, err := ExprValueToValue(v)
		if err != nil {
			return nil, err
		}
		// Special handling for a.b.c format
		if strings.Contains(k, ".") {
			parts := strings.Split(k, ".")
			k = parts[0]
			slices.Reverse(parts)
			for i := 0; i < len(parts)-1; i++ {
				value = ast.NewMapValue(map[ast.Value]ast.Value{
					ast.NewStringValue(parts[i]): value,
				}, ast.StringType, value.Type())
			}
		}
		vars[k] = value
	}
	return vars, nil
}

func RunTestCase(testCase *testpb.SimpleTest) error {
	env := sl.NewStdEnv()

	ast, err := sl.Parse(testCase.GetExpr())
	if err != nil {
		return err
	}

	vars, err := GetTestVariables(testCase)
	if err != nil {
		return err
	}

	varTypes := vars.Type()
	program := sl.NewProgram(ast, varTypes)

	// Fill in default values
	if testCase.GetResultMatcher() == nil {
		testCase.ResultMatcher = &testpb.SimpleTest_Value{
			Value: &expr.Value{
				Kind: &expr.Value_BoolValue{
					BoolValue: true,
				},
			},
		}
	}

	// check
	if !testCase.GetDisableCheck() {
		typ, err := env.Check(program)
		switch m := testCase.GetResultMatcher().(type) {
		case *testpb.SimpleTest_EvalError:
			if err != nil {
				return nil
			}
		case *testpb.SimpleTest_Value:
			if err != nil {
				return fmt.Errorf("Check(%q) error: %v", testCase.GetName(), err)
			}
			if !MatchKind(m.Value, typ) {
				return fmt.Errorf("Check(%q) got %v, want %v", testCase.GetName(), typ, m.Value.Kind)
			}
		default:
			return fmt.Errorf("unexpected matcher kind: %T", testCase.GetResultMatcher())
		}
	}

	// eval
	if !testCase.GetCheckOnly() {
		result, err := env.Run(program, vars)
		switch m := testCase.GetResultMatcher().(type) {
		case *testpb.SimpleTest_EvalError:
			if err == nil {
				return fmt.Errorf("eval: got nil, want %v", m.EvalError)
			}
		case *testpb.SimpleTest_Value:
			if err != nil {
				return fmt.Errorf("eval(%q) error: %v", testCase.GetName(), err)
			}
			val, err := ValueToExprValue(result)
			if err != nil {
				return fmt.Errorf("ValueToExprValue(%q) error: %v", testCase.GetName(), err)
			}

			if diff := cmp.Diff(m.Value, val, protocmp.Transform(), protocmp.SortRepeatedFields(&expr.MapValue{}, "entries")); diff != "" {
				return fmt.Errorf("program.Eval() diff (-want +got):\n%s", diff)
			}

		default:
			return fmt.Errorf("unexpected matcher kind: %T", testCase.GetResultMatcher())
		}
	}
	return nil
}

func MatchKind(i *expr.Value, want ast.ValueType) bool {
	if want.Kind() == ast.TypeKindAny {
		return true
	}

	switch i.GetKind().(type) {
	case *expr.Value_NullValue:
		return want.Kind() == ast.TypeKindNull
	case *expr.Value_BoolValue:
		return want.Kind() == ast.TypeKindBool
	case *expr.Value_Int64Value:
		return want.Kind() == ast.TypeKindInt
	case *expr.Value_Uint64Value:
		return want.Kind() == ast.TypeKindUint
	case *expr.Value_DoubleValue:
		return want.Kind() == ast.TypeKindDouble
	case *expr.Value_StringValue:
		return want.Kind() == ast.TypeKindString
	case *expr.Value_BytesValue:
		return want.Kind() == ast.TypeKindBytes
	case *expr.Value_TypeValue:
		return want.Kind() == ast.TypeKindType
	case *expr.Value_ListValue:
		if want.Kind() != ast.TypeKindList {
			return false
		}
		values := i.GetListValue()
		if len(values.Values) == 0 {
			return true
		}
		return MatchKind(values.Values[0], want.(*ast.ListType).ElementType())
	case *expr.Value_MapValue:
		if want.Kind() != ast.TypeKindMap {
			return false
		}
		values := i.GetMapValue()
		if len(values.Entries) == 0 {
			return true
		}
		return MatchKind(values.Entries[0].Value, want.(*ast.MapType).ValueType())
	default:
		return false
	}
}
