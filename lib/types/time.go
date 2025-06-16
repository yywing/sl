package types

import (
	"fmt"

	"github.com/yywing/sl/ast"
)

const (
	TypeKindTimestamp = "timestamp"
	TypeKindDuration  = "duration"

	// Number of seconds between `0001-01-01T00:00:00Z` and the Unix epoch.
	MinUnixTime int64 = -62135596800
	// Number of seconds between `9999-12-31T23:59:59.999999999Z` and the Unix epoch.
	MaxUnixTime int64 = 253402300799
)

var (
	TimestampType = ast.NewPrimitiveType(TypeKindTimestamp, 0)
	DurationType  = ast.NewPrimitiveType(TypeKindDuration, 0)
)

type TimestampValue struct {
	UnixNano int64
	TZ       string
}

func NewTimestampValue(unixNano int64, tz string) *TimestampValue {
	return &TimestampValue{
		UnixNano: unixNano,
		TZ:       tz,
	}
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

func NewDurationValue(nanosecond int64) *DurationValue {
	return &DurationValue{
		Nanosecond: nanosecond,
	}
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
