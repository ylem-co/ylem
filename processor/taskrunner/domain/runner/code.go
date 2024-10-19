package runner

import (
	"ylem_taskrunner/services/gopyk"

	messaging "github.com/ylem-co/shared-messaging"
)

func CodeTaskRunner(t *messaging.ExecuteCodeTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.TaskType = messaging.TaskTypeCode
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		inst := gopyk.Instance()
		resp, err := inst.Evaluate(gopyk.Request{
			Code:  t.Code,
			Type:  t.Type,
			Input: string(t.Input),
		})

		if err != nil {
			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorExecuteCodeFailure,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			return tr
		}

		tr.Output = resp

		return tr
	})
}
