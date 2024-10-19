package runner

import (
	"context"
	"encoding/json"
	hevaluate "ylem_taskrunner/helpers/evaluate"
	"ylem_taskrunner/helpers/kafka"
	sevaluate "ylem_taskrunner/services/evaluate"

	messaging "github.com/ylem-co/shared-messaging"
)

func AggregateDataTaskRunner(t *messaging.AggregateDataTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.TaskType = messaging.TaskTypeAggregator
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		var i interface{}
		err := json.Unmarshal(t.Input, &i)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskAggregateDataMessageName, err, tr)

			return tr
		}

		ctx := context.WithValue(context.Background(), "ctx", hevaluate.Context{ //nolint:all
			TaskInput:    i,
			EnvVars:      t.Meta.EnvVars,
			PipelineUuid: t.PipelineUuid,
		})
		result, err := sevaluate.AggregateWithContext(ctx, t.Expression, i)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskAggregateDataMessageName, err, tr)

			return tr
		}

		output := []map[string]interface{}{
			{
				t.VariableName: result,
			},
		}
		tr.Output, err = json.Marshal(output)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskAggregateDataMessageName, err, tr)

			return tr
		}

		return tr
	})
}
