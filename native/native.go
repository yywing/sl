package native

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/yywing/sl/ast"
)

// Convert expression language values to Go values
func ValueToGo(v ast.Value) (interface{}, error) {
	switch val := v.(type) {
	case *ast.BoolValue:
		return val.BoolValue, nil
	case *ast.IntValue:
		return val.IntValue, nil
	case *ast.UintValue:
		return val.UintValue, nil
	case *ast.DoubleValue:
		return val.DoubleValue, nil
	case *ast.StringValue:
		return val.StringValue, nil
	case *ast.BytesValue:
		return val.BytesValue, nil
	case *ast.NullValue:
		return nil, nil
	case *ast.ListValue:
		if len(val.ListValue) == 0 {
			return nil, nil
		}

		// Convert the first element first to determine the type
		firstGoVal, err := ValueToGo(val.ListValue[0])
		if err != nil {
			return nil, err
		}

		// Check if all elements are of the same type
		firstType := reflect.TypeOf(firstGoVal)
		allSameType := true
		goVals := make([]interface{}, len(val.ListValue))
		goVals[0] = firstGoVal

		for i := 1; i < len(val.ListValue); i++ {
			goVal, err := ValueToGo(val.ListValue[i])
			if err != nil {
				return nil, err
			}
			goVals[i] = goVal
			if reflect.TypeOf(goVal) != firstType {
				allSameType = false
			}
		}

		// If all elements are of the same type, create a concrete typed slice
		if allSameType && firstType != nil {
			sliceType := reflect.SliceOf(firstType)
			slice := reflect.MakeSlice(sliceType, len(goVals), len(goVals))
			for i, v := range goVals {
				slice.Index(i).Set(reflect.ValueOf(v))
			}
			return slice.Interface(), nil
		} else {
			// Different types or nil values, return []interface{}
			return goVals, nil
		}
	case *ast.MapValue:
		if len(val.MapValue) == 0 {
			return nil, nil
		}

		// Process the first value first to determine the value type
		var firstKeyType reflect.Type
		var firstValueType reflect.Type
		goMap := make(map[interface{}]interface{})
		allSameKeyType := true
		allSameValueType := true

		for k, item := range val.MapValue {
			goKey, err := ValueToGo(k)
			if err != nil {
				return nil, err
			}
			goVal, err := ValueToGo(item)
			if err != nil {
				return nil, err
			}
			goMap[goKey] = goVal

			if firstKeyType == nil && goKey != nil {
				firstKeyType = reflect.TypeOf(goKey)
			} else if goKey != nil && reflect.TypeOf(goKey) != firstKeyType {
				allSameKeyType = false
			}

			if firstValueType == nil && goVal != nil {
				firstValueType = reflect.TypeOf(goVal)
			} else if goVal != nil && reflect.TypeOf(goVal) != firstValueType {
				allSameValueType = false
			}
		}

		// Create corresponding map type based on key and value type consistency
		var keyType, valueType reflect.Type

		if allSameKeyType && firstKeyType != nil {
			keyType = firstKeyType
		} else {
			keyType = reflect.TypeOf((*interface{})(nil)).Elem()
		}

		if allSameValueType && firstValueType != nil {
			valueType = firstValueType
		} else {
			valueType = reflect.TypeOf((*interface{})(nil)).Elem()
		}

		// Create concrete typed map
		mapType := reflect.MapOf(keyType, valueType)
		typedMap := reflect.MakeMap(mapType)
		for k, v := range goMap {
			typedMap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
		return typedMap.Interface(), nil
	default:
		return v, nil
	}
}

// Convert Go values to expression language values
func ValueFromGo(v interface{}) ast.Value {
	if v == nil {
		return ast.NewNullValue()
	}

	if val, ok := v.(ast.Value); ok {
		return val
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Bool:
		return ast.NewBoolValue(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ast.NewIntValue(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ast.NewUintValue(val.Uint())
	case reflect.Float32, reflect.Float64:
		return ast.NewDoubleValue(val.Float())
	case reflect.String:
		return ast.NewStringValue(val.String())
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			return ast.NewBytesValue(val.Bytes())
		}
		// Generic slice - determine element type based on Go reflection type
		values := make([]ast.Value, val.Len())
		elemType := goTypeToValueType(val.Type().Elem()) // Use actual element type
		for i := 0; i < val.Len(); i++ {
			values[i] = ValueFromGo(val.Index(i).Interface())
		}
		return ast.NewListValue(values, elemType)
	case reflect.Map:
		values := make(map[ast.Value]ast.Value)
		// Determine key and value types based on Go reflection type
		keyType := goTypeToValueType(val.Type().Key())
		valueType := goTypeToValueType(val.Type().Elem())

		for _, key := range val.MapKeys() {
			keyVal := ValueFromGo(key.Interface())
			mapVal := ValueFromGo(val.MapIndex(key).Interface())
			values[keyVal] = mapVal
		}
		return ast.NewMapValue(values, keyType, valueType)
	default:
		return ast.NewStringValue(fmt.Sprintf("%v", v))
	}
}

// Convert Go types to expression language types
func goTypeToValueType(t reflect.Type) ast.ValueType {
	switch t.Kind() {
	case reflect.Bool:
		return ast.BoolType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ast.IntType
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ast.UintType
	case reflect.Float32, reflect.Float64:
		return ast.DoubleType
	case reflect.String:
		return ast.StringType
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return ast.BytesType
		}
		return ast.NewListType(goTypeToValueType(t.Elem()))
	case reflect.Map:
		return ast.NewMapType(goTypeToValueType(t.Key()), goTypeToValueType(t.Elem()))
	case reflect.Pointer:
		// Check if pointer type implements ast.Value interface
		if t.Implements(reflect.TypeOf((*ast.Value)(nil)).Elem()) {
			// Create a zero value instance to get type information
			zeroVal := reflect.Zero(t).Interface()
			if val, ok := zeroVal.(ast.Value); ok {
				return val.Type()
			}
		}
		panic(fmt.Sprintf("unsupported pointer type: %v", t))
	default:
		panic(fmt.Sprintf("unsupported type: %T", t))
	}
}

// Convert expression language types to Go types
func goTypeToReflectType(t ast.ValueType) reflect.Type {
	switch t.Kind() {
	case ast.TypeKindBool:
		return reflect.TypeOf(bool(false))
	case ast.TypeKindInt:
		return reflect.TypeOf(int64(0))
	case ast.TypeKindUint:
		return reflect.TypeOf(uint64(0))
	case ast.TypeKindDouble:
		return reflect.TypeOf(float64(0))
	case ast.TypeKindString:
		return reflect.TypeOf("")
	case ast.TypeKindBytes:
		return reflect.TypeOf([]byte{})
	case ast.TypeKindNull:
		return reflect.TypeOf((*interface{})(nil)).Elem()
	case ast.TypeKindAny:
		return reflect.TypeOf(any(nil))
	case ast.TypeKindList:
		// Get concrete ListType
		if listType, ok := t.(*ast.ListType); ok {
			elemType := goTypeToReflectType(listType.ElementType())
			return reflect.SliceOf(elemType)
		}
		// Fallback to []interface{}
		return reflect.TypeOf([]interface{}{})
	case ast.TypeKindMap:
		// Get concrete MapType
		if mapType, ok := t.(*ast.MapType); ok {
			keyType := goTypeToReflectType(mapType.KeyType())
			valueType := goTypeToReflectType(mapType.ValueType())
			return reflect.MapOf(keyType, valueType)
		}
		// Fallback to map[string]interface{}
		return reflect.TypeOf(map[string]interface{}{})
	default:
		return reflect.TypeOf((*interface{})(nil)).Elem()
	}
}

// NativeFunction represents a wrapper for a native Go function
type NativeFunction struct {
	name string
	fn   reflect.Value
	typ  *ast.FunctionType
	// Is a reverse list of parameters
	defaultArgs []reflect.Value
}

func (f *NativeFunction) Name() string {
	return f.name
}

func (f *NativeFunction) Types() []ast.FunctionType {
	var types []ast.FunctionType
	for i := 0; i <= len(f.defaultArgs); i++ {
		types = append(types, *ast.NewFunctionType(
			f.name,
			f.typ.ParamTypes()[:len(f.typ.ParamTypes())-i],
			f.typ.ReturnType(),
		))
	}
	return types
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
	if len(args)+len(f.defaultArgs) < len(f.typ.ParamTypes()) {
		return nil, fmt.Errorf("function %s expects at least %d arguments, got %d",
			f.name, len(f.typ.ParamTypes())-len(f.defaultArgs), len(args))
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
			reflectArgs[i] = reflect.Zero(goTypeToReflectType(f.typ.ParamTypes()[i]))
		}
	}

	// Supplement default parameters
	if len(args) < len(f.typ.ParamTypes()) {
		def := f.defaultArgs[:len(f.typ.ParamTypes())-len(args)]
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

	funcType := ast.NewFunctionType(name, paramTypes, returnType)

	return &NativeFunction{
		name: name,
		fn:   fnVal,
		typ:  funcType,
	}, nil
}

func MustNewNativeFunction(name string, fn interface{}) *NativeFunction {
	nativeFn, err := NewNativeFunction(name, fn)
	if err != nil {
		panic(err)
	}
	return nativeFn
}
