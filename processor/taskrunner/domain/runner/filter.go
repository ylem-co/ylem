package runner

import (
	"encoding/json"
	"ylem_taskrunner/helpers/kafka"
	"ylem_taskrunner/services/transformers"

	messaging "github.com/ylem-co/shared-messaging"
)

func FilterTaskRunner(t *messaging.FilterTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.TaskType = messaging.TaskTypeFilter
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		result := transformers.ExtractFromJsonWithJsonQuery(t.Input, t.Expression)
		newValue, err := json.Marshal(result.Value())

		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskFilterMessageName, err, tr)

			return tr
		}

		tr.Output = newValue

		return tr
	})
}
