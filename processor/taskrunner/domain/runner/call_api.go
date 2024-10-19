package runner

import (
	"ylem_taskrunner/helpers/kafka"
	"ylem_taskrunner/services/api"
	"ylem_taskrunner/services/templater"

	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
)

func CallApiTaskRunner(t *messaging.CallApiTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.TaskType = messaging.TaskTypeApiCall
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		var (
			parsedQueryString string
			inputMap          interface{}
			err               error
			file              *api.File
			parsedHeaders     map[string]string
			parsedPayload     string
		)

		if t.AttachedFileName == "" {
			input, err := kafka.DecodeKafkaTaskValue(t.Task, messaging.TaskSendNotificationMessageName, tr)
			switch in := input.(type) {
			case string, float64:
				inputMap = map[string]interface{}{
					"value": in,
				}
			default:
				inputMap = in
			}
			if err != nil {
				log.Error(err)
				return tr
			}

			parsedQueryString, err = templater.ParseTemplate(t.QueryString, inputMap, t.Meta.EnvVars)

			if err != nil {
				log.Errorf(
					`could not execute task "%s"" with uuid "%s": %v`,
					messaging.TaskCallApiMessageName,
					t.TaskUuid,
					err,
				)

				tr.IsSuccessful = false
				tr.Errors = []messaging.TaskRunError{
					{
						Code:     messaging.ErrorBadRequest,
						Severity: messaging.ErrorSeverityError,
						Message:  err.Error(),
					},
				}

				return tr
			}

			var ok bool
			parsedPayload, ok = parseApiPayload(t, inputMap, tr)
			if !ok {
				return tr
			}
		} else {
			file = &api.File{
				Content: t.Input,
				Name:    t.AttachedFileName,
			}

			parsedPayload = t.Payload
			parsedQueryString = t.QueryString
			inputMap = make(map[string]interface{}, 0)
		}

		url := t.Integration.Value + "?" + parsedQueryString
		parsedHeaders, ok := parseApiHeaders(t, inputMap, tr)
		if !ok {
			return tr
		}

		respBody, err := api.Call(
			url,
			parsedPayload,
			api.Config{
				Headers:               parsedHeaders,
				File:                  file,
				Method:                t.Integration.Method,
				AuthType:              t.Integration.AuthType,
				AuthBearerToken:       t.Integration.AuthBearerToken,
				AuthBasicUserName:     t.Integration.AuthBasicUserName,
				AuthBasicUserPassword: t.Integration.AuthBasicUserPassword,
				AuthHeaderName:        t.Integration.AuthHeaderName,
				AuthHeaderValue:       t.Integration.AuthHeaderValue,
				Severity:              t.Severity,
			},
		)

		if err != nil {
			log.Errorf(
				`could not execute task "%s"" with uuid "%s": %v`,
				messaging.TaskCallApiMessageName,
				t.TaskUuid,
				err,
			)

			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorCallApiTaskFailure,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			return tr
		}

		tr.IsSuccessful = true
		tr.Output = respBody

		return tr
	})
}

func parseApiHeaders(t *messaging.CallApiTask, InputMap interface{}, tr *messaging.TaskRunResult) (map[string]string, bool) {
	parsedHeaders := map[string]string{}
	for k, v := range t.Headers {
		parsedHeaderValue, err := templater.ParseTemplate(v, InputMap, t.Meta.EnvVars)
		parsedHeaders[k] = parsedHeaderValue

		if err != nil {
			log.Errorf(
				`could not execute task "%s"" with uuid "%s": %v`,
				messaging.TaskCallApiMessageName,
				t.TaskUuid,
				err,
			)

			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorBadRequest,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			return nil, false
		}
	}

	return parsedHeaders, true
}

func parseApiPayload(t *messaging.CallApiTask, InputMap interface{}, tr *messaging.TaskRunResult) (string, bool) {
	parsedPayload, err := templater.ParseJsonTemplate(t.Payload, InputMap, t.Meta.EnvVars)

	if err != nil {
		log.Errorf(
			`could not execute task "%s"" with uuid "%s": %v`,
			messaging.TaskCallApiMessageName,
			t.TaskUuid,
			err,
		)

		tr.IsSuccessful = false
		tr.Errors = []messaging.TaskRunError{
			{
				Code:     messaging.ErrorBadRequest,
				Severity: messaging.ErrorSeverityError,
				Message:  err.Error(),
			},
		}

		return "", false
	}

	return parsedPayload, true
}
