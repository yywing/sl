package ast

import (
	"fmt"
	"math"
)

const (
	LogicalAnd    = "_&&_"
	LogicalOr     = "_||_"
	LogicalNot    = "!_"
	Equals        = "_==_"
	NotEquals     = "_!=_"
	Less          = "_<_"
	LessEquals    = "_<=_"
	Greater       = "_>_"
	GreaterEquals = "_>=_"
	Add           = "_+_"
	Subtract      = "_-_"
	Multiply      = "_*_"
	Divide        = "_/_"
	Modulo        = "_%_"
	Negate        = "-_"
	In            = "_in_"

	Size   = "size"
	Type   = "type"
	Bool   = "bool"
	Bytes  = "bytes"
	Double = "double"
	Int    = "int"
	String = "string"
	Uint   = "uint"
)

// Function 表示一个可调用的函数
type Function interface {
	Name() string
	// 支持多种类型入参
	Types() []FunctionType
	Call(args []Value) (Value, error)
}

type BaseFunction struct {
	name  string
	types []FunctionType
	call  func(args []Value) (Value, error)
}

func (f *BaseFunction) Name() string {
	return f.name
}

func (f *BaseFunction) Types() []FunctionType {
	return f.types
}

func (f *BaseFunction) Call(args []Value) (Value, error) {
	return f.call(args)
}

// 预定义函数
func NewBaseFunction(name string, types []FunctionType, call func(args []Value) (Value, error)) Function {
	return &BaseFunction{name: name, types: types, call: call}
}

const ValueTypeParamTypeType = "param"

type ValueTypeParamType struct {
	*PrimitiveType
}

func (t *ValueTypeParamType) String() string {
	return "dyn_" + t.Kind()
}

func (t *ValueTypeParamType) IsDyn() bool {
	return true
}

func NewValueTypeParamType(name string) *ValueTypeParamType {
	return &ValueTypeParamType{PrimitiveType: &PrimitiveType{kind: name}}
}

var (
	paramA  = NewValueTypeParamType("A")
	paramB  = NewValueTypeParamType("B")
	listOfA = NewListType(paramA)
	mapOfAB = NewMapType(paramA, paramB)

	LogicalAndFunction = NewBaseFunction(LogicalAnd, []FunctionType{
		*NewFunctionType(LogicalAnd, []ValueType{BoolType, BoolType}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("logical and expects 2 arguments, got %d", len(args))
		}
		if args[0].Type() != BoolType || args[1].Type() != BoolType {
			return nil, fmt.Errorf("logical and expects bool arguments, got %s and %s", args[0].Type(), args[1].Type())
		}

		return NewBoolValue(args[0].(*BoolValue).BoolValue && args[1].(*BoolValue).BoolValue), nil
	})
	LogicalOrFunction = NewBaseFunction(LogicalOr, []FunctionType{
		*NewFunctionType(LogicalOr, []ValueType{BoolType, BoolType}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("logical or expects 2 arguments, got %d", len(args))
		}
		if args[0].Type() != BoolType || args[1].Type() != BoolType {
			return nil, fmt.Errorf("logical or expects bool arguments, got %s and %s", args[0].Type(), args[1].Type())
		}
		return NewBoolValue(args[0].(*BoolValue).BoolValue || args[1].(*BoolValue).BoolValue), nil
	})
	LogicalNotFunction = NewBaseFunction(LogicalNot, []FunctionType{
		*NewFunctionType(LogicalNot, []ValueType{BoolType}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("logical not expects 1 argument, got %d", len(args))
		}
		if args[0].Type() != BoolType {
			return nil, fmt.Errorf("logical not expects bool argument, got %s", args[0].Type())
		}
		return NewBoolValue(!args[0].(*BoolValue).BoolValue), nil
	})

	EqualsFunction = NewBaseFunction(Equals, []FunctionType{
		*NewFunctionType(Equals, []ValueType{paramA, paramA}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("equals expects 2 arguments, got %d", len(args))
		}
		if !args[0].Type().Equals(args[1].Type()) {
			return nil, fmt.Errorf("equals expects arguments of the same type, got %s and %s", args[0].Type(), args[1].Type())
		}

		return NewBoolValue(args[0].Equal(args[1])), nil
	})
	NotEqualsFunction = NewBaseFunction(NotEquals, []FunctionType{
		*NewFunctionType(NotEquals, []ValueType{paramA, paramA}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("not equals expects 2 arguments, got %d", len(args))
		}
		if !args[0].Type().Equals(args[1].Type()) {
			return nil, fmt.Errorf("not equals expects arguments of the same type, got %s and %s", args[0].Type(), args[1].Type())
		}

		return NewBoolValue(!args[0].Equal(args[1])), nil
	})

	AddFunction = NewBaseFunction(Add, []FunctionType{
		*NewFunctionType(Add, []ValueType{BytesType, BytesType}, BytesType),
		*NewFunctionType(Add, []ValueType{DoubleType, DoubleType}, DoubleType),
		*NewFunctionType(Add, []ValueType{IntType, IntType}, IntType),
		*NewFunctionType(Add, []ValueType{UintType, UintType}, UintType),
		*NewFunctionType(Add, []ValueType{StringType, StringType}, StringType),
		*NewFunctionType(Add, []ValueType{listOfA, listOfA}, listOfA),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("add expects 2 arguments, got %d", len(args))
		}

		if !args[0].Type().Equals(args[1].Type()) {
			return nil, fmt.Errorf("add expects arguments of the same type, got %s and %s", args[0].Type(), args[1].Type())
		}

		switch args[0].Type().Kind() {
		case TypeKindBytes:
			x := args[0].(*BytesValue).BytesValue
			y := args[1].(*BytesValue).BytesValue
			return NewBytesValue(append(x, y...)), nil
		case TypeKindInt:
			x := args[0].(*IntValue).IntValue
			y := args[1].(*IntValue).IntValue
			if (y > 0 && x > math.MaxInt64-y) || (y < 0 && x < math.MinInt64-y) {
				return nil, fmt.Errorf("int overflow")
			}
			return NewIntValue(x + y), nil
		case TypeKindUint:
			x := args[0].(*UintValue).UintValue
			y := args[1].(*UintValue).UintValue
			if y > math.MaxUint64-x {
				return nil, fmt.Errorf("uint overflow")
			}
			return NewUintValue(x + y), nil
		case TypeKindDouble:
			return NewDoubleValue(args[0].(*DoubleValue).DoubleValue + args[1].(*DoubleValue).DoubleValue), nil
		case TypeKindString:
			return NewStringValue(args[0].(*StringValue).StringValue + args[1].(*StringValue).StringValue), nil
		case TypeKindList:
			return NewListValue(append(args[0].(*ListValue).ListValue, args[1].(*ListValue).ListValue...), args[0].(*ListValue).ElementType()), nil
		}
		return nil, fmt.Errorf("unsupported add: %T + %T", args[0], args[1])
	})
	SubtractFunction = NewBaseFunction(Subtract, []FunctionType{
		*NewFunctionType(Subtract, []ValueType{IntType, IntType}, IntType),
		*NewFunctionType(Subtract, []ValueType{UintType, UintType}, UintType),
		*NewFunctionType(Subtract, []ValueType{DoubleType, DoubleType}, DoubleType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("subtract expects 2 arguments, got %d", len(args))
		}

		if !args[0].Type().Equals(args[1].Type()) {
			return nil, fmt.Errorf("subtract expects arguments of the same type, got %s and %s", args[0].Type(), args[1].Type())
		}

		switch args[0].Type() {
		case IntType:
			x := args[0].(*IntValue).IntValue
			y := args[1].(*IntValue).IntValue
			if (y < 0 && x > math.MaxInt64+y) || (y > 0 && x < math.MinInt64+y) {
				return nil, fmt.Errorf("int overflow")
			}
			return NewIntValue(x - y), nil
		case UintType:
			x := args[0].(*UintValue).UintValue
			y := args[1].(*UintValue).UintValue
			if y > x {
				return nil, fmt.Errorf("uint overflow")
			}
			return NewUintValue(x - y), nil
		case DoubleType:
			return NewDoubleValue(args[0].(*DoubleValue).DoubleValue - args[1].(*DoubleValue).DoubleValue), nil
		}
		return nil, fmt.Errorf("unsupported subtract: %T - %T", args[0], args[1])
	})
	MultiplyFunction = NewBaseFunction(Multiply, []FunctionType{
		*NewFunctionType(Multiply, []ValueType{IntType, IntType}, IntType),
		*NewFunctionType(Multiply, []ValueType{UintType, UintType}, UintType),
		*NewFunctionType(Multiply, []ValueType{DoubleType, DoubleType}, DoubleType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("multiply expects 2 arguments, got %d", len(args))
		}

		if !args[0].Type().Equals(args[1].Type()) {
			return nil, fmt.Errorf("multiply expects arguments of the same type, got %s and %s", args[0].Type(), args[1].Type())
		}

		switch args[0].Type() {
		case IntType:
			x := args[0].(*IntValue).IntValue
			y := args[1].(*IntValue).IntValue
			if (x == -1 && y == math.MinInt64) || (y == -1 && x == math.MinInt64) ||
				// x is positive, y is positive
				(x > 0 && y > 0 && x > math.MaxInt64/y) ||
				// x is positive, y is negative
				(x > 0 && y < 0 && y < math.MinInt64/x) ||
				// x is negative, y is positive
				(x < 0 && y > 0 && x < math.MinInt64/y) ||
				// x is negative, y is negative
				(x < 0 && y < 0 && y < math.MaxInt64/x) {
				return nil, fmt.Errorf("int overflow")
			}
			return NewIntValue(x * y), nil
		case UintType:
			x := args[0].(*UintValue).UintValue
			y := args[1].(*UintValue).UintValue
			if y != 0 && x > math.MaxUint64/y {
				return nil, fmt.Errorf("uint overflow")
			}
			return NewUintValue(x * y), nil
		case DoubleType:
			return NewDoubleValue(args[0].(*DoubleValue).DoubleValue * args[1].(*DoubleValue).DoubleValue), nil
		}
		return nil, fmt.Errorf("unsupported multiply: %T * %T", args[0], args[1])
	})
	DivideFunction = NewBaseFunction(Divide, []FunctionType{
		*NewFunctionType(Divide, []ValueType{IntType, IntType}, IntType),
		*NewFunctionType(Divide, []ValueType{UintType, UintType}, UintType),
		*NewFunctionType(Divide, []ValueType{DoubleType, DoubleType}, DoubleType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("divide expects 2 arguments, got %d", len(args))
		}

		if !args[0].Type().Equals(args[1].Type()) {
			return nil, fmt.Errorf("divide expects arguments of the same type, got %s and %s", args[0].Type(), args[1].Type())
		}

		switch args[0].Type() {
		case IntType:
			x := args[0].(*IntValue).IntValue
			y := args[1].(*IntValue).IntValue
			if y == 0 {
				return nil, fmt.Errorf("divide by zero")
			}
			if x == math.MinInt64 && y == -1 {
				return nil, fmt.Errorf("int overflow")
			}
			return NewIntValue(x / y), nil
		case UintType:
			if args[1].(*UintValue).UintValue == 0 {
				return nil, fmt.Errorf("divide by zero")
			}
			return NewUintValue(args[0].(*UintValue).UintValue / args[1].(*UintValue).UintValue), nil
		case DoubleType:
			return NewDoubleValue(args[0].(*DoubleValue).DoubleValue / args[1].(*DoubleValue).DoubleValue), nil
		}
		return nil, fmt.Errorf("unsupported divide: %T / %T", args[0], args[1])
	})
	ModuloFunction = NewBaseFunction(Modulo, []FunctionType{
		*NewFunctionType(Modulo, []ValueType{IntType, IntType}, IntType),
		*NewFunctionType(Modulo, []ValueType{UintType, UintType}, UintType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("modulo expects 2 arguments, got %d", len(args))
		}

		if !args[0].Type().Equals(args[1].Type()) {
			return nil, fmt.Errorf("modulo expects arguments of the same type, got %s and %s", args[0].Type(), args[1].Type())
		}

		switch args[0].Type() {
		case IntType:
			x := args[0].(*IntValue).IntValue
			y := args[1].(*IntValue).IntValue
			if y == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}
			if x == math.MinInt64 && y == -1 {
				return nil, fmt.Errorf("int overflow")
			}
			return NewIntValue(x % y), nil
		case UintType:
			if args[1].(*UintValue).UintValue == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}
			return NewUintValue(args[0].(*UintValue).UintValue % args[1].(*UintValue).UintValue), nil
		}
		return nil, fmt.Errorf("unsupported modulo: %T %T", args[0], args[1])
	})
	NegateFunction = NewBaseFunction(Negate, []FunctionType{
		*NewFunctionType(Negate, []ValueType{IntType}, IntType),
		*NewFunctionType(Negate, []ValueType{DoubleType}, DoubleType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("negate expects 1 argument, got %d", len(args))
		}

		if args[0].Type() != IntType && args[0].Type() != DoubleType {
			return nil, fmt.Errorf("negate expects int or double argument, got %T", args[0])
		}

		switch args[0].Type() {
		case IntType:
			x := args[0].(*IntValue).IntValue
			if x == math.MinInt64 {
				return nil, fmt.Errorf("int overflow")
			}
			return NewIntValue(-args[0].(*IntValue).IntValue), nil
		case DoubleType:
			return NewDoubleValue(-args[0].(*DoubleValue).DoubleValue), nil
		}
		return nil, fmt.Errorf("unsupported negate: -%T", args[0])
	})

	LessFunction = NewBaseFunction(Less, []FunctionType{
		*NewFunctionType(Less, []ValueType{IntType, IntType}, BoolType),
		*NewFunctionType(Less, []ValueType{IntType, DoubleType}, BoolType),
		*NewFunctionType(Less, []ValueType{IntType, UintType}, BoolType),
		*NewFunctionType(Less, []ValueType{UintType, UintType}, BoolType),
		*NewFunctionType(Less, []ValueType{UintType, IntType}, BoolType),
		*NewFunctionType(Less, []ValueType{UintType, DoubleType}, BoolType),
		*NewFunctionType(Less, []ValueType{DoubleType, DoubleType}, BoolType),
		*NewFunctionType(Less, []ValueType{DoubleType, IntType}, BoolType),
		*NewFunctionType(Less, []ValueType{DoubleType, UintType}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("less expects 2 arguments, got %d", len(args))
		}

		left, err := resolveNumber(args[0])
		if err != nil {
			return nil, err
		}
		right, err := resolveNumber(args[1])
		if err != nil {
			return nil, err
		}

		return NewBoolValue(left < right), nil
	})
	LessEqualsFunction = NewBaseFunction(LessEquals, []FunctionType{
		*NewFunctionType(LessEquals, []ValueType{IntType, IntType}, BoolType),
		*NewFunctionType(LessEquals, []ValueType{IntType, DoubleType}, BoolType),
		*NewFunctionType(LessEquals, []ValueType{IntType, UintType}, BoolType),
		*NewFunctionType(LessEquals, []ValueType{UintType, UintType}, BoolType),
		*NewFunctionType(LessEquals, []ValueType{UintType, IntType}, BoolType),
		*NewFunctionType(LessEquals, []ValueType{UintType, DoubleType}, BoolType),
		*NewFunctionType(LessEquals, []ValueType{DoubleType, DoubleType}, BoolType),
		*NewFunctionType(LessEquals, []ValueType{DoubleType, IntType}, BoolType),
		*NewFunctionType(LessEquals, []ValueType{DoubleType, UintType}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("less equals expects 2 arguments, got %d", len(args))
		}
		left, err := resolveNumber(args[0])
		if err != nil {
			return nil, err
		}
		right, err := resolveNumber(args[1])
		if err != nil {
			return nil, err
		}
		return NewBoolValue(left <= right), nil
	})

	GreaterFunction = NewBaseFunction(Greater, []FunctionType{
		*NewFunctionType(Greater, []ValueType{IntType, IntType}, BoolType),
		*NewFunctionType(Greater, []ValueType{IntType, DoubleType}, BoolType),
		*NewFunctionType(Greater, []ValueType{IntType, UintType}, BoolType),
		*NewFunctionType(Greater, []ValueType{UintType, UintType}, BoolType),
		*NewFunctionType(Greater, []ValueType{UintType, IntType}, BoolType),
		*NewFunctionType(Greater, []ValueType{UintType, DoubleType}, BoolType),
		*NewFunctionType(Greater, []ValueType{DoubleType, DoubleType}, BoolType),
		*NewFunctionType(Greater, []ValueType{DoubleType, IntType}, BoolType),
		*NewFunctionType(Greater, []ValueType{DoubleType, UintType}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("greater expects 2 arguments, got %d", len(args))
		}
		left, err := resolveNumber(args[0])
		if err != nil {
			return nil, err
		}
		right, err := resolveNumber(args[1])
		if err != nil {
			return nil, err
		}
		return NewBoolValue(left > right), nil
	})

	GreaterEqualsFunction = NewBaseFunction(GreaterEquals, []FunctionType{
		*NewFunctionType(GreaterEquals, []ValueType{IntType, IntType}, BoolType),
		*NewFunctionType(GreaterEquals, []ValueType{IntType, DoubleType}, BoolType),
		*NewFunctionType(GreaterEquals, []ValueType{IntType, UintType}, BoolType),
		*NewFunctionType(GreaterEquals, []ValueType{UintType, UintType}, BoolType),
		*NewFunctionType(GreaterEquals, []ValueType{UintType, IntType}, BoolType),
		*NewFunctionType(GreaterEquals, []ValueType{UintType, DoubleType}, BoolType),
		*NewFunctionType(GreaterEquals, []ValueType{DoubleType, DoubleType}, BoolType),
		*NewFunctionType(GreaterEquals, []ValueType{DoubleType, IntType}, BoolType),
		*NewFunctionType(GreaterEquals, []ValueType{DoubleType, UintType}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("greater equals expects 2 arguments, got %d", len(args))
		}
		left, err := resolveNumber(args[0])
		if err != nil {
			return nil, err
		}
		right, err := resolveNumber(args[1])
		if err != nil {
			return nil, err
		}
		return NewBoolValue(left >= right), nil
	})
	InFunction = NewBaseFunction(In, []FunctionType{
		*NewFunctionType(In, []ValueType{paramA, listOfA}, BoolType),
		*NewFunctionType(In, []ValueType{paramA, mapOfAB}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("in expects 2 arguments, got %d", len(args))
		}

		switch args[1].Type().Kind() {
		case TypeKindList:
			list, ok := args[1].(*ListValue)
			if !ok {
				return nil, fmt.Errorf("in expects list argument, got %T", args[1])
			}
			for _, elem := range list.ListValue {
				if elem.Equal(args[0]) {
					return NewBoolValue(true), nil
				}
			}
			return NewBoolValue(false), nil
		case TypeKindMap:
			m, ok := args[1].(*MapValue)
			if !ok {
				return nil, fmt.Errorf("in expects map argument, got %T", args[1])
			}
			_, exists := m.Get(args[0])
			return NewBoolValue(exists), nil
		default:
			return nil, fmt.Errorf("in expects list or map argument, got %T", args[1])
		}
	})
	SizeFunction = NewBaseFunction(Size, []FunctionType{
		*NewFunctionType(Size, []ValueType{BytesType}, IntType),
		*NewFunctionType(Size, []ValueType{StringType}, IntType),
		*NewFunctionType(Size, []ValueType{listOfA}, IntType),
		*NewFunctionType(Size, []ValueType{mapOfAB}, IntType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("size expects 1 argument, got %d", len(args))
		}
		switch args[0].Type().Kind() {
		case TypeKindBytes:
			return NewIntValue(int64(len(args[0].(*BytesValue).BytesValue))), nil
		case TypeKindString:
			return NewIntValue(int64(len([]rune(args[0].(*StringValue).StringValue)))), nil
		case TypeKindList:
			return NewIntValue(int64(len(args[0].(*ListValue).ListValue))), nil
		case TypeKindMap:
			return NewIntValue(int64(len(args[0].(*MapValue).MapValue))), nil
		default:
			return nil, fmt.Errorf("size expects bytes, string, list or map argument, got %T", args[0])
		}
	})
	TypeFunction = NewBaseFunction(Type, []FunctionType{
		*NewFunctionType(Type, []ValueType{paramA}, TypeType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("type expects 1 argument, got %d", len(args))
		}
		return NewTypeValue(args[0].Type().Kind()), nil
	})

	// TOOD: 这部分需要优化，添加类型的时候如何动态的添加转换函数
	BoolFunction = NewBaseFunction(Bool, []FunctionType{
		*NewFunctionType(Bool, []ValueType{BoolType}, BoolType),
		*NewFunctionType(Bool, []ValueType{StringType}, BoolType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("bool expects 1 argument, got %d", len(args))
		}
		return args[0].Type().ConvertTo(args[0], BoolType)
	})
	BytesFunction = NewBaseFunction(Bytes, []FunctionType{
		*NewFunctionType(Bytes, []ValueType{BytesType}, BytesType),
		*NewFunctionType(Bytes, []ValueType{StringType}, BytesType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("bytes expects 1 argument, got %d", len(args))
		}
		return args[0].Type().ConvertTo(args[0], BytesType)
	})
	DoubleFunction = NewBaseFunction(Double, []FunctionType{
		*NewFunctionType(Double, []ValueType{IntType}, DoubleType),
		*NewFunctionType(Double, []ValueType{UintType}, DoubleType),
		*NewFunctionType(Double, []ValueType{DoubleType}, DoubleType),
		*NewFunctionType(Double, []ValueType{StringType}, DoubleType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("double expects 1 argument, got %d", len(args))
		}
		return args[0].Type().ConvertTo(args[0], DoubleType)
	})
	IntFunction = NewBaseFunction(Int, []FunctionType{
		*NewFunctionType(Int, []ValueType{DoubleType}, IntType),
		*NewFunctionType(Int, []ValueType{UintType}, IntType),
		*NewFunctionType(Int, []ValueType{IntType}, IntType),
		*NewFunctionType(Int, []ValueType{StringType}, IntType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("int expects 1 argument, got %d", len(args))
		}
		return args[0].Type().ConvertTo(args[0], IntType)
	})
	UintFunction = NewBaseFunction(Uint, []FunctionType{
		*NewFunctionType(Uint, []ValueType{DoubleType}, UintType),
		*NewFunctionType(Uint, []ValueType{IntType}, UintType),
		*NewFunctionType(Uint, []ValueType{UintType}, UintType),
		*NewFunctionType(Uint, []ValueType{StringType}, UintType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("uint expects 1 argument, got %d", len(args))
		}
		return args[0].Type().ConvertTo(args[0], UintType)
	})
	StringFunction = NewBaseFunction(String, []FunctionType{
		*NewFunctionType(String, []ValueType{StringType}, StringType),
		*NewFunctionType(String, []ValueType{BytesType}, StringType),
		*NewFunctionType(String, []ValueType{BoolType}, StringType),
		*NewFunctionType(String, []ValueType{DoubleType}, StringType),
		*NewFunctionType(String, []ValueType{IntType}, StringType),
		*NewFunctionType(String, []ValueType{UintType}, StringType),
	}, func(args []Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("string expects 1 argument, got %d", len(args))
		}
		return args[0].Type().ConvertTo(args[0], StringType)
	})
)

func resolveNumber(val Value) (float64, error) {
	switch val.Type() {
	case IntType:
		return float64(val.(*IntValue).IntValue), nil
	case UintType:
		return float64(val.(*UintValue).UintValue), nil
	case DoubleType:
		return val.(*DoubleValue).DoubleValue, nil
	default:
		return 0, fmt.Errorf("unsupported number: %T", val)
	}
}

var BuiltinFunctions = map[string]Function{
	LogicalAnd:    LogicalAndFunction,
	LogicalOr:     LogicalOrFunction,
	LogicalNot:    LogicalNotFunction,
	Equals:        EqualsFunction,
	NotEquals:     NotEqualsFunction,
	Less:          LessFunction,
	LessEquals:    LessEqualsFunction,
	Greater:       GreaterFunction,
	GreaterEquals: GreaterEqualsFunction,
	Add:           AddFunction,
	Subtract:      SubtractFunction,
	Multiply:      MultiplyFunction,
	Divide:        DivideFunction,
	Modulo:        ModuloFunction,
	Negate:        NegateFunction,
	In:            InFunction,

	Size:   SizeFunction,
	Type:   TypeFunction,
	Bool:   BoolFunction,
	Bytes:  BytesFunction,
	Double: DoubleFunction,
	Int:    IntFunction,
	Uint:   UintFunction,
	String: StringFunction,
}
