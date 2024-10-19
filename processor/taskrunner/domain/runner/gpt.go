package runner

import (
	"fmt"
	"ylem_taskrunner/services/openai"

	messaging "github.com/ylem-co/shared-messaging"
)

func GptTaskRunner(t *messaging.CallOpenapiGptTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.TaskType = messaging.TaskTypeOpenapiGpt
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		const MaxJSONPayloadKb = 16
		if len(t.Input) > MaxJSONPayloadKb*1024 {
			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorCallOpenapiGptTaskFailure,
					Severity: messaging.ErrorSeverityError,
					Message: fmt.Sprintf(
						"Too large JSON. Maximum size is %dKb. Please use either a filter or transformer",
						MaxJSONPayloadKb,
					),
				},
			}

			return tr
		}

		inst := openai.Instance()
		resp, err := inst.CompleteText(openai.Completion{
			JSON:       string(t.Input),
			UserPrompt: t.Prompt,
		})

		if err != nil {
			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorCallOpenapiGptTaskFailure,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			return tr
		}

		tr.Output = []byte(resp)

		return tr
	})
}
