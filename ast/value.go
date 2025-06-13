package ast

import "fmt"

// Value 表示运行时的值
type Value interface {
	Type() ValueType
	Equal(other Value) bool
	String() string
}

// 基础值类型
type BoolValue struct {
	BoolValue bool
}

func (v *BoolValue) Type() ValueType { return BoolType }
func (v *BoolValue) Equal(other Value) bool {
	if o, ok := other.(*BoolValue); ok {
		return v.BoolValue == o.BoolValue
	}
	return false
}
func (v *BoolValue) String() string { return fmt.Sprintf("%t", v.BoolValue) }

type IntValue struct {
	IntValue int64
}

func (v *IntValue) Type() ValueType { return IntType }
func (v *IntValue) Equal(other Value) bool {
	if o, ok := other.(*IntValue); ok {
		return v.IntValue == o.IntValue
	}
	return false
}
func (v *IntValue) String() string { return fmt.Sprintf("%d", v.IntValue) }

type UintValue struct {
	UintValue uint64
}

func (v *UintValue) Type() ValueType { return UintType }
func (v *UintValue) Equal(other Value) bool {
	if o, ok := other.(*UintValue); ok {
		return v.UintValue == o.UintValue
	}
	return false
}
func (v *UintValue) String() string { return fmt.Sprintf("%d", v.UintValue) }

type DoubleValue struct {
	DoubleValue float64
}

func (v *DoubleValue) Type() ValueType { return DoubleType }
func (v *DoubleValue) Equal(other Value) bool {
	if o, ok := other.(*DoubleValue); ok {
		return v.DoubleValue == o.DoubleValue
	}
	return false
}
func (v *DoubleValue) String() string { return fmt.Sprintf("%g", v.DoubleValue) }

type StringValue struct {
	StringValue string
}

func (v *StringValue) Type() ValueType { return StringType }
func (v *StringValue) Equal(other Value) bool {
	if o, ok := other.(*StringValue); ok {
		return v.StringValue == o.StringValue
	}
	return false
}
func (v *StringValue) String() string { return fmt.Sprintf("%q", v.StringValue) }

type BytesValue struct {
	BytesValue []byte
}

func (v *BytesValue) Type() ValueType { return BytesType }
func (v *BytesValue) Equal(other Value) bool {
	if o, ok := other.(*BytesValue); ok {
		return string(v.BytesValue) == string(o.BytesValue)
	}
	return false
}
func (v *BytesValue) String() string { return fmt.Sprintf("b%q", v.BytesValue) }

type ListValue struct {
	ListValue   []Value
	elementType ValueType
}

func (v *ListValue) Type() ValueType { return NewListType(v.elementType) }
func (v *ListValue) Equal(other Value) bool {
	if o, ok := other.(*ListValue); ok {
		if len(v.ListValue) != len(o.ListValue) {
			return false
		}
		for i, val := range v.ListValue {
			if !val.Equal(o.ListValue[i]) {
				return false
			}
		}
		return true
	}
	return false
}
func (v *ListValue) String() string {
	values := make([]string, len(v.ListValue))
	for i, val := range v.ListValue {
		values[i] = val.String()
	}
	return fmt.Sprintf("[%s]", fmt.Sprintf("%v", values))
}
func (v *ListValue) ElementType() ValueType { return v.elementType }

type MapValue struct {
	MapValue  map[Value]Value
	keyType   ValueType
	valueType ValueType
}

func (v *MapValue) Type() ValueType { return NewMapType(v.keyType, v.valueType) }
func (v *MapValue) Equal(other Value) bool {
	if o, ok := other.(*MapValue); ok {
		if len(v.MapValue) != len(o.MapValue) {
			return false
		}
		for k, val := range v.MapValue {
			if otherVal, exists := o.Get(k); !exists || !val.Equal(otherVal) {
				return false
			}
		}
		return true
	}
	return false
}
func (v *MapValue) String() string       { return fmt.Sprintf("%v", v.MapValue) }
func (v *MapValue) KeyType() ValueType   { return v.keyType }
func (v *MapValue) ValueType() ValueType { return v.valueType }

func (v *MapValue) Get(key Value) (Value, bool) {
	for k, v := range v.MapValue {
		if k.Equal(key) {
			return v, true
		}
	}
	return nil, false
}

func (m *MapValue) Set(key Value, value Value) {
	for k := range m.MapValue {
		if k.Equal(key) {
			m.MapValue[k] = value
			return
		}
	}
	m.MapValue[key] = value
}

type NullValue struct{}

func (v *NullValue) Type() ValueType        { return NullType }
func (v *NullValue) Equal(other Value) bool { _, ok := other.(*NullValue); return ok }
func (v *NullValue) String() string         { return "null" }

type TypeValue struct {
	Value string
}

func (v *TypeValue) Type() ValueType { return TypeType }
func (v *TypeValue) Equal(other Value) bool {
	if o, ok := other.(*TypeValue); ok {
		return v.Value == o.Value
	}
	return false
}
func (v *TypeValue) String() string { return fmt.Sprintf("type<%s>", v.Value) }

// 创建值的便利函数
func NewBoolValue(v bool) *BoolValue        { return &BoolValue{BoolValue: v} }
func NewIntValue(v int64) *IntValue         { return &IntValue{IntValue: v} }
func NewUintValue(v uint64) *UintValue      { return &UintValue{UintValue: v} }
func NewDoubleValue(v float64) *DoubleValue { return &DoubleValue{DoubleValue: v} }
func NewStringValue(v string) *StringValue  { return &StringValue{StringValue: v} }
func NewBytesValue(v []byte) *BytesValue    { return &BytesValue{BytesValue: v} }
func NewNullValue() *NullValue              { return &NullValue{} }
func NewTypeValue(v string) *TypeValue      { return &TypeValue{Value: v} }

func NewListValue(values []Value, elementType ValueType) *ListValue {
	return &ListValue{ListValue: values, elementType: elementType}
}

func NewMapValue(values map[Value]Value, keyType, valueType ValueType) *MapValue {
	return &MapValue{MapValue: values, keyType: keyType, valueType: valueType}
}
