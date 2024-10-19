package runner

import (
	"encoding/json"
	"ylem_taskrunner/helpers/kafka"

	messaging "github.com/ylem-co/shared-messaging"
)

func RunPipelineTaskRunner(t *messaging.RunPipelineTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.TaskType = messaging.TaskTypeRunPipeline
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		result, err := json.Marshal(t.PipelineToRunUuid)
		if err != nil {
			kafka.HandleBadRequestError(t.TaskUuid, messaging.TaskRunPipelineMessageName, err, tr)

			return tr
		}

		tr.Output = result

		return tr
	})
}
