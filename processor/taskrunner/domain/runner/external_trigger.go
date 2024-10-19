package runner

import (
	messaging "github.com/ylem-co/shared-messaging"
)

func ExternalTriggerTaskRunner(t *messaging.ExternalTriggerTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.TaskType = messaging.TaskTypeExternalTrigger
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta
		tr.IsSuccessful = true

		tr.Output = t.Input

		return tr
	})
}
