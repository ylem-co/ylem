package runner

import (
	"context"
	"fmt"
	"ylem_taskrunner/helpers"
	"ylem_taskrunner/helpers/kafka"
	"ylem_taskrunner/services/aws"
	"ylem_taskrunner/services/google/sheets"
	"ylem_taskrunner/services/hubspot"
	"ylem_taskrunner/services/incidentio"
	"ylem_taskrunner/services/jenkins"
	"ylem_taskrunner/services/jira"
	"ylem_taskrunner/services/opsgenie"
	"ylem_taskrunner/services/salesforce"
	"ylem_taskrunner/services/slack"
	"ylem_taskrunner/services/tableau"
	"ylem_taskrunner/services/templater"
	"ylem_taskrunner/services/twilio"

	hubspotclient "github.com/ylem-co/hubspot-client"
	salesforceclient "github.com/ylem-co/salesforce-client"
	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
)

func SendNotificationTaskRunner(t *messaging.SendNotificationTask, ctx context.Context) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.TaskType = messaging.TaskTypeNotification
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		switch t.Type {
		case messaging.NotificationTypeSms, messaging.NotificationTypeEmail:
			if !t.IsConfirmed {
				code := map[string]uint{
					messaging.NotificationTypeSms:   messaging.ErrorSendNotificationTaskUnconfirmedSms,
					messaging.NotificationTypeEmail: messaging.ErrorSendNotificationTaskUnconfirmedEmail,
				}[t.Type]

				tr.IsSuccessful = false
				tr.Errors = []messaging.TaskRunError{
					{
						Code:     code,
						Severity: messaging.ErrorSeverityWarning,
						Message: fmt.Sprintf(
							`Destination "%s" is not confirmed. Please confirm your %s to be able to send notifications to it.`,
							t.Integration.Name,
							t.Integration.Type,
						),
					},
				}

				return tr
			}
		}

		var (
			inputMap interface{}
			err      error
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
		}

		parsedPayload, err := templater.ParseTemplate(t.Body, inputMap, t.Meta.EnvVars)

		if err != nil {
			log.Errorf(
				`could not execute task "%s"" with uuid "%s": %v`,
				messaging.TaskSendNotificationMessageName,
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

		switch t.Type {
		case messaging.NotificationTypeSms:
			err = twilio.SendSms(t.Integration.Value, parsedPayload)
		case messaging.NotificationTypeWhatsApp:
			parsedPayload, _ = parseWhatsAppPayload(t, inputMap, tr)
			err = twilio.SendWhatsAppMessage(t.WhatsAppConfiguration.ContentSid, t.Integration.Value, parsedPayload)
		case messaging.NotificationTypeEmail:
			file := createNotificationFile(t)
			_, err = aws.SendEmail(t.Integration.Value, t.TaskName, parsedPayload, file)
		case messaging.NotificationTypeSlack:
			err = slack.SendSlackMessage(
				t.SlackConfiguration.SlackChannelId,
				t.TaskName,
				parsedPayload,
				t.Severity,
				t.SlackConfiguration.AccessToken,
			)
		case messaging.NotificationTypeJira:
			err = jira.CreateTask(
				ctx,
				jira.Issue{
					ProjectKey:  t.Integration.Value,
					IssueType:   t.JiraConfiguration.IssueType,
					Summary:     t.TaskName,
					Description: parsedPayload,
				},
				jira.Authentication{
					CloudId:              t.JiraConfiguration.Url,
					EncryptedDataKey:     t.JiraConfiguration.DataKey,
					EncryptedAccessToken: t.JiraConfiguration.AccessToken,
				},
			)
		case messaging.NotificationTypeHubspot:
			err = hubspot.CreateTicket(
				ctx,
				hubspotclient.CreateTicketRequest{
					Properties: hubspotclient.CreateTicketRequestProperties{
						Pipeline:         t.Integration.Value,
						PipelineStage:    t.HubspotConfiguration.PipelineStageCode,
						HsTicketPriority: t.Severity,
						HsOwnerId:        t.HubspotConfiguration.OwnerCode,
						Subject:          t.TaskName,
						Content:          parsedPayload,
					},
				},
				hubspot.Authentication{
					EncryptedDataKey:     t.HubspotConfiguration.DataKey,
					EncryptedAccessToken: t.HubspotConfiguration.AccessToken,
				},
			)
		case messaging.NotificationTypeSalesforce:
			err = salesforce.CreateCase(
				ctx,
				salesforceclient.CreateCaseRequest{
					Subject:     t.TaskName,
					Description: parsedPayload,
					Priority:    t.Severity,
				},
				salesforce.Authentication{
					EncryptedDataKey:     t.SalesforceConfiguration.DataKey,
					EncryptedAccessToken: t.SalesforceConfiguration.AccessToken,
				},
				t.SalesforceConfiguration.Domain,
			)
		case messaging.NotificationTypeIncidentIo:
			err = incidentio.DecryptKeyAndCreateIncident(
				ctx,
				t.IncidentIoConfiguration.DataKey,
				t.IncidentIoConfiguration.ApiKey,
				incidentio.Incident{
					IdempotencyKey: t.PipelineRunUuid.String(),
					Mode:           t.IncidentIoConfiguration.Mode,
					Name:           t.TaskName,
					SeverityId:     t.Severity,
					Status:         "triage",
					Summary:        parsedPayload,
					Visibility:     t.IncidentIoConfiguration.Visibility,
				},
			)
		case messaging.NotificationTypeOpsgenie:
			err = opsgenie.DecryptKeyAndCreateAlert(
				ctx,
				t.OpsgenieConfiguration.DataKey,
				t.OpsgenieConfiguration.ApiKey,
				opsgenie.Alert{
					Message:     t.TaskName,
					Description: parsedPayload,
					Priority:    t.Severity,
				},
			)
		case messaging.NotificationTypeTableau:
			var data []map[string]interface{}
			data, err = inputToMap(inputMap)

			if err == nil {
				var (
					username string
					password string
				)
				username, err = helpers.DecryptData(ctx, t.TableauConfiguration.DataKey, t.TableauConfiguration.Username)
				if err != nil {
					break
				}
				password, err = helpers.DecryptData(ctx, t.TableauConfiguration.DataKey, t.TableauConfiguration.Password)
				if err != nil {
					break
				}
				var columns []tableau.Column
				columns, err = tableau.Columns(t.Meta.SqlQueryColumnOrder, data)
				if err != nil {
					break
				}
				err = tableau.NewClient().Insert(
					t.TableauConfiguration.Server,
					username,
					password,
					t.TableauConfiguration.Sitename,
					t.TableauConfiguration.ProjectName,
					t.TableauConfiguration.DatasourceName,
					t.TableauConfiguration.Mode,
					columns,
					tableau.Rows(t.Meta.SqlQueryColumnOrder, data),
				)

				if err != nil {
					break
				}
			}

		case messaging.NotificationTypeGoogleSheets:
			var data []map[string]interface{}
			data, err = inputToMap(inputMap)
			if err != nil {
				break
			}

			var credentials string
			gs := t.GoogleSheetsConfiguration
			credentials, err = helpers.DecryptData(ctx, gs.DataKey, gs.Credentials)
			if err != nil {
				break
			}

			var c sheets.Client
			c, err = sheets.NewClient(ctx, []byte(credentials))
			if err != nil {
				break
			}

			err = c.WriteData(
				gs.SpreadsheetId,
				gs.SheetId,
				gs.Mode,
				t.Meta.SqlQueryColumnOrder,
				data,
				gs.WriteHeader,
			)

			if err != nil {
				break
			}
		case messaging.NotificationTypeJenkins:
			var token string
			token, err = helpers.DecryptData(ctx, t.JenkinsConfiguration.DataKey, t.JenkinsConfiguration.Token)
			if err != nil {
				break
			}

			err = jenkins.Instance(ctx).RunBuild(t.JenkinsConfiguration.BaseUrl, t.Integration.Value, token)
			if err != nil {
				break
			}
		}

		if err != nil {
			log.Errorf(
				`could not execute task "%s"" with uuid "%s": %v`,
				messaging.TaskSendNotificationMessageName,
				t.TaskUuid,
				err,
			)

			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorSendNotificationTaskFailure,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			return tr
		}

		tr.IsSuccessful = true

		return tr
	})
}

func createNotificationFile(t *messaging.SendNotificationTask) *aws.EmailFile {
	var file *aws.EmailFile
	if t.AttachedFileName != "" {
		file = &aws.EmailFile{
			Content: t.Input,
			Name:    t.AttachedFileName,
		}
	}

	return file
}

func inputToMap(input interface{}) ([]map[string]interface{}, error) {
	data, ok := input.([]map[string]interface{})
	if !ok {
		data2, ok := input.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid payload format: %#v", input)
		}
		data = make([]map[string]interface{}, len(data2))
		for k, v := range data2 {
			data[k], ok = v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid payload format: %#v", v)
			}
		}
	}

	return data, nil
}

func parseWhatsAppPayload(t *messaging.SendNotificationTask, InputMap interface{}, tr *messaging.TaskRunResult) (string, bool) {
	parsedPayload, err := templater.ParseJsonTemplate(t.Body, InputMap, t.Meta.EnvVars)

	if err != nil {
		log.Errorf(
			`could not execute task "%s"" with uuid "%s": %v`,
			messaging.TaskSendNotificationMessageName,
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
