package runner

import (
	"fmt"
	"ylem_taskrunner/config"
	"ylem_taskrunner/services/anthropic"
	"ylem_taskrunner/services/openai"
	"ylem_taskrunner/services/gemini"

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

		var err error
		var resp string
		switch config.Cfg().AIProvider {
		case "openai":
			inst := openai.Instance()
			resp, err = inst.CompleteText(openai.Completion{
				JSON:       string(t.Input),
				UserPrompt: t.Prompt,
			})
		case "gemini":
			resp, err = gemini.Process(string(t.Input), t.Prompt)
		case "anthropic":
			inst := anthropic.Instance()
			resp, err = inst.CompleteText(anthropic.Completion{
				JSON:       string(t.Input),
				UserPrompt: t.Prompt,
			})
		default:
			err = fmt.Errorf("unsupported AI provider %q", config.Cfg().AIProvider)
		}

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
