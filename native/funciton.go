package native

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/yywing/sl/ast"
)

// NativeFunction represents a wrapper for a native Go function
type NativeFunction struct {
	name       string
	fn         reflect.Value
	paramTypes []ast.ValueType
	returnType ast.ValueType
	// Is a reverse list of parameters
	defaultArgs []reflect.Value
}

func (f *NativeFunction) Definitions() []ast.Definition {
	var defs []ast.Definition
	for i := 0; i <= len(f.defaultArgs); i++ {
		defs = append(defs, ast.Definition{
			Type: *ast.NewFunctionType(
				f.name,
				f.paramTypes[:len(f.paramTypes)-i],
				f.returnType,
			),
			Call: func(args []ast.Value) (ast.Value, error) {
				return f.call(args)
			},
		})
	}
	return defs
}

// From back to front in reverse order
func (f *NativeFunction) WithDefaultArg(value interface{}) *NativeFunction {
	f.defaultArgs = append(f.defaultArgs, reflect.ValueOf(value))
	return f
}

func (f *NativeFunction) Call(args []ast.Value) (result ast.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("function %s panicked: %v", f.name, r)
		}
	}()

	return f.call(args)
}

func (f *NativeFunction) call(args []ast.Value) (ast.Value, error) {
	// Check parameter count
	if len(args)+len(f.defaultArgs) < len(f.paramTypes) {
		return nil, fmt.Errorf("function %s expects at least %d arguments, got %d",
			f.name, len(f.paramTypes)-len(f.defaultArgs), len(args))
	}
	// Convert parameters to reflect.Value
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		goVal, err := ValueToGo(arg)
		if err != nil {
			return nil, fmt.Errorf("error converting argument %d for function %s: %v", i, f.name, err)
		}
		// Handle empty map and empty slice
		if goVal != nil {
			reflectArgs[i] = reflect.ValueOf(goVal)
		} else {
			reflectArgs[i] = reflect.Zero(goTypeToReflectType(f.paramTypes[i]))
		}
	}

	// Supplement default parameters
	if len(args) < len(f.paramTypes) {
		def := f.defaultArgs[:len(f.paramTypes)-len(args)]
		slices.Reverse(def)
		reflectArgs = append(reflectArgs, def...)
	}

	// Call function
	results := f.fn.Call(reflectArgs)

	// Handle return values
	if len(results) == 0 {
		return ast.NewNullValue(), nil
	}

	if len(results) == 1 {
		return ValueFromGo(results[0].Interface()), nil
	}

	// Multiple return values, check if the last one is error
	lastResult := results[len(results)-1]
	if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !lastResult.IsNil() {
			return nil, lastResult.Interface().(error)
		}
		// Remove error return value
		results = results[:len(results)-1]
	}

	if len(results) == 1 {
		return ValueFromGo(results[0].Interface()), nil
	}

	// Multiple return values, wrap as list
	values := make([]ast.Value, len(results))
	var elemType ast.ValueType = ast.StringType
	for i, res := range results {
		values[i] = ValueFromGo(res.Interface())
		if i == 0 {
			elemType = values[i].Type()
		}
	}
	return ast.NewListValue(values, elemType), nil
}

// Create Function from Go function
func NewNativeFunction(name string, fn interface{}) (*NativeFunction, error) {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected function, got %T", fn)
	}

	fnType := fnVal.Type()

	// Build parameter types
	paramTypes := make([]ast.ValueType, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramTypes[i] = goTypeToValueType(fnType.In(i))
	}

	// Build return type
	var returnType ast.ValueType = ast.NullType
	if fnType.NumOut() > 0 {
		// If the last return value is error, ignore it
		numOut := fnType.NumOut()
		if fnType.Out(numOut - 1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			numOut--
		}

		if numOut == 1 {
			returnType = goTypeToValueType(fnType.Out(0))
		} else if numOut > 1 {
			// Multiple return values, create list type
			returnType = ast.NewListType(ast.StringType) // Simplified handling
		}
	}

	return &NativeFunction{
		name:       name,
		fn:         fnVal,
		paramTypes: paramTypes,
		returnType: returnType,
	}, nil
}

func MustNewNativeFunction(name string, fn interface{}) *NativeFunction {
	nativeFn, err := NewNativeFunction(name, fn)
	if err != nil {
		panic(err)
	}
	return nativeFn
}
