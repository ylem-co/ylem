package runner

import (
	"context"
	"encoding/json"
	hevaluate "ylem_taskrunner/helpers/evaluate"
	"ylem_taskrunner/helpers/kafka"
	"ylem_taskrunner/services/evaluate"

	messaging "github.com/ylem-co/shared-messaging"
)

func CheckConditionTaskRunner(t *messaging.CheckConditionTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.TaskType = messaging.TaskTypeCondition
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		var i interface{}
		err := json.Unmarshal(t.Input, &i)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskCheckConditionMessageName, err, tr)

			return tr
		}

		switch in := i.(type) {
		case float64:
			i = map[string]float64{
				"value": in,
			}
		case string:
			i = map[string]string{
				"value": in,
			}
		}

		ctx := context.WithValue(context.Background(), "ctx", hevaluate.Context{ //nolint:all
			TaskInput:    i,
			EnvVars:      t.Meta.EnvVars,
			PipelineUuid: t.PipelineUuid,
		})
		result, err := evaluate.ConditionWithContext(ctx, t.Expression, i)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskCheckConditionMessageName, err, tr)

			return tr
		}

		response := &messaging.ConditionResult{
			Result:        result,
			OriginalInput: t.Input,
		}

		tr.Output, err = json.Marshal(response)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskCheckConditionMessageName, err, tr)

			return tr
		}

		return tr
	})
}
