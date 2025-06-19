package functions

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/yywing/sl/ast"
	"github.com/yywing/sl/lib/types"
	"github.com/yywing/sl/native"
)

func init() {
	ast.AddFunction.Combine(AddFunction)
	ast.SubtractFunction.Combine(SubtractFunction)
	ast.LessFunction.Combine(LessFunction)
	ast.LessEqualsFunction.Combine(LessEqualsFunction)
	ast.GreaterFunction.Combine(GreaterFunction)
	ast.GreaterEqualsFunction.Combine(GreaterEqualsFunction)
	ast.IntFunction.Combine(IntFunction)
	ast.StringFunction.Combine(StringFunction)

	LibFunctions["now"] = ast.NewBaseFunction("now", native.MustNewNativeFunction("now", Now).Definitions())
	LibFunctions["getFullYear"] = ast.NewBaseFunction("getFullYear", native.MustNewNativeFunction("getFullYear", GetFullYear).WithDefaultArg("").Definitions())
	LibFunctions["getMonth"] = ast.NewBaseFunction("getMonth", native.MustNewNativeFunction("getMonth", GetMonth).WithDefaultArg("").Definitions())
	LibFunctions["getDayOfYear"] = ast.NewBaseFunction("getDayOfYear", native.MustNewNativeFunction("getDayOfYear", GetDayOfYear).WithDefaultArg("").Definitions())
	LibFunctions["getDate"] = ast.NewBaseFunction("getDate", native.MustNewNativeFunction("getDate", GetDayOfMonthOneBased).WithDefaultArg("").Definitions())
	LibFunctions["getDayOfMonth"] = ast.NewBaseFunction("getDayOfMonth", native.MustNewNativeFunction("getDayOfMonth", GetDayOfMonthZeroBased).WithDefaultArg("").Definitions())
	LibFunctions["getDayOfWeek"] = ast.NewBaseFunction("getDayOfWeek", native.MustNewNativeFunction("getDayOfWeek", GetDayOfWeek).WithDefaultArg("").Definitions())
	LibFunctions["getHours"] = GetHoursFunction
	LibFunctions["getMinutes"] = GetMinutesFunction
	LibFunctions["getSeconds"] = GetSecondsFunction
	LibFunctions["getMilliseconds"] = GetMillisecondsFunction
	LibFunctions["duration"] = DurationFunction
	LibFunctions["timestamp"] = TimestampFunction
}

const (
	FunctionGetHours        = "getHours"
	FunctionGetMinutes      = "getMinutes"
	FunctionGetSeconds      = "getSeconds"
	FunctionGetMilliseconds = "getMilliseconds"
	FunctionDuration        = "duration"
	FunctionTimestamp       = "timestamp"
	FunctionNow             = "now"
)

var (
	AddFunction = ast.NewBaseFunction(
		ast.Add,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.Add, []ast.ValueType{types.DurationType, types.DurationType}, types.DurationType),
				Call: func(args []ast.Value) (ast.Value, error) {
					x := args[0].(*types.DurationValue).Nanosecond
					y := args[1].(*types.DurationValue).Nanosecond
					if (y > 0 && x > math.MaxInt64-y) || (y < 0 && x < math.MinInt64-y) {
						return nil, fmt.Errorf("int overflow")
					}
					return types.NewDurationValue(x + y), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.Add, []ast.ValueType{types.DurationType, types.TimestampType}, types.TimestampType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t, err := loadTimestamp(args[1].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					d := loadDuration(args[0].(*types.DurationValue))
					result := t.Add(d)
					if result.Unix() < types.MinUnixTime || result.Unix() > types.MaxUnixTime {
						return nil, fmt.Errorf("timestamp overflow")
					}
					return exportTimestamp(&result), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.Add, []ast.ValueType{types.TimestampType, types.DurationType}, types.TimestampType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					d := loadDuration(args[1].(*types.DurationValue))
					result := t.Add(d)
					if result.Unix() < types.MinUnixTime || result.Unix() > types.MaxUnixTime {
						return nil, fmt.Errorf("timestamp overflow")
					}
					return exportTimestamp(&result), nil
				},
			},
		},
	)
	SubtractFunction = ast.NewBaseFunction(
		ast.Subtract,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.Subtract, []ast.ValueType{types.DurationType, types.DurationType}, types.DurationType),
				Call: func(args []ast.Value) (ast.Value, error) {
					x := args[0].(*types.DurationValue).Nanosecond
					y := args[1].(*types.DurationValue).Nanosecond
					if (y > 0 && x > math.MaxInt64-y) || (y < 0 && x < math.MinInt64-y) {
						return nil, fmt.Errorf("int overflow")
					}
					return types.NewDurationValue(x - y), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.Add, []ast.ValueType{types.TimestampType, types.DurationType}, types.TimestampType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					d := loadDuration(args[1].(*types.DurationValue))
					result := t.Add(-d)
					if result.Unix() < types.MinUnixTime || result.Unix() > types.MaxUnixTime {
						return nil, fmt.Errorf("timestamp overflow")
					}
					return exportTimestamp(&result), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.Subtract, []ast.ValueType{types.TimestampType, types.TimestampType}, types.DurationType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t1, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					t2, err := loadTimestamp(args[1].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}

					d := t1.Sub(*t2)
					if !t1.Add(-d).Equal(*t2) {
						return nil, fmt.Errorf("duration overflow")
					}
					return types.NewDurationValue(d.Nanoseconds()), nil
				},
			},
		},
	)
	LessFunction = ast.NewBaseFunction(
		ast.Less,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.Less, []ast.ValueType{types.DurationType, types.DurationType}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewBoolValue(args[0].(*types.DurationValue).Nanosecond < args[1].(*types.DurationValue).Nanosecond), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.Subtract, []ast.ValueType{types.TimestampType, types.TimestampType}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t1, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					t2, err := loadTimestamp(args[1].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewBoolValue(t1.Before(*t2)), nil
				},
			},
		},
	)
	LessEqualsFunction = ast.NewBaseFunction(
		ast.LessEquals,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.LessEquals, []ast.ValueType{types.DurationType, types.DurationType}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewBoolValue(args[0].(*types.DurationValue).Nanosecond <= args[1].(*types.DurationValue).Nanosecond), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.LessEquals, []ast.ValueType{types.TimestampType, types.TimestampType}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t1, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					t2, err := loadTimestamp(args[1].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewBoolValue(t1.Before(*t2) || t1.Equal(*t2)), nil
				},
			},
		},
	)
	GreaterFunction = ast.NewBaseFunction(
		ast.Greater,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.Greater, []ast.ValueType{types.DurationType, types.DurationType}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewBoolValue(args[0].(*types.DurationValue).Nanosecond > args[1].(*types.DurationValue).Nanosecond), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.Greater, []ast.ValueType{types.TimestampType, types.TimestampType}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t1, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					t2, err := loadTimestamp(args[1].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewBoolValue(t1.After(*t2)), nil
				},
			},
		},
	)
	GreaterEqualsFunction = ast.NewBaseFunction(
		ast.GreaterEquals,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.GreaterEquals, []ast.ValueType{types.DurationType, types.DurationType}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewBoolValue(args[0].(*types.DurationValue).Nanosecond >= args[1].(*types.DurationValue).Nanosecond), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.GreaterEquals, []ast.ValueType{types.TimestampType, types.TimestampType}, ast.BoolType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t1, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					t2, err := loadTimestamp(args[1].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewBoolValue(t1.After(*t2) || t1.Equal(*t2)), nil
				},
			},
		},
	)

	DurationFunction = ast.NewBaseFunction(
		FunctionDuration,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(FunctionDuration, []ast.ValueType{types.DurationType}, types.DurationType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return args[0], nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionDuration, []ast.ValueType{ast.IntType}, types.DurationType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return types.NewDurationValue(args[0].(*ast.IntValue).IntValue), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionDuration, []ast.ValueType{ast.StringType}, types.DurationType),
				Call: func(args []ast.Value) (ast.Value, error) {
					d, err := time.ParseDuration(args[0].(*ast.StringValue).StringValue)
					if err != nil {
						return nil, err
					}
					return types.NewDurationValue(int64(d.Nanoseconds())), nil
				},
			},
		},
	)

	TimestampFunction = ast.NewBaseFunction(
		FunctionTimestamp,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(FunctionTimestamp, []ast.ValueType{types.TimestampType}, types.TimestampType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return args[0], nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionTimestamp, []ast.ValueType{ast.IntType}, types.TimestampType),
				Call: func(args []ast.Value) (ast.Value, error) {
					i := args[0].(*ast.IntValue).IntValue
					// The maximum positive value that can be passed to time.Unix is math.MaxInt64 minus the number
					// of seconds between year 1 and year 1970. See comments on unixToInternal.
					if int64(i) < types.MinUnixTime || int64(i) > types.MaxUnixTime {
						return nil, fmt.Errorf("timestamp overflow")
					}
					t := time.Unix(int64(i), 0).In(time.UTC)
					return exportTimestamp(&t), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionTimestamp, []ast.ValueType{ast.StringType}, types.TimestampType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t, err := time.Parse(time.RFC3339, args[0].(*ast.StringValue).StringValue)
					if err != nil {
						return nil, err
					}
					if t.Unix() < types.MinUnixTime || t.Unix() > types.MaxUnixTime {
						return nil, fmt.Errorf("timestamp overflow")
					}
					return exportTimestamp(&t), nil
				},
			},
		},
	)

	IntFunction = ast.NewBaseFunction(
		ast.Int,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.Int, []ast.ValueType{types.DurationType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewIntValue(args[0].(*types.DurationValue).Nanosecond), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.Int, []ast.ValueType{types.TimestampType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(t.Unix()), nil
				},
			},
		},
	)

	StringFunction = ast.NewBaseFunction(
		ast.String,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(ast.Int, []ast.ValueType{types.DurationType}, ast.StringType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewStringValue(strconv.FormatFloat(time.Duration(args[0].(*types.DurationValue).Nanosecond).Seconds(), 'f', -1, 64) + "s"), nil
				},
			},
			{
				Type: *ast.NewFunctionType(ast.Int, []ast.ValueType{types.TimestampType}, ast.StringType),
				Call: func(args []ast.Value) (ast.Value, error) {
					t, err := loadTimestamp(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewStringValue(t.Format(time.RFC3339Nano)), nil
				},
			},
		},
	)

	GetHoursFunction = ast.NewBaseFunction(
		FunctionGetHours,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(FunctionGetHours, []ast.ValueType{types.DurationType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewIntValue(DurationGetHours(args[0].(*types.DurationValue))), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGetHours, []ast.ValueType{types.TimestampType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					result, err := TimestampGetHours(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(result), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGetHours, []ast.ValueType{types.TimestampType, ast.StringType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					result, err := TimestampGetHours(args[0].(*types.TimestampValue), args[1].(*ast.StringValue).StringValue)
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(result), nil
				},
			},
		},
	)

	GetMinutesFunction = ast.NewBaseFunction(
		FunctionGetMinutes,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(FunctionGetMinutes, []ast.ValueType{types.DurationType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewIntValue(DurationGetMinutes(args[0].(*types.DurationValue))), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGetMinutes, []ast.ValueType{types.TimestampType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					result, err := TimestampGetMinutes(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(result), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGetMinutes, []ast.ValueType{types.TimestampType, ast.StringType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					result, err := TimestampGetMinutes(args[0].(*types.TimestampValue), args[1].(*ast.StringValue).StringValue)
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(result), nil
				},
			},
		},
	)

	GetSecondsFunction = ast.NewBaseFunction(
		FunctionGetSeconds,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(FunctionGetSeconds, []ast.ValueType{types.DurationType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewIntValue(DurationGetSeconds(args[0].(*types.DurationValue))), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGetSeconds, []ast.ValueType{types.TimestampType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					result, err := TimestampGetSeconds(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(result), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGetSeconds, []ast.ValueType{types.TimestampType, ast.StringType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					result, err := TimestampGetSeconds(args[0].(*types.TimestampValue), args[1].(*ast.StringValue).StringValue)
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(result), nil
				},
			},
		},
	)

	GetMillisecondsFunction = ast.NewBaseFunction(
		FunctionGetMilliseconds,
		[]ast.Definition{
			{
				Type: *ast.NewFunctionType(FunctionGetMilliseconds, []ast.ValueType{types.DurationType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					return ast.NewIntValue(DurationGetMilliseconds(args[0].(*types.DurationValue))), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGetMilliseconds, []ast.ValueType{types.TimestampType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					result, err := TimestampGetMilliseconds(args[0].(*types.TimestampValue), "")
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(result), nil
				},
			},
			{
				Type: *ast.NewFunctionType(FunctionGetMilliseconds, []ast.ValueType{types.TimestampType, ast.StringType}, ast.IntType),
				Call: func(args []ast.Value) (ast.Value, error) {
					result, err := TimestampGetMilliseconds(args[0].(*types.TimestampValue), args[1].(*ast.StringValue).StringValue)
					if err != nil {
						return nil, err
					}
					return ast.NewIntValue(result), nil
				},
			},
		},
	)
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
	t := time.Unix(v.Sec, v.NSec).In(loc)
	return &t, nil
}

func exportTimestamp(t *time.Time) *types.TimestampValue {
	return types.NewTimestampValue(t.Unix(), t.Sub(time.Unix(t.Unix(), 0)).Nanoseconds(), t.Location().String())
}

func Now() *types.TimestampValue {
	t := time.Now()
	return exportTimestamp(&t)
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
