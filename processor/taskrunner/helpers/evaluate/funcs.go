package evaluate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	go_math "math"
	"strconv"
	"strings"
	"ylem_taskrunner/services/ylem_statistics"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func sum(in []interface{}) (interface{}, error) {
	decimals, err := interfaceToDecimal(in)
	if err != nil {
		return nil, err
	}

	sum, _ := decimal.Sum(decimal.Zero, decimals...).Float64()

	return sum, nil
}

func count(in []interface{}) interface{} {
	return len(in)
}

func avg(in []interface{}) (interface{}, error) {
	decimals, err := interfaceToDecimal(in)
	if err != nil {
		return nil, err
	}

	count := decimal.New(int64(len(in)), 0)
	sum := decimal.Sum(decimal.Zero, decimals...)

	avg, _ := sum.Div(count).Float64()

	return avg, nil
}

func min(in []interface{}) (interface{}, error) {
	return minmaxWrapper(in, func(slice []decimal.Decimal) decimal.Decimal {
		return decimal.Min(decimal.Max(decimal.Zero, slice...), slice...)
	})
}

func max(in []interface{}) (interface{}, error) {
	return minmaxWrapper(in, func(slice []decimal.Decimal) decimal.Decimal {
		return decimal.Max(decimal.Zero, slice...)
	})
}

func minmaxWrapper(in []interface{}, cb func(slice []decimal.Decimal) decimal.Decimal) (interface{}, error) {
	bFound := false
	for _, v := range in {
		if v == nil {
			continue
		}

		_, err := date(v)

		if err != nil {
			log.Debugf(`The string %s is not time. %s`, in[0], err.Error())

			break
		}

		bFound = true
		break
	}

	if !bFound {
		decimals, err := interfaceToDecimal(in)
		if err != nil {
			return nil, err
		}

		f, _ := cb(decimals).Float64()

		return f, nil
	}

	unixSlice := make([]int64, 0)
	unixMap := make(map[int64]string, 0)
	for _, v := range in {
		if v == nil {
			continue
		}

		t, err := date(v)
		if err != nil {
			return nil, fmt.Errorf(`string "%s" is not time. %s`, v, err.Error())
		}

		unix := t.UnixMilli()
		unixSlice = append(unixSlice, unix)
		unixMap[unix] = v.(string)
	}

	decimals, err := int64ToDecimal(unixSlice)
	if err != nil {
		return nil, err
	}

	selectedTime := cb(decimals).IntPart()

	return unixMap[selectedTime], err
}

func first(in []interface{}) (interface{}, error) {
	if len(in) == 0 {
		return nil, fmt.Errorf("no elements")
	}

	return in[0], nil
}

func date(in interface{}) (*time.Time, error) {
	t, ok := in.(time.Time)
	if ok {
		return &t, nil
	}

	s, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("date() expects exactly one string argument")
	}
	for _, format := range [...]string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.Kitchen,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",                         // RFC 3339
		"2006-01-02 15:04",                   // RFC 3339 with minutes
		"2006-01-02 15:04:05",                // RFC 3339 with seconds
		"2006-01-02 15:04:05-07:00",          // RFC 3339 with seconds and timezone
		"2006-01-02T15Z0700",                 // ISO8601 with hour
		"2006-01-02T15Z",                     // ISO8601 with hour without UTC offset
		"2006-01-02T15:04Z0700",              // ISO8601 with minutes
		"2006-01-02T15:04Z",                  // ISO8601 with minutes without UTC offset
		"2006-01-02T15:04:05Z0700",           // ISO8601 with seconds
		"2006-01-02T15:04:05Z",               // ISO8601 with seconds without UTC offset
		"2006-01-02T15:04:05.999999999Z0700", // ISO8601 with nanoseconds
		"2006-01-02T15:04:05.999999999Z",     // ISO8601 with nanoseconds without UTC offset
	} {
		ret, err := time.ParseInLocation(format, s, time.Local)
		if err == nil {
			return &ret, nil
		}
	}
	return nil, fmt.Errorf("date() could not parse %s", s)
}

func last(in []interface{}) (interface{}, error) {
	if len(in) == 0 {
		return nil, fmt.Errorf("no elements")
	}

	return in[len(in)-1], nil
}

func now() time.Time {
	return time.Now()
}

func neg(number float64) (interface{}, error) {
	dc := decimal.NewFromFloat(number)

	return dc.Neg(), nil
}

func sign(number float64) (interface{}, error) {
	dc := decimal.NewFromFloat(number)

	return dc.Sign(), nil
}

func stringFn(number float64) (interface{}, error) {
	dc := decimal.NewFromFloat(number)

	return dc.String(), nil
}

func intPart(number float64) (interface{}, error) {
	dc := decimal.NewFromFloat(number)

	return dc.IntPart(), nil
}

var funcs = gval.NewLanguage(
	gval.Function("AVG", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("AVG")
		_, ok := args[0].(noIdentifierPresented)
		if ok {
			return 0, nil
		}

		switch in := args[0].(type) { // means, they try to aggregate on single value
		case float64, int32:
			return in, nil
		}

		value, err := avg(args[0].([]interface{}))

		return value, err
	}),
	gval.Function("SUM", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("SUM")
		_, ok := args[0].(noIdentifierPresented)
		if ok {
			return 0, nil
		}

		switch in := args[0].(type) {
		case float64, int32:
			return in, nil
		}

		value, err := sum(args[0].([]interface{}))

		return value, err
	}),
	gval.Function("COUNT", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("COUNT")

		if len(args) == 0 {
			ctxValue, ok := ctx.Value("ctx").(Context)
			if ctxValue.TaskInput == nil {
				return 0, nil
			}

			if !ok {
				return nil, fmt.Errorf("invalid context")
			}

			switch in := ctxValue.TaskInput.(type) {
			case []interface{}:
				return len(in), nil
			case map[string]interface{}, float64, string, int32:
				return 1, nil
			default:
				return nil, fmt.Errorf("unknown task input")
			}
		}

		_, ok := args[0].(noIdentifierPresented)
		if ok {
			return 0, nil
		}

		switch in := args[0].(type) {
		case float64, int32:
			return 1, nil
		case string:
			return 1, nil
		default:
			_ = in
		}

		value := count(args[0].([]interface{}))

		return value, nil
	}),
	gval.Function("MIN", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("MIN")
		_, ok := args[0].(noIdentifierPresented)
		if ok {
			return 0, nil
		}

		switch in := args[0].(type) {
		case float64, int32:
			return in, nil
		case string:
			_, err := date(in)
			if err != nil {
				errf := fmt.Errorf(`the string %s is not time. %s`, args[0], err.Error())
				log.Debug(errf.Error())

				return nil, errf
			}

			return in, nil
		}

		value, err := min(args[0].([]interface{}))

		return value, err
	}),
	gval.Function("MAX", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("MAX")
		_, ok := args[0].(noIdentifierPresented)
		if ok {
			return 0, nil
		}

		switch in := args[0].(type) {
		case float64, int32:
			return in, nil
		case string:
			_, err := date(in)
			if err != nil {
				errf := fmt.Errorf(`the string %s is not time. %s`, args[0], err.Error())
				log.Debug(errf.Error())

				return nil, errf
			}

			return in, nil
		}

		value, err := max(args[0].([]interface{}))

		return value, err
	}),
	gval.Function("FIRST", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("FIRST")
		_, ok := args[0].(noIdentifierPresented)
		if ok {
			return nil, nil
		}

		switch in := args[0].(type) {
		case float64, int32:
			return in, nil
		}

		value, err := first(args[0].([]interface{}))

		return value, err
	}),
	gval.Function("LAST", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("LAST")
		_, ok := args[0].(noIdentifierPresented)
		if ok {
			return nil, nil
		}

		switch in := args[0].(type) {
		case float64, int32:
			return in, nil
		}

		value, err := last(args[0].([]interface{}))

		return value, err
	}),
	gval.Function("NOW", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("NOW")
		return now().String(), nil
	}),
	gval.Function("INPUT", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("INPUT")

		if len(args) > 0 {
			return nil, fmt.Errorf("INPUT doesn't take any arguments")
		}

		ctxValue, ok := ctx.Value("ctx").(Context)
		if !ok {
			return nil, fmt.Errorf("invalid context")
		}

		j, err := json.Marshal(ctxValue.TaskInput)
		if err != nil {
			return nil, err
		}

		return string(j), nil
	}),
	gval.Function("ROUND", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("ROUND")

		if len(args) != 3 {
			return nil, fmt.Errorf(`ROUND accepts 3 arguments: number, precision (e.g., 2), rounding strategy ("floor" or "ceil")`)
		}

		number, ok := args[0].(float64)
		if !ok {
			return nil, fmt.Errorf(`ROUND the first argument is not a number`)
		}

		precision, ok := args[1].(float64)
		if !ok {
			return nil, fmt.Errorf(`ROUND the second argument is not a number`)
		}

		strategy, ok := args[2].(string)
		if !ok {
			return nil, fmt.Errorf(`ROUND the third argument should be string`)
		}

		strategy = strings.ToLower(strategy)
		if strategy != "ceil" && strategy != "floor" {
			return nil, fmt.Errorf(`ROUND the third argument should be either "ceil" or "floor"`)
		}

		factor := go_math.Pow10(int(precision))

		if strategy == "ceil" {
			return go_math.Ceil(number*factor) / factor, nil
		}

		return go_math.Floor(number*factor) / factor, nil
	}),
	gval.Function("ABS", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("ABS")

		if len(args) != 1 {
			return nil, fmt.Errorf(`ABS accepts 1 argument: number`)
		}

		number, ok := args[0].(float64)
		if !ok {
			return nil, fmt.Errorf(`ABS the argument is not a number`)
		}

		return go_math.Abs(number), nil
	}),
	gval.Function("NEG", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("NEG")

		if len(args) != 1 {
			return nil, fmt.Errorf(`NEG accepts 1 argument: number`)
		}

		number, ok := args[0].(float64)
		if !ok {
			return nil, fmt.Errorf(`NEG the argument is not a number`)
		}

		return neg(number)
	}),
	gval.Function("STRING", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("STRING")

		if len(args) != 1 {
			return nil, fmt.Errorf(`STRING accepts 1 argument: number`)
		}

		number, ok := args[0].(float64)
		if !ok {
			return nil, fmt.Errorf(`STRING the argument is not a number`)
		}

		return stringFn(number)
	}),
	gval.Function("INT", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("INT")

		if len(args) != 1 {
			return nil, fmt.Errorf(`INT accepts 1 argument: number`)
		}

		number, ok := args[0].(float64)
		if !ok {
			return nil, fmt.Errorf(`INT the argument is not a number`)
		}

		return intPart(number)
	}),
	gval.Function("SIGN", func(args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("SIGN")

		if len(args) != 1 {
			return nil, fmt.Errorf(`SIGN accepts 1 argument: number`)
		}

		number, ok := args[0].(float64)
		if !ok {
			return nil, fmt.Errorf(`SIGN the argument is not a number`)
		}

		return sign(number)
	}),

	gval.Function("METRIC_AVG", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("METRIC_AVG")
		if len(args) != 2 {
			return 0, errors.New("METRIC_AVG expects exactly two arguments, e.g. METRIC_AVG(\"day\", 7)")
		}

		ctxValue, ok := ctx.Value("ctx").(Context)
		if !ok {
			return 0, fmt.Errorf("invalid context")
		}

		period, ok := args[0].(string)
		if !ok {
			return 0, errors.New("first argument must be a string")
		}

		periodCount, ok := args[1].(float64)
		if !ok {
			return 0, errors.New("second argument must be an integer")
		}

		client := ylem_statistics.NewClient()

		val, err := client.GetAverageMetricValue(ctxValue.PipelineUuid, period, int(periodCount))

		log.Tracef("METRIC_AVG(%s, %f) = %f", period, periodCount, val)

		return val, err
	}),

	gval.Function("METRIC_MEDIAN", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("METRIC_MEDIAN")
		if len(args) != 2 {
			return 0, errors.New("METRIC_MEDIAN expects exactly two arguments, e.g. METRIC_MEDIAN(\"day\", 7)")
		}

		ctxValue, ok := ctx.Value("ctx").(Context)
		if !ok {
			return 0, fmt.Errorf("invalid context")
		}

		period, ok := args[0].(string)
		if !ok {
			return 0, errors.New("first argument must be a string")
		}

		periodCount, ok := args[1].(float64)
		if !ok {
			return 0, errors.New("second argument must be an integer")
		}

		client := ylem_statistics.NewClient()

		val, err := client.GetMetricValueQuantile(ctxValue.PipelineUuid, 0.5, period, int(periodCount))

		log.Tracef("METRIC_MEDIAN(%s, %f) = %f", period, periodCount, val)

		return val, err
	}),

	gval.Function("METRIC_QUANTILE", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		defer recoverGvalFunc("METRIC_QUANTILE")
		if len(args) != 3 {
			return 0, errors.New("METRIC_QUANTILE expects exactly three arguments, e.g. METRIC_QUANTILE(0.5, \"day\", 7)")
		}

		ctxValue, ok := ctx.Value("ctx").(Context)
		if !ok {
			return 0, fmt.Errorf("invalid context")
		}

		level, ok := args[0].(float64)
		if !ok {
			return 0, errors.New("first argument must be a decimal")
		}

		period, ok := args[1].(string)
		if !ok {
			return 0, errors.New("second argument must be a string")
		}

		periodCount, ok := args[2].(float64)
		if !ok {
			return 0, errors.New("third argument must be an integer")
		}

		client := ylem_statistics.NewClient()

		val, err := client.GetMetricValueQuantile(ctxValue.PipelineUuid, level, period, int(periodCount))

		log.Tracef("METRIC_QUANTILE(%f, %s, %f) = %f", level, period, periodCount, val)

		return val, err
	}),
)

func interfaceToDecimal(in []interface{}) ([]decimal.Decimal, error) {
	decimals := make([]decimal.Decimal, 0)
	for _, v := range in {
		switch v := v.(type) {
		case float64:
			decimals = append(decimals, decimal.NewFromFloat(v))
		case int64:
			decimals = append(decimals, decimal.New(v, 0))
		case string:
			floatval, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return decimals, err
			}
			decimals = append(decimals, decimal.NewFromFloat(floatval))
		case nil:
			continue
		default:
			return nil, fmt.Errorf("%v expected to be float64 or int64, got %T", v, v)
		}
	}

	return decimals, nil
}

func int64ToDecimal(in []int64) ([]decimal.Decimal, error) {
	decimals := make([]decimal.Decimal, len(in))
	for k, v := range in {
		decimals[k] = decimal.New(v, 0)
	}

	return decimals, nil
}
