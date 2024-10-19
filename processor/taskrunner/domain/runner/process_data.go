package runner

import (
	"encoding/json"
	"ylem_taskrunner/helpers/kafka"

	"github.com/itchyny/gojq"
	messaging "github.com/ylem-co/shared-messaging"
)

func ProcessDataTaskRunner(t *messaging.ProcessDataTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.TaskType = messaging.TaskTypeProcessor
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		query, err := gojq.Parse(t.Expression)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskProcessDataMessageName, err, tr)

			return tr
		}

		code, err := gojq.Compile(query)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskProcessDataMessageName, err, tr)

			return tr
		}

		var inputValue interface{}
		err = json.Unmarshal(t.Input, &inputValue)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskProcessDataMessageName, err, tr)

			return tr
		}

		iter := code.Run(inputValue)
		var result interface{}
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskProcessDataMessageName, err, tr)

				return tr
			}

			result = v;
		}

		if result == nil {
			result = "{}"
		}

		newValue, err := json.Marshal(result)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskProcessDataMessageName, err, tr)

			return tr
		}

		tr.Output = newValue

		return tr
	})
}
