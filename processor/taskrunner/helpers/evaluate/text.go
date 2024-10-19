package evaluate

import (
	"github.com/PaesslerAG/gval"
)

var text = gval.NewLanguage(
	gval.InfixTextOperator("==", func(a, b string) (interface{}, error) {
		defer recoverGvalFunc("==")

		return a == b, nil
	}),
	gval.InfixTextOperator("===", func(a, b string) (interface{}, error) {
		defer recoverGvalFunc("===")

		return a == b, nil
	}),
	gval.InfixTextOperator("!=", func(a, b string) (interface{}, error) {
		defer recoverGvalFunc("!=")

		return a != b, nil
	}),
	gval.InfixTextOperator("!==", func(a, b string) (interface{}, error) {
		defer recoverGvalFunc("!==")

		return a != b, nil
	}),
)
