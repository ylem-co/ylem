package runner

import (
	"encoding/json"
	"ylem_taskrunner/helpers/kafka"

	messaging "github.com/ylem-co/shared-messaging"
)

func RunForEachTaskRunner(t *messaging.RunForEachTask) []*messaging.TaskRunResult {
	return runMeasuredMultiOutput(func() []*messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.TaskType = messaging.TaskTypeForEach
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		var outputs []interface{}
		trs := make([]*messaging.TaskRunResult, 0)
		err := json.Unmarshal(t.Input, &outputs)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskRunForEachMessageName, err, tr)
			trs = append(trs, tr)

			return trs
		}

		if len(outputs) == 0 {
			tr.IsSuccessful = true
			tr.IsFinalTask = true
			trs = append(trs, tr)

			return trs
		}

		for _, output := range outputs {
			newTr := *tr
			var newOutputs []interface{}
			newOutputs = append(newOutputs, output)
			newTr.Output, err = json.Marshal(newOutputs)
			if err != nil {
				kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskRunForEachMessageName, err, tr)

				return trs
			}

			trs = append(trs, &newTr)
		}

		return trs
	})
}
