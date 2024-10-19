package helpers

import (
	"context"
	"regexp"
	"ylem_taskrunner/helpers/evaluate"

	"github.com/PaesslerAG/gval"
)

func EvaluateGValExpressionWithContext(ctx context.Context, expression string, data interface{}) (interface{}, error) {
	return evaluateWithContext(ctx, expression, data)
}

func evaluateWithContext(ctx context.Context, expression string, data interface{}) (interface{}, error) {
	rx, err := regexp.Compile("COUNT *\\( *\\* *\\)") //nolint:all
	if err != nil {
		return nil, err
	}
	replacedExpression := rx.ReplaceAllString(expression, "COUNT()")

	return gval.EvaluateWithContext(
		ctx,
		replacedExpression,
		data,
		evaluate.Language(),
	)
}
