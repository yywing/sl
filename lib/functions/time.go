package functions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/lib/types"
	"github.com/yywing/sl/native"
)

func init() {
	LibFunctions["getFullYear"] = native.MustNewNativeFunction("getFullYear", GetFullYear).WithDefaultArg("")
	LibFunctions["getMonth"] = native.MustNewNativeFunction("getMonth", GetMonth).WithDefaultArg("")
	LibFunctions["getDayOfYear"] = native.MustNewNativeFunction("getDayOfYear", GetDayOfYear).WithDefaultArg("")
	LibFunctions["getDate"] = native.MustNewNativeFunction("getDate", GetDayOfMonthZeroBased).WithDefaultArg("")
	LibFunctions["getDayOfMonth"] = native.MustNewNativeFunction("getDayOfMonth", GetDayOfMonthOneBased).WithDefaultArg("")
	LibFunctions["getDayOfWeek"] = native.MustNewNativeFunction("getDayOfWeek", GetDayOfWeek).WithDefaultArg("")
	LibFunctions["getHours"] = GetHoursFunction
	LibFunctions["getMinutes"] = GetMinutesFunction
	LibFunctions["getSeconds"] = GetSecondsFunction
	LibFunctions["getMilliseconds"] = GetMillisecondsFunction
}

var (
	GetHoursFunction = ast.NewBaseFunction("getHours", []ast.FunctionType{
		*ast.NewFunctionType("getHours", []ast.ValueType{types.TimestampType}, ast.IntType),
		*ast.NewFunctionType("getHours", []ast.ValueType{types.DurationType}, ast.IntType),
		*ast.NewFunctionType("getHours", []ast.ValueType{types.TimestampType, ast.StringType}, ast.IntType),
	}, func(args []ast.Value) (ast.Value, error) {
		switch len(args) {
		case 1:
			switch args[0].Type().Kind() {
			case types.TypeKindTimestamp:
				result, err := TimestampGetHours(args[0].(*types.TimestampValue), "")
				if err != nil {
					return nil, err
				}
				return ast.NewIntValue(result), nil
			case types.TypeKindDuration:
				return ast.NewIntValue(DurationGetHours(args[0].(*types.DurationValue))), nil
			default:
				return nil, fmt.Errorf("getHours expects timestamp or duration argument, got %s", args[0].Type())
			}
		case 2:
			t, ok := args[0].(*types.TimestampValue)
			if !ok {
				return nil, fmt.Errorf("getHours expects timestamp argument, got %s", args[0].Type())
			}
			tz, ok := args[1].(*ast.StringValue)
			if !ok {
				return nil, fmt.Errorf("getHours expects string argument, got %s", args[1].Type())
			}
			result, err := TimestampGetHours(t, tz.StringValue)
			if err != nil {
				return nil, err
			}
			return ast.NewIntValue(result), nil
		default:
			return nil, fmt.Errorf("getHours expects 1 or 2 arguments, got %d", len(args))
		}
	})
	GetMinutesFunction = ast.NewBaseFunction("getMinutes", []ast.FunctionType{
		*ast.NewFunctionType("getMinutes", []ast.ValueType{types.TimestampType}, ast.IntType),
		*ast.NewFunctionType("getMinutes", []ast.ValueType{types.DurationType}, ast.IntType),
		*ast.NewFunctionType("getMinutes", []ast.ValueType{types.TimestampType, ast.StringType}, ast.IntType),
	}, func(args []ast.Value) (ast.Value, error) {
		switch len(args) {
		case 1:
			switch args[0].Type().Kind() {
			case types.TypeKindTimestamp:
				result, err := TimestampGetMinutes(args[0].(*types.TimestampValue), "")
				if err != nil {
					return nil, err
				}
				return ast.NewIntValue(result), nil
			case types.TypeKindDuration:
				return ast.NewIntValue(DurationGetMinutes(args[0].(*types.DurationValue))), nil
			default:
				return nil, fmt.Errorf("getMinutes expects timestamp or duration argument, got %s", args[0].Type())
			}
		case 2:
			t, ok := args[0].(*types.TimestampValue)
			if !ok {
				return nil, fmt.Errorf("getMinutes expects timestamp argument, got %s", args[0].Type())
			}
			tz, ok := args[1].(*ast.StringValue)
			if !ok {
				return nil, fmt.Errorf("getMinutes expects string argument, got %s", args[1].Type())
			}
			result, err := TimestampGetMinutes(t, tz.StringValue)
			if err != nil {
				return nil, err
			}
			return ast.NewIntValue(result), nil
		default:
			return nil, fmt.Errorf("getMinutes expects 1 or 2 arguments, got %d", len(args))
		}
	})
	GetSecondsFunction = ast.NewBaseFunction("getSeconds", []ast.FunctionType{
		*ast.NewFunctionType("getSeconds", []ast.ValueType{types.TimestampType}, ast.IntType),
		*ast.NewFunctionType("getSeconds", []ast.ValueType{types.DurationType}, ast.IntType),
		*ast.NewFunctionType("getSeconds", []ast.ValueType{types.TimestampType, ast.StringType}, ast.IntType),
	}, func(args []ast.Value) (ast.Value, error) {
		switch len(args) {
		case 1:
			switch args[0].Type().Kind() {
			case types.TypeKindTimestamp:
				result, err := TimestampGetSeconds(args[0].(*types.TimestampValue), "")
				if err != nil {
					return nil, err
				}
				return ast.NewIntValue(result), nil
			case types.TypeKindDuration:
				return ast.NewIntValue(DurationGetSeconds(args[0].(*types.DurationValue))), nil
			default:
				return nil, fmt.Errorf("getSeconds expects timestamp or duration argument, got %s", args[0].Type())
			}
		case 2:
			t, ok := args[0].(*types.TimestampValue)
			if !ok {
				return nil, fmt.Errorf("getSeconds expects timestamp argument, got %s", args[0].Type())
			}
			tz, ok := args[1].(*ast.StringValue)
			if !ok {
				return nil, fmt.Errorf("getSeconds expects string argument, got %s", args[1].Type())
			}
			result, err := TimestampGetSeconds(t, tz.StringValue)
			if err != nil {
				return nil, err
			}
			return ast.NewIntValue(result), nil
		default:
			return nil, fmt.Errorf("getSeconds expects 1 or 2 arguments, got %d", len(args))
		}
	})
	GetMillisecondsFunction = ast.NewBaseFunction("getMilliseconds", []ast.FunctionType{
		*ast.NewFunctionType("getMilliseconds", []ast.ValueType{types.TimestampType}, ast.IntType),
		*ast.NewFunctionType("getMilliseconds", []ast.ValueType{types.DurationType}, ast.IntType),
		*ast.NewFunctionType("getMilliseconds", []ast.ValueType{types.TimestampType, ast.StringType}, ast.IntType),
	}, func(args []ast.Value) (ast.Value, error) {
		switch len(args) {
		case 1:
			switch args[0].Type().Kind() {
			case types.TypeKindTimestamp:
				result, err := TimestampGetMilliseconds(args[0].(*types.TimestampValue), "")
				if err != nil {
					return nil, err
				}
				return ast.NewIntValue(result), nil
			case types.TypeKindDuration:
				return ast.NewIntValue(DurationGetMilliseconds(args[0].(*types.DurationValue))), nil
			default:
				return nil, fmt.Errorf("getMilliseconds expects timestamp or duration argument, got %s", args[0].Type())
			}
		case 2:
			t, ok := args[0].(*types.TimestampValue)
			if !ok {
				return nil, fmt.Errorf("getMilliseconds expects timestamp argument, got %s", args[0].Type())
			}
			tz, ok := args[1].(*ast.StringValue)
			if !ok {
				return nil, fmt.Errorf("getMilliseconds expects string argument, got %s", args[1].Type())
			}
			result, err := TimestampGetMilliseconds(t, tz.StringValue)
			if err != nil {
				return nil, err
			}
			return ast.NewIntValue(result), nil
		default:
			return nil, fmt.Errorf("getMilliseconds expects 1 or 2 arguments, got %d", len(args))
		}
	})
)

func timeZone(val string) (*time.Location, error) {
	ind := strings.Index(val, ":")
	if ind == -1 {
		loc, err := time.LoadLocation(val)
		if err != nil {
			return nil, err
		}
		return loc, nil
	}

	// If the input is not the name of a timezone (for example, 'US/Central'), it should be a numerical offset from UTC
	// in the format ^(+|-)(0[0-9]|1[0-4]):[0-5][0-9]$. The numerical input is parsed in terms of hours and minutes.
	hr, err := strconv.Atoi(string(val[0:ind]))
	if err != nil {
		return nil, err
	}
	min, err := strconv.Atoi(string(val[ind+1:]))
	if err != nil {
		return nil, err
	}
	var offset int
	if string(val[0]) == "-" {
		offset = hr*60 - min
	} else {
		offset = hr*60 + min
	}
	secondsEastOfUTC := int((time.Duration(offset) * time.Minute).Seconds())
	return time.FixedZone(val, secondsEastOfUTC), nil
}

func loadTimestamp(v *types.TimestampValue, tz string) (*time.Time, error) {
	useTZ := v.TZ
	if tz != "" {
		useTZ = tz
	}

	loc, err := timeZone(useTZ)
	if err != nil {
		return nil, err
	}
	t := time.Unix(0, v.UnixNano).In(loc)
	return &t, nil
}

func GetFullYear(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Year()), nil
}
func GetMonth(v *types.TimestampValue, tz string) (int64, error) {
	// CEL spec indicates that the month should be 0-based, but the Time value
	// for Month() is 1-based.
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Month() - 1), nil
}
func GetDayOfYear(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.YearDay() - 1), nil
}
func GetDayOfMonthZeroBased(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Day() - 1), nil
}
func GetDayOfMonthOneBased(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Day()), nil
}
func GetDayOfWeek(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Weekday()), nil
}
func TimestampGetHours(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Hour()), nil
}
func TimestampGetMinutes(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Minute()), nil
}
func TimestampGetSeconds(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Second()), nil
}
func TimestampGetMilliseconds(v *types.TimestampValue, tz string) (int64, error) {
	t, err := loadTimestamp(v, tz)
	if err != nil {
		return 0, err
	}
	return int64(t.Nanosecond() / 1000000), nil
}

func loadDuration(v *types.DurationValue) time.Duration {
	return time.Duration(v.Nanosecond)
}

// DurationGetHours returns the duration in hours.
func DurationGetHours(val *types.DurationValue) int64 {
	return int64(loadDuration(val).Hours())
}

// DurationGetMinutes returns duration in minutes.
func DurationGetMinutes(val *types.DurationValue) int64 {
	return int64(loadDuration(val).Minutes())
}

// DurationGetSeconds returns duration in seconds.
func DurationGetSeconds(val *types.DurationValue) int64 {
	return int64(loadDuration(val).Seconds())
}

// DurationGetMilliseconds returns duration in milliseconds.
func DurationGetMilliseconds(val *types.DurationValue) int64 {
	return int64(loadDuration(val).Milliseconds())
}

func GetHours(v *types.TimestampValue, tz string) (int64, error) {
	return TimestampGetHours(v, tz)
}
func GetMinutes(v *types.TimestampValue, tz string) (int64, error) {
	return TimestampGetMinutes(v, tz)
}
func GetSeconds(v *types.TimestampValue, tz string) (int64, error) {
	return TimestampGetSeconds(v, tz)
}
