package evaluate

import (
	"context"
	"fmt"
	"reflect"
	"ylem_taskrunner/helpers"

	log "github.com/sirupsen/logrus"
)

func ConditionWithContext(ctx context.Context, expression string, data interface{}) (bool, error) {
	transformedData := transform(data)
	result, err := helpers.EvaluateGValExpressionWithContext(ctx, expression, transformedData)

	if err != nil {
		return false, err
	}

	switch in := result.(type) {
	case bool:
		return in, nil
	default:
		return false, fmt.Errorf("non boolean evaluation")
	}
}

func AggregateWithContext(ctx context.Context, expression string, data interface{}) (interface{}, error) {
	result, err := helpers.EvaluateGValExpressionWithContext(ctx, expression, transform(data))

	if err != nil {
		return false, err
	}

	return result, nil
}

func transform(in interface{}) interface{} {
	transformed := map[string][]interface{}{}

	switch in := in.(type) {
	case []interface{}:
		for _, object := range in {
			objectInterface, _ := object.(map[string]interface{})

			if len(in) == 1 {
				return objectInterface
			}

			for k, v := range objectInterface {
				transformed[k] = append(transformed[k], v)
			}
		}
	case map[string]interface{}:
		return in
	case map[string]float64, map[string]string:
		return in
	case nil:
		return map[string]interface{}{}
	default:
		if in == nil {
			log.Debugf(`Could not transform a type of nil`)
		} else {
			log.Debugf(`Could not transform a type of "%s"`, reflect.TypeOf(in).String())
		}
	}

	return transformed
}
