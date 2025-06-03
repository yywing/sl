package ast

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	// SelectorType
	SelectorType = 1 << iota
)

type Selector interface {
	GetValueType(string) ValueType
}

// ValueType 表示表达式语言中的类型
type ValueType interface {
	Kind() string
	String() string
	HasTrait(trait int) bool
	Equals(other ValueType) bool
	ConvertTo(v Value, ty ValueType) (Value, error)

	// dyn
	IsDyn() bool

	// selector
	Member(name string) ValueType
}

const (
	TypeKindBool     = "bool"
	TypeKindInt      = "int"
	TypeKindUint     = "uint"
	TypeKindDouble   = "double"
	TypeKindString   = "string"
	TypeKindBytes    = "bytes"
	TypeKindList     = "list"
	TypeKindMap      = "map"
	TypeKindFunction = "function"
	TypeKindNull     = "null_type"
	TypeKindType     = "type"
	TypeKindAny      = "any"
)

// 基础类型实现
type PrimitiveType struct {
	kind      string
	traitMask int
	convert   func(v Value, ty ValueType) (Value, error)
	equals    func(other ValueType) bool
}

func (t *PrimitiveType) Kind() string {
	return t.kind
}

// HasTrait implements the ref.Type interface method.
func (t *PrimitiveType) HasTrait(trait int) bool {
	return trait&t.traitMask == trait
}

func (t *PrimitiveType) Equals(other ValueType) bool {
	if t.equals != nil {
		return t.equals(other)
	}

	if t.kind == TypeKindAny || other.Kind() == TypeKindAny {
		return true
	}

	if o, ok := other.(*PrimitiveType); ok {
		return t.kind == o.kind
	}
	return false
}

func (t *PrimitiveType) ConvertTo(v Value, ty ValueType) (Value, error) {
	if v.Type().Equals(ty) {
		return v, nil
	}
	if t.convert != nil {
		return t.convert(v, ty)
	}
	return nil, fmt.Errorf("cannot convert %s to %s", v.Type().String(), ty.String())
}

func (t *PrimitiveType) Member(name string) ValueType {
	panic("not implemented")
}

func (t *PrimitiveType) IsDyn() bool {
	return false
}

func (t *PrimitiveType) String() string {
	return t.kind
}

// 预定义类型
var (
	BoolType = &PrimitiveType{
		kind:      TypeKindBool,
		traitMask: 0,
		convert: func(v Value, ty ValueType) (Value, error) {
			switch ty.Kind() {
			case TypeKindString:
				return NewStringValue(strconv.FormatBool(v.(*BoolValue).BoolValue)), nil
			default:
				return nil, fmt.Errorf("cannot convert %s to %s", v.Type().String(), ty.String())
			}
		},
	}
	IntType = &PrimitiveType{
		kind:      TypeKindInt,
		traitMask: 0,
		convert: func(v Value, ty ValueType) (Value, error) {
			switch ty.Kind() {
			case TypeKindDouble:
				return NewDoubleValue(float64(v.(*IntValue).IntValue)), nil
			case TypeKindUint:
				if v.(*IntValue).IntValue < 0 {
					return nil, fmt.Errorf("int value %d is too small to convert to uint", v.(*IntValue).IntValue)
				}
				return NewUintValue(uint64(v.(*IntValue).IntValue)), nil
			case TypeKindString:
				return NewStringValue(strconv.FormatInt(v.(*IntValue).IntValue, 10)), nil
			default:
				return nil, fmt.Errorf("cannot convert %s to %s", v.Type().String(), ty.String())
			}
		},
	}
	UintType = &PrimitiveType{
		kind:      TypeKindUint,
		traitMask: 0,
		convert: func(v Value, ty ValueType) (Value, error) {
			switch ty.Kind() {
			case TypeKindInt:
				if v.(*UintValue).UintValue > math.MaxInt64 {
					return nil, fmt.Errorf("uint value %d is too large to convert to int", v.(*UintValue).UintValue)
				}
				return NewIntValue(int64(v.(*UintValue).UintValue)), nil
			case TypeKindDouble:
				return NewDoubleValue(float64(v.(*UintValue).UintValue)), nil
			case TypeKindString:
				return NewStringValue(strconv.FormatUint(v.(*UintValue).UintValue, 10)), nil
			default:
				return nil, fmt.Errorf("cannot convert %s to %s", v.Type().String(), ty.String())
			}
		},
	}
	DoubleType = &PrimitiveType{
		kind:      TypeKindDouble,
		traitMask: 0,
		convert: func(v Value, ty ValueType) (Value, error) {
			switch ty.Kind() {
			case TypeKindInt:
				if v.(*DoubleValue).DoubleValue >= math.MaxInt64 {
					return nil, fmt.Errorf("double value %f is too large to convert to int", v.(*DoubleValue).DoubleValue)
				}
				if v.(*DoubleValue).DoubleValue <= math.MinInt64 {
					return nil, fmt.Errorf("double value %f is too small to convert to int", v.(*DoubleValue).DoubleValue)
				}
				return NewIntValue(int64(v.(*DoubleValue).DoubleValue)), nil
			case TypeKindUint:
				if v.(*DoubleValue).DoubleValue < 0 {
					return nil, fmt.Errorf("double value %f is too small to convert to uint", v.(*DoubleValue).DoubleValue)
				}
				if v.(*DoubleValue).DoubleValue >= math.MaxUint64 {
					return nil, fmt.Errorf("double value %f is too large to convert to uint", v.(*DoubleValue).DoubleValue)
				}
				return NewUintValue(uint64(v.(*DoubleValue).DoubleValue)), nil
			case TypeKindString:
				return NewStringValue(strconv.FormatFloat(v.(*DoubleValue).DoubleValue, 'f', -1, 64)), nil
			default:
				return nil, fmt.Errorf("cannot convert %s to %s", v.Type().String(), ty.String())
			}
		},
	}
	StringType = &PrimitiveType{
		kind:      TypeKindString,
		traitMask: 0,
		convert: func(v Value, ty ValueType) (Value, error) {
			switch ty.Kind() {
			case TypeKindBool:
				b, err := strconv.ParseBool(v.(*StringValue).StringValue)
				if err != nil {
					return nil, err
				}
				return NewBoolValue(b), nil
			case TypeKindBytes:
				return NewBytesValue([]byte(v.(*StringValue).StringValue)), nil
			case TypeKindDouble:
				f, err := strconv.ParseFloat(v.(*StringValue).StringValue, 64)
				if err != nil {
					return nil, err
				}
				return NewDoubleValue(f), nil
			case TypeKindInt:
				i, err := strconv.ParseInt(v.(*StringValue).StringValue, 10, 64)
				if err != nil {
					return nil, err
				}
				return NewIntValue(i), nil
			case TypeKindUint:
				u, err := strconv.ParseUint(v.(*StringValue).StringValue, 10, 64)
				if err != nil {
					return nil, err
				}
				return NewUintValue(u), nil
			default:
				return nil, fmt.Errorf("cannot convert %s to %s", v.Type().String(), ty.String())
			}
		},
	}
	BytesType = &PrimitiveType{
		kind:      TypeKindBytes,
		traitMask: 0,
		convert: func(v Value, ty ValueType) (Value, error) {
			switch ty.Kind() {
			case TypeKindString:
				if !utf8.Valid(v.(*BytesValue).BytesValue) {
					return nil, fmt.Errorf("invalid UTF-8 in bytes, cannot convert to string")
				}
				return NewStringValue(string(v.(*BytesValue).BytesValue)), nil
			default:
				return nil, fmt.Errorf("cannot convert %s to %s", v.Type().String(), ty.String())
			}
		},
	}
	NullType = &PrimitiveType{
		kind:      TypeKindNull,
		traitMask: 0,
	}
	TypeType = &PrimitiveType{
		kind:      TypeKindType,
		traitMask: 0,
	}
	AnyType = &PrimitiveType{
		kind:      TypeKindAny,
		traitMask: 0,
		equals: func(other ValueType) bool {
			return true
		},
	}
)

// 列表类型
type ListType struct {
	*PrimitiveType
	elementType ValueType
}

func (t *ListType) Equals(other ValueType) bool {
	if other.Kind() == TypeKindAny {
		return true
	}

	if o, ok := other.(*ListType); ok {
		return t.elementType.Equals(o.elementType)
	}
	return false
}

func (t *ListType) String() string {
	return fmt.Sprintf("list<%s>", t.elementType.String())
}

func (t *ListType) ElementType() ValueType {
	return t.elementType
}

func (t *ListType) IsDyn() bool {
	return t.elementType.IsDyn()
}

func NewListType(elementType ValueType) *ListType {
	return &ListType{PrimitiveType: &PrimitiveType{kind: TypeKindList, traitMask: 0}, elementType: elementType}
}

// 映射类型
type MapType struct {
	*PrimitiveType
	keyType   ValueType
	valueType ValueType
}

func (t *MapType) Equals(other ValueType) bool {
	if other.Kind() == TypeKindAny {
		return true
	}

	if o, ok := other.(*MapType); ok {
		return t.keyType.Equals(o.keyType) && t.valueType.Equals(o.valueType)
	}
	return false
}

func (t *MapType) String() string {
	return fmt.Sprintf("map<%s, %s>", t.keyType.String(), t.valueType.String())
}

func (t *MapType) KeyType() ValueType {
	return t.keyType
}

func (t *MapType) ValueType() ValueType {
	return t.valueType
}

func (t *MapType) Member(name string) ValueType {
	return t.valueType
}

func (t *MapType) IsDyn() bool {
	return t.keyType.IsDyn() || t.valueType.IsDyn()
}

func NewMapType(keyType, valueType ValueType) *MapType {
	return &MapType{PrimitiveType: &PrimitiveType{kind: TypeKindMap, traitMask: SelectorType}, keyType: keyType, valueType: valueType}
}

// 函数类型
type FunctionType struct {
	*PrimitiveType
	name       string
	paramTypes []ValueType
	returnType ValueType
}

func (t *FunctionType) Equals(other ValueType) bool {
	if o, ok := other.(*FunctionType); ok {
		if t.name != o.name {
			return false
		}
		if len(t.paramTypes) != len(o.paramTypes) {
			return false
		}
		for i, p := range t.paramTypes {
			if !p.Equals(o.paramTypes[i]) {
				return false
			}
		}
		return t.returnType.Equals(o.returnType)
	}
	return false
}

func (t *FunctionType) Name() string {
	return t.name
}

func (t *FunctionType) ParamTypes() []ValueType {
	return t.paramTypes
}

func (t *FunctionType) ReturnType() ValueType {
	return t.returnType
}

func (t *FunctionType) String() string {
	params := make([]string, len(t.paramTypes))
	for i, param := range t.paramTypes {
		params[i] = param.String()
	}
	return fmt.Sprintf("%s(%s) -> %s", t.name, strings.Join(params, ", "), t.returnType.String())
}

func (t *FunctionType) IsDyn() bool {
	for _, param := range t.paramTypes {
		if param.IsDyn() {
			return true
		}
	}
	return t.returnType.IsDyn()
}

func NewFunctionType(name string, paramTypes []ValueType, returnType ValueType) *FunctionType {
	return &FunctionType{PrimitiveType: &PrimitiveType{kind: TypeKindFunction, traitMask: 0}, name: name, paramTypes: paramTypes, returnType: returnType}
}
