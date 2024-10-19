package messaging

import (
	"context"
	"fmt"
	"reflect"
	"ylem_pipelines/app/task"
	"ylem_pipelines/services/ylem_integrations"

	messaging "github.com/ylem-co/shared-messaging"
	"github.com/google/uuid"
)

type NotificationTaskMessageFactory struct {
	ctx         context.Context
	integrationsClient ylem_integrations.Client
}

func (f *NotificationTaskMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.Notification)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.ApiCall{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	tm, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	msg := &messaging.SendNotificationTask{
		Task:             tm,
		Type:             impl.Type,
		Body:             impl.Body,
		Severity:         t.Severity,
		AttachedFileName: impl.AttachedFileName,
	}

	dUid, err := uuid.Parse(impl.DestinationUuid)
	if err != nil {
		return nil, err
	}

	switch impl.Type {
	case task.NotificationTypeEmail:
		d, err := f.integrationsClient.GetEmailIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.IsConfirmed = d.IsConfirmed

	case task.NotificationTypeSms:
		d, err := f.integrationsClient.GetSmsIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.IsConfirmed = d.IsConfirmed

	case task.NotificationTypeSlack:
		d, err := f.integrationsClient.GetSlackIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.SlackConfiguration.AccessToken = d.SlackAuthorization.AccessToken
		msg.SlackConfiguration.SlackChannelId = d.SlackChannelId

	case task.NotificationTypeJira:
		d, err := f.integrationsClient.GetJiraIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.JiraConfiguration.AccessToken = d.AccessToken
		msg.JiraConfiguration.DataKey = d.DataKey
		msg.JiraConfiguration.ProjectKey = d.Integration.Value
		msg.JiraConfiguration.IssueType = d.IssueType
		msg.JiraConfiguration.Url = d.CloudId

	case task.NotificationTypeIncidentIo:
		d, err := f.integrationsClient.GetIncidentIoIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.IncidentIoConfiguration.ApiKey = d.ApiKey
		msg.IncidentIoConfiguration.DataKey = d.DataKey
		msg.IncidentIoConfiguration.Mode = d.Mode
		msg.IncidentIoConfiguration.Visibility = d.Visibility

	case task.NotificationTypeOpsgenie:
		d, err := f.integrationsClient.GetOpsgenieIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.OpsgenieConfiguration.ApiKey = d.ApiKey
		msg.OpsgenieConfiguration.DataKey = d.DataKey

	case task.NotificationTypeTableau:
		d, err := f.integrationsClient.GetTableauIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.TableauConfiguration.DataKey = d.DataKey
		msg.TableauConfiguration.Server = d.Server
		msg.TableauConfiguration.Username = d.Username
		msg.TableauConfiguration.Password = d.Password
		msg.TableauConfiguration.Sitename = d.Sitename
		msg.TableauConfiguration.ProjectName = d.ProjectName
		msg.TableauConfiguration.DatasourceName = d.DatasourceName
		msg.TableauConfiguration.Mode = d.Mode

	case task.NotificationTypeHubspot:
		d, err := f.integrationsClient.GetHubspotIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.HubspotConfiguration.DataKey = d.DataKey
		msg.HubspotConfiguration.AccessToken = d.AccessToken
		msg.HubspotConfiguration.PipelineStageCode = d.PipelineStageCode
		msg.HubspotConfiguration.OwnerCode = d.OwnerCode

	case task.NotificationTypeGoogleSheets:
		d, err := f.integrationsClient.GetGoogleSheetsIntegration(dUid)

		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.GoogleSheetsConfiguration.DataKey = d.DataKey
		msg.GoogleSheetsConfiguration.SpreadsheetId = d.SpreadsheetId
		msg.GoogleSheetsConfiguration.SheetId = d.SheetId
		msg.GoogleSheetsConfiguration.Mode = d.Mode
		msg.GoogleSheetsConfiguration.Credentials = d.Credentials
		msg.GoogleSheetsConfiguration.WriteHeader = d.WriteHeader

	case task.NotificationTypeSalesforce:
		d, err := f.integrationsClient.GetSalesforceIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.SalesforceConfiguration.DataKey = d.DataKey
		msg.SalesforceConfiguration.AccessToken = d.AccessToken
		msg.SalesforceConfiguration.Domain = d.Domain

	case task.NotificationTypeJenkins:
		d, err := f.integrationsClient.GetJenkinsIntegration(dUid)
		if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
			return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
		} else if err != nil {
			return nil, err
		}

		err = f.setIntegration(msg, d.Integration)
		if err != nil {
			return nil, err
		}
		msg.JenkinsConfiguration.DataKey = d.DataKey
		msg.JenkinsConfiguration.Token = d.Token
		msg.JenkinsConfiguration.BaseUrl = d.BaseUrl

	default:
		return nil, NewErrorNonRepeatable(fmt.Sprintf("unknown notification type %s", impl.Type))
	}

	return messaging.NewEnvelope(msg), nil
}

func (f *NotificationTaskMessageFactory) setIntegration(msg *messaging.SendNotificationTask, dest ylem_integrations.Integration) error {
	var err error
	msg.Integration, err = createIntegrationMessage(dest)

	return err
}

func NewNotificationTaskMessageFactory(ctx context.Context) (*NotificationTaskMessageFactory, error) {
	ycl, err := ylem_integrations.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	f := &NotificationTaskMessageFactory{
		ctx:         ctx,
		integrationsClient: ycl,
	}

	return f, nil
}
