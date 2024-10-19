package runner

import (
	"context"
	"fmt"
	"ylem_taskrunner/helpers/kafka"
	"ylem_taskrunner/services/aws/kms"
	"ylem_taskrunner/services/sqlIntegrations"
	"ylem_taskrunner/services/templater"

	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

func RunQueryTaskRunner(t *messaging.RunQueryTask) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.TaskType = messaging.TaskTypeQuery
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta

		err := kms.DecryptSource(&t.Source, context.TODO())

		if err != nil {
			log.Errorf("failed to decrypt a source: %s", err.Error())
			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorRunQueryTaskFailure,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			return tr
		}

		connection, err := sqlIntegrations.CreateSQLIntegrationConnection(
			t.Source.Type,
			sqlIntegrations.DefaultSQLIntegrationConnectionConfiguration{
				Host:        string(t.Source.Host),
				Port:        uint16(t.Source.Port),
				User:        t.Source.User,
				Password:    string(t.Source.Password),
				Database:    t.Source.Database,
				ProjectId:   t.Source.ProjectId,
				Credentials: string(t.Source.Credentials),
				SslEnabled:  t.Source.SslEnabled,
				EsVersion:   &t.Source.EsVersion,
			},
		)

		if err != nil {
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

		if t.Source.ConnectionType == messaging.SQLIntegrationConnectionTypeSsh {
			sshConn, ok := connection.(sqlIntegrations.ViaSshConnection)
			if !ok {
				err = fmt.Errorf("%s connection doesn't support SSH", t.Source.Type)
			} else {
				err = sshConn.OpenSsh(string(t.Source.SshHost), uint16(t.Source.SshPort), t.Source.SshUser)
			}
		} else {
			err = connection.Open()
		}

		if err != nil {
			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorRunQueryTaskFailure,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			tr.Output, _ = t.Source.Uuid.MarshalBinary()

			return tr
		}
		defer connection.Close()

		var inputMap interface{}
		input, err := kafka.DecodeKafkaTaskValue(t.Task, messaging.TaskRunQueryMessageName, tr)
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

		query, err := templater.ParseTemplate(t.Query, inputMap, t.Meta.EnvVars)
		if err != nil {
			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorRunQueryTaskFailure,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			tr.Output, _ = t.Source.Uuid.MarshalBinary()

			return tr
		}

		var collectedData []byte
		var columnNames []string
		if sqlConnection, ok := connection.(sqlIntegrations.SQLDriverConnection); ok {
			collectedData, columnNames, err = sqlIntegrations.CollectDataFromSQLSQLIntegrationAsJSON(sqlConnection, query)
		} else if queryableConnection, ok := connection.(sqlIntegrations.QueryableConnection); ok {
			collectedData, columnNames, err = sqlIntegrations.CollectDataFromQueryableSQLIntegrationAsJSON(queryableConnection, query)
		} else {
			err = fmt.Errorf("%s connection doesn't support querying", t.Source.Type)
		}

		if err != nil {
			tr.IsSuccessful = false
			tr.Errors = []messaging.TaskRunError{
				{
					Code:     messaging.ErrorRunQueryTaskFailure,
					Severity: messaging.ErrorSeverityError,
					Message:  err.Error(),
				},
			}

			tr.Output, _ = t.Source.Uuid.MarshalBinary()

			return tr
		}

		tr.Meta.SqlQueryColumnOrder = columnNames
		tr.Output = collectedData

		return tr
	})
}
