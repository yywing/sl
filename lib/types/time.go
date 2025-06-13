package types

import (
	"fmt"

	"github.com/yywing/sl/ast"
)

const (
	TypeKindTimestamp = "timestamp"
	TypeKindDuration  = "duration"
)

var (
	TimestampType = ast.NewPrimitiveType(TypeKindTimestamp, 0)
	DurationType  = ast.NewPrimitiveType(TypeKindDuration, 0)
)

type TimestampValue struct {
	UnixNano int64
	TZ       string
}

func (v *TimestampValue) Type() ast.ValueType {
	return TimestampType
}

func (v *TimestampValue) String() string {
	return fmt.Sprintf("%d", v.UnixNano)
}

func (v *TimestampValue) Equal(other ast.Value) bool {
	otherValue, ok := other.(*TimestampValue)
	if !ok {
		return false
	}
	return v.UnixNano == otherValue.UnixNano && v.TZ == otherValue.TZ
}

type DurationValue struct {
	Nanosecond int64
}

func (v *DurationValue) Type() ast.ValueType {
	return DurationType
}

func (v *DurationValue) String() string {
	return fmt.Sprintf("%d", v.Nanosecond)
}

func (v *DurationValue) Equal(other ast.Value) bool {
	otherValue, ok := other.(*DurationValue)
	if !ok {
		return false
	}
	return v.Nanosecond == otherValue.Nanosecond
}
