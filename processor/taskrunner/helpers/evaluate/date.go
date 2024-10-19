package evaluate

import (
	"fmt"
	"reflect"
	"github.com/PaesslerAG/gval"
)

var dates = gval.NewLanguage(
	gval.InfixOperator(">", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc(">")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.After(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) > (%T)", a, b)
	}),
	gval.InfixOperator(">=", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc(">=")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.After(*date2) || date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) >= (%T)", a, b)
	}),
	gval.InfixOperator(">==", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc(">==")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.After(*date2) || date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) >== (%T)", a, b)
	}),
	gval.InfixOperator(">==", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc(">==")
		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.After(*date2) || date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) >== (%T)", a, b)
	}),
	gval.InfixOperator("<", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("<")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.Before(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) < (%T)", a, b)
	}),
	gval.InfixOperator("<=", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("<=")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.Before(*date2) || date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) <= (%T)", a, b)
	}),
	gval.InfixOperator("<==", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("<==")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.Before(*date2) || date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) <== (%T)", a, b)
	}),
	gval.InfixOperator("<==", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("<==")
		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.Before(*date2) || date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) <== (%T)", a, b)
	}),
	gval.InfixOperator("==", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("==")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) == (%T)", a, b)
	}),
	gval.InfixOperator("===", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("===")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) === (%T)", a, b)
	}),
	gval.InfixOperator("===", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("===")
		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) === (%T)", a, b)
	}),
	gval.InfixOperator("!=", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("!=")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return !date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) != (%T)", a, b)
	}),
	gval.InfixOperator("!==", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("!==")

		if !reflect.ValueOf(a).IsValid() {
			return false, nil
		}

		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return false, nil
		}

		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return !date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) !== (%T)", a, b)
	}),
	gval.InfixOperator("!==", func(a, b interface{}) (interface{}, error) {
		defer recoverGvalFunc("!==")
		date1, err1 := date(a)
		date2, err2 := date(b)

		if err1 == nil && err2 == nil {
			return !date1.Equal(*date2), nil
		}

		return nil, fmt.Errorf("unexpected operands types (%T) !== (%T)", a, b)
	}),
)
