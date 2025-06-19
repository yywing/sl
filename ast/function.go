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

// Function represents a callable function
type Function interface {
	Name() string
	// Supports multiple input parameter types
	Types() []FunctionType
	Call(args []Value) (Value, error)
}

type FunctionCall func(args []Value) (Value, error)

type Definition struct {
	Type FunctionType
	Call FunctionCall
}

type BaseFunction struct {
	name        string
	Definitions []Definition
}

func (f *BaseFunction) Name() string {
	return f.name
}

func (f *BaseFunction) Types() []FunctionType {
	results := make([]FunctionType, 0, len(f.Definitions))
	for _, d := range f.Definitions {
		results = append(results, d.Type)
	}
	return results
}

func (f *BaseFunction) Call(args []Value) (Value, error) {
	var argTypes []ValueType
	for _, arg := range args {
		argTypes = append(argTypes, arg.Type())
	}

	for _, d := range f.Definitions {
		ft := d.Type.ParamTypes()

		if len(ft) != len(args) {
			continue
		}

		_, ok := MatchFunctionTypes(ft, argTypes)
		if !ok {
			continue
		}

		return d.Call(args)

	}
	return nil, fmt.Errorf("no matching function definition found, with args %v", argTypes)
}

func (f *BaseFunction) AddDefinition(n Definition) error {
	for _, d := range f.Definitions {
		if d.Type.Equals(&n.Type) {
			return fmt.Errorf("function definition already exists")
		}
	}
	f.Definitions = append(f.Definitions, n)
	return nil
}

func (f *BaseFunction) Combine(other *BaseFunction) *BaseFunction {
	for _, d := range other.Definitions {
		if err := f.AddDefinition(d); err != nil {
			panic(err)
		}
	}
	return nil
}

// Predefined functions
func NewBaseFunction(name string, d []Definition) *BaseFunction {
	return &BaseFunction{name: name, Definitions: d}
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

func ResolveDynamicType(env map[string]ValueType, dyn ValueType, rt ValueType) (ValueType, error) {
	if !dyn.IsDyn() {
		return dyn, nil
	}

	switch dyn := dyn.(type) {
	case *ListType:
		var innerType ValueType
		if rt != nil {
			e, ok := rt.(*ListType)
			if !ok {
				return nil, fmt.Errorf("cannot resolve dynamic type %s to %s", dyn.String(), rt.String())
			}
			innerType = e.ElementType()
		}
		result, err := ResolveDynamicType(env, dyn.ElementType(), innerType)
		if err != nil {
			return nil, err
		}
		return NewListType(result), nil
	case *MapType:
		var keyType ValueType
		var valueType ValueType
		if rt != nil {
			m, ok := rt.(*MapType)
			if !ok {
				return nil, fmt.Errorf("cannot resolve dynamic type %s to %s", dyn.String(), rt.String())
			}
			keyType = m.KeyType()
			valueType = m.ValueType()
		}
		key, err := ResolveDynamicType(env, dyn.KeyType(), keyType)
		if err != nil {
			return nil, err
		}
		value, err := ResolveDynamicType(env, dyn.ValueType(), valueType)
		if err != nil {
			return nil, err
		}
		return NewMapType(key, value), nil
	default:
		t, ok := env[dyn.Kind()]
		if rt == nil {
			if !ok {
				return nil, fmt.Errorf("dynamic type %s not found", dyn.String())
			}
			return t, nil
		}

		if ok {
			if !t.Equals(rt) {
				return nil, fmt.Errorf("dynamic type should %s but got %s", t.String(), rt.String())
			}
			env[dyn.Kind()] = GetDeterministicType(t, rt)
		} else {
			env[dyn.Kind()] = rt
		}

		return env[dyn.Kind()], nil
	}
}

func MatchFunctionTypes(functionType []ValueType, argTypes []ValueType) (map[string]ValueType, bool) {
	if len(functionType) != len(argTypes) {
		return nil, false
	}

	var err error
	env := make(map[string]ValueType)
	for i, argType := range argTypes {
		paramType := functionType[i]
		if paramType.IsDyn() {
			paramType, err = ResolveDynamicType(env, paramType, argType)
			if err != nil {
				return nil, false
			}
		}
		if !TypeEquals(argType, paramType) {
			return nil, false
		}
	}
	return env, true
}

var (
	paramA  = NewValueTypeParamType("A")
	paramB  = NewValueTypeParamType("B")
	listOfA = NewListType(paramA)
	mapOfAB = NewMapType(paramA, paramB)

	LogicalAndFunction = NewBaseFunction(
		LogicalAnd,
		[]Definition{
			{
				Type: *NewFunctionType(LogicalAnd, []ValueType{BoolType, BoolType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*BoolValue).BoolValue && args[1].(*BoolValue).BoolValue), nil
				},
			},
		},
	)

	LogicalOrFunction = NewBaseFunction(
		LogicalOr,
		[]Definition{
			{
				Type: *NewFunctionType(LogicalOr, []ValueType{BoolType, BoolType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*BoolValue).BoolValue || args[1].(*BoolValue).BoolValue), nil
				},
			},
		},
	)

	LogicalNotFunction = NewBaseFunction(
		LogicalNot,
		[]Definition{
			{
				Type: *NewFunctionType(LogicalNot, []ValueType{BoolType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(!args[0].(*BoolValue).BoolValue), nil
				},
			},
		},
	)

	EqualsFunction = NewBaseFunction(
		Equals,
		[]Definition{
			{
				Type: *NewFunctionType(Equals, []ValueType{paramA, paramA}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].Equal(args[1])), nil
				},
			},
		},
	)

	NotEqualsFunction = NewBaseFunction(
		NotEquals,
		[]Definition{
			{
				Type: *NewFunctionType(NotEquals, []ValueType{paramA, paramA}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(!args[0].Equal(args[1])), nil
				},
			},
		},
	)

	AddFunction = NewBaseFunction(
		Add,
		[]Definition{
			{
				Type: *NewFunctionType(Add, []ValueType{BytesType, BytesType}, BytesType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*BytesValue).BytesValue
					y := args[1].(*BytesValue).BytesValue
					return NewBytesValue(append(x, y...)), nil
				},
			},
			{
				Type: *NewFunctionType(Add, []ValueType{DoubleType, DoubleType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return NewDoubleValue(args[0].(*DoubleValue).DoubleValue + args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(Add, []ValueType{IntType, IntType}, IntType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*IntValue).IntValue
					y := args[1].(*IntValue).IntValue
					if (y > 0 && x > math.MaxInt64-y) || (y < 0 && x < math.MinInt64-y) {
						return nil, fmt.Errorf("int overflow")
					}
					return NewIntValue(x + y), nil
				},
			},
			{
				Type: *NewFunctionType(Add, []ValueType{UintType, UintType}, UintType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*UintValue).UintValue
					y := args[1].(*UintValue).UintValue
					if y > math.MaxUint64-x {
						return nil, fmt.Errorf("uint overflow")
					}
					return NewUintValue(x + y), nil
				},
			},
			{
				Type: *NewFunctionType(Add, []ValueType{StringType, StringType}, StringType),
				Call: func(args []Value) (Value, error) {
					return NewStringValue(args[0].(*StringValue).StringValue + args[1].(*StringValue).StringValue), nil
				},
			},
			{
				Type: *NewFunctionType(Add, []ValueType{listOfA, listOfA}, listOfA),
				Call: func(args []Value) (Value, error) {
					return NewListValue(append(args[0].(*ListValue).ListValue, args[1].(*ListValue).ListValue...), args[0].(*ListValue).ElementType()), nil
				},
			},
		},
	)

	SubtractFunction = NewBaseFunction(
		Subtract,
		[]Definition{
			{
				Type: *NewFunctionType(Subtract, []ValueType{IntType, IntType}, IntType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*IntValue).IntValue
					y := args[1].(*IntValue).IntValue
					if (y < 0 && x > math.MaxInt64+y) || (y > 0 && x < math.MinInt64+y) {
						return nil, fmt.Errorf("int overflow")
					}
					return NewIntValue(x - y), nil
				},
			},
			{
				Type: *NewFunctionType(Subtract, []ValueType{UintType, UintType}, UintType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*UintValue).UintValue
					y := args[1].(*UintValue).UintValue
					if y > x {
						return nil, fmt.Errorf("uint overflow")
					}
					return NewUintValue(x - y), nil
				},
			},
			{
				Type: *NewFunctionType(Subtract, []ValueType{DoubleType, DoubleType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return NewDoubleValue(args[0].(*DoubleValue).DoubleValue - args[1].(*DoubleValue).DoubleValue), nil
				},
			},
		},
	)

	MultiplyFunction = NewBaseFunction(
		Multiply,
		[]Definition{
			{
				Type: *NewFunctionType(Multiply, []ValueType{IntType, IntType}, IntType),
				Call: func(args []Value) (Value, error) {
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
				},
			},
			{
				Type: *NewFunctionType(Multiply, []ValueType{UintType, UintType}, UintType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*UintValue).UintValue
					y := args[1].(*UintValue).UintValue
					if y != 0 && x > math.MaxUint64/y {
						return nil, fmt.Errorf("uint overflow")
					}
					return NewUintValue(x * y), nil
				},
			},
			{
				Type: *NewFunctionType(Multiply, []ValueType{DoubleType, DoubleType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return NewDoubleValue(args[0].(*DoubleValue).DoubleValue * args[1].(*DoubleValue).DoubleValue), nil
				},
			},
		},
	)

	DivideFunction = NewBaseFunction(
		Divide,
		[]Definition{
			{
				Type: *NewFunctionType(Divide, []ValueType{IntType, IntType}, IntType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*IntValue).IntValue
					y := args[1].(*IntValue).IntValue
					if y == 0 {
						return nil, fmt.Errorf("divide by zero")
					}
					if x == math.MinInt64 && y == -1 {
						return nil, fmt.Errorf("int overflow")
					}
					return NewIntValue(x / y), nil
				},
			},
			{
				Type: *NewFunctionType(Divide, []ValueType{UintType, UintType}, UintType),
				Call: func(args []Value) (Value, error) {
					if args[1].(*UintValue).UintValue == 0 {
						return nil, fmt.Errorf("divide by zero")
					}
					return NewUintValue(args[0].(*UintValue).UintValue / args[1].(*UintValue).UintValue), nil
				},
			},
			{
				Type: *NewFunctionType(Divide, []ValueType{DoubleType, DoubleType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return NewDoubleValue(args[0].(*DoubleValue).DoubleValue / args[1].(*DoubleValue).DoubleValue), nil
				},
			},
		},
	)

	ModuloFunction = NewBaseFunction(
		Modulo,
		[]Definition{
			{
				Type: *NewFunctionType(Modulo, []ValueType{IntType, IntType}, IntType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*IntValue).IntValue
					y := args[1].(*IntValue).IntValue
					if y == 0 {
						return nil, fmt.Errorf("modulo by zero")
					}
					if x == math.MinInt64 && y == -1 {
						return nil, fmt.Errorf("int overflow")
					}
					return NewIntValue(x % y), nil
				},
			},
			{
				Type: *NewFunctionType(Modulo, []ValueType{UintType, UintType}, UintType),
				Call: func(args []Value) (Value, error) {
					if args[1].(*UintValue).UintValue == 0 {
						return nil, fmt.Errorf("modulo by zero")
					}
					return NewUintValue(args[0].(*UintValue).UintValue % args[1].(*UintValue).UintValue), nil
				},
			},
		},
	)

	NegateFunction = NewBaseFunction(
		Negate,
		[]Definition{
			{
				Type: *NewFunctionType(Negate, []ValueType{IntType}, IntType),
				Call: func(args []Value) (Value, error) {
					x := args[0].(*IntValue).IntValue
					if x == math.MinInt64 {
						return nil, fmt.Errorf("int overflow")
					}
					return NewIntValue(-args[0].(*IntValue).IntValue), nil
				},
			},
			{
				Type: *NewFunctionType(Negate, []ValueType{DoubleType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return NewDoubleValue(-args[0].(*DoubleValue).DoubleValue), nil
				},
			},
		},
	)

	LessFunction = NewBaseFunction(
		Less,
		[]Definition{
			{
				Type: *NewFunctionType(Less, []ValueType{IntType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*IntValue).IntValue < args[1].(*IntValue).IntValue), nil
				},
			},
			{
				Type: *NewFunctionType(Less, []ValueType{IntType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*IntValue).IntValue) < args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(Less, []ValueType{IntType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*IntValue).IntValue) < float64(args[1].(*UintValue).UintValue)), nil
				},
			},
			{
				Type: *NewFunctionType(Less, []ValueType{UintType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*UintValue).UintValue < args[1].(*UintValue).UintValue), nil
				},
			},
			{
				Type: *NewFunctionType(Less, []ValueType{UintType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*UintValue).UintValue) < float64(args[1].(*IntValue).IntValue)), nil
				},
			},
			{
				Type: *NewFunctionType(Less, []ValueType{UintType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*UintValue).UintValue) < args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(Less, []ValueType{DoubleType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue < args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(Less, []ValueType{DoubleType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue < float64(args[1].(*IntValue).IntValue)), nil
				},
			},
			{
				Type: *NewFunctionType(Less, []ValueType{DoubleType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue < float64(args[1].(*UintValue).UintValue)), nil
				},
			},
		},
	)

	LessEqualsFunction = NewBaseFunction(
		LessEquals,
		[]Definition{
			{
				Type: *NewFunctionType(LessEquals, []ValueType{IntType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*IntValue).IntValue <= args[1].(*IntValue).IntValue), nil
				},
			},
			{
				Type: *NewFunctionType(LessEquals, []ValueType{IntType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*IntValue).IntValue) <= args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(LessEquals, []ValueType{IntType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*IntValue).IntValue) <= float64(args[1].(*UintValue).UintValue)), nil
				},
			},
			{
				Type: *NewFunctionType(LessEquals, []ValueType{UintType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*UintValue).UintValue <= args[1].(*UintValue).UintValue), nil
				},
			},
			{
				Type: *NewFunctionType(LessEquals, []ValueType{UintType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*UintValue).UintValue) <= float64(args[1].(*IntValue).IntValue)), nil
				},
			},
			{
				Type: *NewFunctionType(LessEquals, []ValueType{UintType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*UintValue).UintValue) <= args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(LessEquals, []ValueType{DoubleType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue <= args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(LessEquals, []ValueType{DoubleType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue <= float64(args[1].(*IntValue).IntValue)), nil
				},
			},
			{
				Type: *NewFunctionType(LessEquals, []ValueType{DoubleType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue <= float64(args[1].(*UintValue).UintValue)), nil
				},
			},
		},
	)

	GreaterFunction = NewBaseFunction(
		Greater,
		[]Definition{
			{
				Type: *NewFunctionType(Greater, []ValueType{IntType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*IntValue).IntValue > args[1].(*IntValue).IntValue), nil
				},
			},
			{
				Type: *NewFunctionType(Greater, []ValueType{IntType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*IntValue).IntValue) > args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(Greater, []ValueType{IntType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*IntValue).IntValue) > float64(args[1].(*UintValue).UintValue)), nil
				},
			},
			{
				Type: *NewFunctionType(Greater, []ValueType{UintType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*UintValue).UintValue > args[1].(*UintValue).UintValue), nil
				},
			},
			{
				Type: *NewFunctionType(Greater, []ValueType{UintType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*UintValue).UintValue) > float64(args[1].(*IntValue).IntValue)), nil
				},
			},
			{
				Type: *NewFunctionType(Greater, []ValueType{UintType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*UintValue).UintValue) > args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(Greater, []ValueType{DoubleType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue > args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(Greater, []ValueType{DoubleType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue > float64(args[1].(*IntValue).IntValue)), nil
				},
			},
			{
				Type: *NewFunctionType(Greater, []ValueType{DoubleType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue > float64(args[1].(*UintValue).UintValue)), nil
				},
			},
		},
	)

	GreaterEqualsFunction = NewBaseFunction(
		GreaterEquals,
		[]Definition{
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{IntType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*IntValue).IntValue >= args[1].(*IntValue).IntValue), nil
				},
			},
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{IntType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*IntValue).IntValue) >= args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{IntType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*IntValue).IntValue) >= float64(args[1].(*UintValue).UintValue)), nil
				},
			},
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{UintType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*UintValue).UintValue >= args[1].(*UintValue).UintValue), nil
				},
			},
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{UintType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*UintValue).UintValue) >= float64(args[1].(*IntValue).IntValue)), nil
				},
			},
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{UintType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(float64(args[0].(*UintValue).UintValue) >= args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{DoubleType, DoubleType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue >= args[1].(*DoubleValue).DoubleValue), nil
				},
			},
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{DoubleType, IntType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue >= float64(args[1].(*IntValue).IntValue)), nil
				},
			},
			{
				Type: *NewFunctionType(GreaterEquals, []ValueType{DoubleType, UintType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return NewBoolValue(args[0].(*DoubleValue).DoubleValue >= float64(args[1].(*UintValue).UintValue)), nil
				},
			},
		},
	)

	InFunction = NewBaseFunction(
		In,
		[]Definition{
			{
				Type: *NewFunctionType(In, []ValueType{paramA, listOfA}, BoolType),
				Call: func(args []Value) (Value, error) {
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
				},
			},
			{
				Type: *NewFunctionType(In, []ValueType{paramA, mapOfAB}, BoolType),
				Call: func(args []Value) (Value, error) {
					m, ok := args[1].(*MapValue)
					if !ok {
						return nil, fmt.Errorf("in expects map argument, got %T", args[1])
					}
					_, exists := m.Get(args[0])
					return NewBoolValue(exists), nil
				},
			},
		},
	)

	SizeFunction = NewBaseFunction(
		Size, []Definition{
			{
				Type: *NewFunctionType(Size, []ValueType{BytesType}, IntType),
				Call: func(args []Value) (Value, error) {
					return NewIntValue(int64(len(args[0].(*BytesValue).BytesValue))), nil
				},
			},
			{
				Type: *NewFunctionType(Size, []ValueType{StringType}, IntType),
				Call: func(args []Value) (Value, error) {
					return NewIntValue(int64(len([]rune(args[0].(*StringValue).StringValue)))), nil
				},
			},
			{
				Type: *NewFunctionType(Size, []ValueType{listOfA}, IntType),
				Call: func(args []Value) (Value, error) {
					return NewIntValue(int64(len(args[0].(*ListValue).ListValue))), nil
				},
			},
			{
				Type: *NewFunctionType(Size, []ValueType{mapOfAB}, IntType),
				Call: func(args []Value) (Value, error) {
					return NewIntValue(int64(len(args[0].(*MapValue).MapValue))), nil
				},
			},
		},
	)

	TypeFunction = NewBaseFunction(
		Type,
		[]Definition{
			{
				Type: *NewFunctionType(Type, []ValueType{paramA}, TypeType),
				Call: func(args []Value) (Value, error) {
					return NewTypeValue(args[0].Type().Kind()), nil
				},
			},
		},
	)

	BoolFunction = NewBaseFunction(
		Bool,
		[]Definition{
			{
				Type: *NewFunctionType(Bool, []ValueType{BoolType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], BoolType)
				},
			},
			{
				Type: *NewFunctionType(Bool, []ValueType{StringType}, BoolType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], BoolType)
				},
			},
		},
	)

	BytesFunction = NewBaseFunction(
		Bytes,
		[]Definition{
			{
				Type: *NewFunctionType(Bytes, []ValueType{BytesType}, BytesType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], BytesType)
				},
			},
			{
				Type: *NewFunctionType(Bytes, []ValueType{StringType}, BytesType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], BytesType)
				},
			},
		},
	)

	DoubleFunction = NewBaseFunction(
		Double,
		[]Definition{
			{
				Type: *NewFunctionType(Double, []ValueType{IntType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], DoubleType)
				},
			},
			{
				Type: *NewFunctionType(Double, []ValueType{UintType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], DoubleType)
				},
			},
			{
				Type: *NewFunctionType(Double, []ValueType{DoubleType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], DoubleType)
				},
			},
			{
				Type: *NewFunctionType(Double, []ValueType{StringType}, DoubleType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], DoubleType)
				},
			},
		},
	)

	IntFunction = NewBaseFunction(
		Int,
		[]Definition{
			{
				Type: *NewFunctionType(Int, []ValueType{DoubleType}, IntType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], IntType)
				},
			},
			{
				Type: *NewFunctionType(Int, []ValueType{UintType}, IntType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], IntType)
				},
			},
			{
				Type: *NewFunctionType(Int, []ValueType{IntType}, IntType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], IntType)
				},
			},
			{
				Type: *NewFunctionType(Int, []ValueType{StringType}, IntType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], IntType)
				},
			},
		},
	)

	UintFunction = NewBaseFunction(
		Uint,
		[]Definition{
			{
				Type: *NewFunctionType(Uint, []ValueType{DoubleType}, UintType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], UintType)
				},
			},
			{
				Type: *NewFunctionType(Uint, []ValueType{UintType}, UintType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], UintType)
				},
			},
			{
				Type: *NewFunctionType(Uint, []ValueType{IntType}, UintType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], UintType)
				},
			},
			{
				Type: *NewFunctionType(Uint, []ValueType{StringType}, UintType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], UintType)
				},
			},
		},
	)

	StringFunction = NewBaseFunction(
		String,
		[]Definition{
			{
				Type: *NewFunctionType(String, []ValueType{StringType}, StringType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], StringType)
				},
			},
			{
				Type: *NewFunctionType(String, []ValueType{BytesType}, StringType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], StringType)
				},
			},
			{
				Type: *NewFunctionType(String, []ValueType{BoolType}, StringType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], StringType)
				},
			},
			{
				Type: *NewFunctionType(String, []ValueType{DoubleType}, StringType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], StringType)
				},
			},
			{
				Type: *NewFunctionType(String, []ValueType{IntType}, StringType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], StringType)
				},
			},
			{
				Type: *NewFunctionType(String, []ValueType{UintType}, StringType),
				Call: func(args []Value) (Value, error) {
					return args[0].Type().ConvertTo(args[0], StringType)
				},
			},
		},
	)
)

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
