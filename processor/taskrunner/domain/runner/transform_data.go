package runner

import (
	"fmt"
	"encoding/json"
	"ylem_taskrunner/helpers/kafka"
	"ylem_taskrunner/services/transformers"

	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
)

func TransformDataTaskRunner(t *messaging.TransformDataTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.TaskType = messaging.TaskTypeTransformer
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		var (
			inputValue interface{}
			value      []byte
			err        error
		)

		inputValue, err = kafka.DecodeKafkaTaskValue(t.Task, messaging.TaskTransformDataMessageName, tr)
		if err != nil {
			return tr
		}

		transformerFound := true
		switch in := inputValue.(type) {
		// -------------------------------------------------------------------
		case string:
			switch t.Type {
			case transformers.TransformerTypeStrSplit:
				strings := transformers.SplitString(in, t.Delimiter)
				value, err = json.Marshal(strings)

				if err != nil {
					kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)

					return tr
				}
			case transformers.TransformerTypeCastTo:
				if t.CastToType == transformers.TransformerTypeCastToInteger {
					number, err := transformers.CastStringToInteger(in)

					if err != nil {
						log.Info(err.Error())
						kafka.HandleBadRequestError(
							t.TaskUuid,
							messaging.TaskTransformDataMessageName,
							fmt.Errorf(`sorry, such conversion is not supported. Can't cast "%s" to integer`, in),
							tr,
						)

						return tr
					}

					value, err = json.Marshal(number)

					if err != nil {
						kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)

						return tr
					}
				} else if t.CastToType == transformers.TransformerTypeCastToString {
					value, err = json.Marshal(in)

					if err != nil {
						kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)

						return tr
					}
				} else {
					transformerFound = false
				}
			default:
				transformerFound = false
			}
			// -------------------------------------------------------------------
		case float64:
			switch t.Type {
			case transformers.TransformerTypeCastTo:
				if t.CastToType == transformers.TransformerTypeCastToInteger {
					number := transformers.CastFloatToInteger(in)
					value, err = json.Marshal(number)

					if err != nil {
						kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)

						return tr
					}
				} else if t.CastToType == transformers.TransformerTypeCastToString {
					number := transformers.CastToStringType(in)
					value, err = json.Marshal(number)

					if err != nil {
						kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)

						return tr
					}
				} else {
					transformerFound = false
				}
			default:
				transformerFound = false
			}
			// -------------------------------------------------------------------
		case []interface{}:
			switch t.Type {
			case transformers.TransformerTypeExtractFromJSON:
				tr.Output, err = runExtractJsonTransformer(t.Input, t, tr)

				if err != nil {
					kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)

					return tr
				}

				return tr
			case transformers.TransformerTypeEncode:
				if t.EncodeFormat == transformers.TransformerTypeEncodeToXML {
					tr.Output, err = transformers.EncodeToXml(t.Input)
				} else if t.EncodeFormat == transformers.TransformerTypeEncodeToCSV {
					tr.Output, err = transformers.EncodeToCsv(t.Input, t.Delimiter, t.Meta.SqlQueryColumnOrder)
				} else {
					kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, fmt.Errorf("unknown encode type %s", t.EncodeFormat), tr)

					return tr
				}

				if err != nil {
					kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)
				}

				return tr
			default:
				transformerFound = false
			}
			// -------------------------------------------------------------------
		case map[string]interface{}:
			switch t.Type {
			case transformers.TransformerTypeExtractFromJSON:
				tr.Output, err = runExtractJsonTransformer(t.Input, t, tr)

				if err != nil {
					kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)

					return tr
				}

				return tr
			default:
				transformerFound = false
			}
			// -------------------------------------------------------------------
		default:
			transformerFound = false
		}

		if !transformerFound {
			log.Errorf(
				`could not execute task "%s"" with uuid "%s": %v`,
				messaging.TaskTransformDataMessageName,
				t.TaskUuid,
				fmt.Errorf("could not find a transformer"),
			)

			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorBadRequest,
					Severity: messaging.ErrorSeverityWarning,
					Message:  "sorry, such conversion is not supported",
				},
			}

			return tr
		}

		tr.Output = value

		return tr
	})
}

func runExtractJsonTransformer(value []byte, t *messaging.TransformDataTask, tr *messaging.TaskRunResult) ([]byte, error) {
	result := transformers.ExtractFromJsonWithJsonQuery(value, t.JsonQueryExpression)
	rawValue := result.Value()
	newValue, err := json.Marshal(rawValue)

	if err != nil {
		kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskTransformDataMessageName, err, tr)

		return nil, err
	}

	return newValue, nil
}
