package native

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/yywing/sl/ast"
)

// 将表达式语言的值转换为Go值
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

		// 先转换第一个元素，确定类型
		firstGoVal, err := ValueToGo(val.ListValue[0])
		if err != nil {
			return nil, err
		}

		// 检查所有元素是否为同一类型
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

		// 如果所有元素类型相同，创建具体类型的slice
		if allSameType && firstType != nil {
			sliceType := reflect.SliceOf(firstType)
			slice := reflect.MakeSlice(sliceType, len(goVals), len(goVals))
			for i, v := range goVals {
				slice.Index(i).Set(reflect.ValueOf(v))
			}
			return slice.Interface(), nil
		} else {
			// 类型不同或有nil值，返回[]interface{}
			return goVals, nil
		}
	case *ast.MapValue:
		if len(val.MapValue) == 0 {
			return nil, nil
		}

		// 先处理第一个值，确定值的类型
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

		// 根据key和value类型的一致性，创建相应类型的map
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

		// 创建具体类型的map
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

// 从Go值转换为表达式语言值
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
		// 通用slice - 根据Go反射类型确定元素类型
		values := make([]ast.Value, val.Len())
		elemType := goTypeToValueType(val.Type().Elem()) // 使用实际的元素类型
		for i := 0; i < val.Len(); i++ {
			values[i] = ValueFromGo(val.Index(i).Interface())
		}
		return ast.NewListValue(values, elemType)
	case reflect.Map:
		values := make(map[ast.Value]ast.Value)
		// 根据Go反射类型确定key和value类型
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

// 将Go类型转换为表达式语言类型
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
		// 检查指针类型是否实现了ast.Value接口
		if t.Implements(reflect.TypeOf((*ast.Value)(nil)).Elem()) {
			// 创建一个零值实例来获取类型信息
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

// 将表达式语言类型转换为Go类型
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
		// 获取具体的ListType
		if listType, ok := t.(*ast.ListType); ok {
			elemType := goTypeToReflectType(listType.ElementType())
			return reflect.SliceOf(elemType)
		}
		// 回退到[]interface{}
		return reflect.TypeOf([]interface{}{})
	case ast.TypeKindMap:
		// 获取具体的MapType
		if mapType, ok := t.(*ast.MapType); ok {
			keyType := goTypeToReflectType(mapType.KeyType())
			valueType := goTypeToReflectType(mapType.ValueType())
			return reflect.MapOf(keyType, valueType)
		}
		// 回退到map[string]interface{}
		return reflect.TypeOf(map[string]interface{}{})
	default:
		return reflect.TypeOf((*interface{})(nil)).Elem()
	}
}

// NativeFunction 表示一个原生Go函数的包装
type NativeFunction struct {
	name string
	fn   reflect.Value
	typ  *ast.FunctionType
	// 是参数的倒序列表
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

// 从后倒序往前
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
	// 检查参数数量
	if len(args)+len(f.defaultArgs) < len(f.typ.ParamTypes()) {
		return nil, fmt.Errorf("function %s expects at least %d arguments, got %d",
			f.name, len(f.typ.ParamTypes())-len(f.defaultArgs), len(args))
	}
	// 将参数转换为reflect.Value
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		goVal, err := ValueToGo(arg)
		if err != nil {
			return nil, fmt.Errorf("error converting argument %d for function %s: %v", i, f.name, err)
		}
		// 处理一下 空map 和 空slice
		if goVal != nil {
			reflectArgs[i] = reflect.ValueOf(goVal)
		} else {
			reflectArgs[i] = reflect.Zero(goTypeToReflectType(f.typ.ParamTypes()[i]))
		}
	}

	// 补足默认参数
	if len(args) < len(f.typ.ParamTypes()) {
		def := f.defaultArgs[:len(f.typ.ParamTypes())-len(args)]
		slices.Reverse(def)
		reflectArgs = append(reflectArgs, def...)
	}

	// 调用函数
	results := f.fn.Call(reflectArgs)

	// 处理返回值
	if len(results) == 0 {
		return ast.NewNullValue(), nil
	}

	if len(results) == 1 {
		return ValueFromGo(results[0].Interface()), nil
	}

	// 多个返回值，检查最后一个是否是error
	lastResult := results[len(results)-1]
	if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !lastResult.IsNil() {
			return nil, lastResult.Interface().(error)
		}
		// 移除error返回值
		results = results[:len(results)-1]
	}

	if len(results) == 1 {
		return ValueFromGo(results[0].Interface()), nil
	}

	// 多个返回值，包装为列表
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

// 从Go函数创建Function
func NewNativeFunction(name string, fn interface{}) (*NativeFunction, error) {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected function, got %T", fn)
	}

	fnType := fnVal.Type()

	// 构建参数类型
	paramTypes := make([]ast.ValueType, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramTypes[i] = goTypeToValueType(fnType.In(i))
	}

	// 构建返回类型
	var returnType ast.ValueType = ast.NullType
	if fnType.NumOut() > 0 {
		// 如果最后一个返回值是error，忽略它
		numOut := fnType.NumOut()
		if fnType.Out(numOut - 1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			numOut--
		}

		if numOut == 1 {
			returnType = goTypeToValueType(fnType.Out(0))
		} else if numOut > 1 {
			// 多个返回值，创建列表类型
			returnType = ast.NewListType(ast.StringType) // 简化处理
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
