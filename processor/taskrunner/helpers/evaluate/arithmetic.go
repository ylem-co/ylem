package evaluate

import (
	"github.com/PaesslerAG/gval"
	"github.com/shopspring/decimal"
)

var arithmetic = gval.NewLanguage(
	gval.InfixNumberOperator("+", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("+")

		return add(a, b), nil
	}),
	gval.InfixNumberOperator("-", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("-")

		return sub(a, b), nil
	}),
	gval.InfixNumberOperator("*", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("*")

		return mul(a, b), nil
	}),
	gval.InfixNumberOperator("/", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("/")

		return div(a, b), nil
	}),

	gval.InfixNumberOperator(">", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc(">")

		return gt(a, b), nil
	}),
	gval.InfixNumberOperator(">=", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc(">=")

		return gte(a, b), nil
	}),
	gval.InfixNumberOperator(">==", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc(">==")

		return gte(a, b), nil
	}),
	gval.InfixNumberOperator("<", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("<")

		return lt(a, b), nil
	}),
	gval.InfixNumberOperator("<=", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("<=")

		return lte(a, b), nil
	}),
	gval.InfixNumberOperator("<==", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("<==")

		return lte(a, b), nil
	}),
	gval.InfixNumberOperator("==", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("==")

		return eq(a, b), nil
	}),
	gval.InfixNumberOperator("===", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("===")

		return eq(a, b), nil
	}),
	gval.InfixNumberOperator("!=", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("!=")

		return !eq(a, b), nil
	}),
	gval.InfixNumberOperator("!==", func(a, b float64) (interface{}, error) {
		defer recoverGvalFunc("!==")

		return !eq(a, b), nil
	}),
)

func sub(a, b float64) float64 {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	diff := dA.Sub(dB)
	f64, _ := diff.Float64()

	return f64
}

func add(a, b float64) float64 {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	diff := dA.Add(dB)
	f64, _ := diff.Float64()

	return f64
}

func mul(a, b float64) float64 {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	diff := dA.Mul(dB)
	f64, _ := diff.Float64()

	return f64
}

func div(a, b float64) float64 {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	diff := dA.Div(dB)
	f64, _ := diff.Float64()

	return f64
}

func gt(a, b float64) bool {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	return dA.GreaterThan(dB)
}

func gte(a, b float64) bool {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	return dA.GreaterThanOrEqual(dB)
}

func lt(a, b float64) bool {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	return dA.LessThan(dB)
}

func lte(a, b float64) bool {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	return dA.LessThanOrEqual(dB)
}

func eq(a, b float64) bool {
	dA := decimal.NewFromFloat(a)
	dB := decimal.NewFromFloat(b)

	return dA.Equal(dB)
}
