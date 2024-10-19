package server

import (
	"context"
	"net"
	"net/http"
	"ylem_integrations/api"
	"ylem_integrations/config"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Listen string
}

func (s *Server) Run(ctx context.Context) error {
	log.Info("Starting server listening on " + s.Listen)

	rtr := mux.NewRouter()

	rtr.HandleFunc("/api", api.CreateApiIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/api/{uuid}", api.UpdateApiIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/api/{uuid}", api.GetApiIntegration).Methods(http.MethodGet)

	rtr.HandleFunc("/sms", api.CreateSmsIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/sms/{uuid}", api.UpdateSmsIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/sms/{uuid}", api.GetSmsIntegration).Methods(http.MethodGet)
	rtr.HandleFunc("/sms/{uuid}/confirm", api.ConfirmSmsIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/sms/{uuid}/resend", api.ResendConfirmationSms).Methods(http.MethodPost)

	rtr.HandleFunc("/slack/authorize", api.AuthorizeSlack).
		Queries("code", "{code}").
		Queries("state", "{state}").
		Methods(http.MethodGet)
	rtr.HandleFunc("/slack/authorization/{uuid}", api.UpdateSlackAuthorization).Methods(http.MethodPost)
	rtr.HandleFunc("/slack/authorization/{uuid}", api.GetSlackAuthorization).Methods(http.MethodGet)
	rtr.HandleFunc("/slack/authorization/{uuid}/reauthorize", api.ReauthorizeSlackAuthorization).Methods(http.MethodPost)
	rtr.HandleFunc("/slack", api.CreateSlackIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/slack/{uuid}", api.UpdateSlackIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/slack/{uuid}", api.GetSlackIntegration).Methods(http.MethodGet)

	rtr.HandleFunc("/email", api.CreateEmailIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/email/{uuid}/confirm", api.ConfirmEmailIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/email/{uuid}", api.UpdateEmailIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/email/{uuid}", api.GetEmailIntegration).Methods(http.MethodGet)
	rtr.HandleFunc("/email/{uuid}/resend", api.ResendConfirmationEmail).Methods(http.MethodPost)

	rtr.HandleFunc("/integration/{uuid}/delete", api.DeleteIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/organization/{uuid}/integrations/{io_type}", api.GetAllIntegrations).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/slack/authorizations", api.GetAllSlackAuthorizations).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/slack/authorization", api.CreateSlackAuthorization).Methods(http.MethodPost)

	rtr.HandleFunc("/integration/sql/test/{type}", api.TestSQLIntegrationConnection).Methods(http.MethodPost)
	rtr.HandleFunc("/integration/sql/{uuid}/test", api.TestExistingSQLIntegrationConnection).Methods(http.MethodPost)
	rtr.HandleFunc("/integration/sql/{uuid}/db/{db}/table/{table}", api.DescribeSQLIntegrationTables).Methods(http.MethodGet)
	rtr.HandleFunc("/integration/sql/{uuid}/db/{db}/tables", api.ShowSQLIntegrationTables).Methods(http.MethodGet)
	rtr.HandleFunc("/integration/sql/{uuid}/dbs", api.ShowSQLIntegrationDatabases).Methods(http.MethodGet)

	rtr.HandleFunc("/jira/authorize", api.AuthorizeJira).
		Queries("code", "{code}").
		Queries("state", "{state}").
		Methods(http.MethodGet)
	rtr.HandleFunc("/jira/authorization/{uuid}", api.UpdateJiraAuthorization).Methods(http.MethodPost)
	rtr.HandleFunc("/jira/authorization/{uuid}", api.GetJiraAuthorization).Methods(http.MethodGet)
	rtr.HandleFunc("/jira", api.CreateJiraIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/jira/{uuid}", api.UpdateJiraIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/jira/{uuid}", api.GetJiraIntegration).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/jira/authorizations", api.GetAllJiraAuthorizations).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/jira/authorization", api.CreateJiraAuthorization).Methods(http.MethodPost)

	rtr.HandleFunc("/salesforce/authorize", api.AuthorizeSalesforce).
		Queries("code", "{code}").
		Queries("state", "{state}").
		Methods(http.MethodGet)
	rtr.HandleFunc("/salesforce/authorization/{uuid}", api.UpdateSalesforceAuthorization).Methods(http.MethodPost)
	rtr.HandleFunc("/salesforce/authorization/{uuid}", api.GetSalesforceAuthorization).Methods(http.MethodGet)
	rtr.HandleFunc("/salesforce", api.CreateSalesforceIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/salesforce/{uuid}", api.UpdateSalesforceIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/salesforce/{uuid}", api.GetSalesforceIntegration).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/salesforce/authorizations", api.GetAllSalesforceAuthorizations).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/salesforce/authorization", api.CreateSalesforceAuthorization).Methods(http.MethodPost)

	rtr.HandleFunc("/incidentio", api.CreateIncidentIoIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/incidentio/{uuid}", api.UpdateIncidentIoIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/incidentio/{uuid}", api.GetIncidentIoIntegration).Methods(http.MethodGet)
	rtr.HandleFunc("/incidentio/{uuid}/severities", api.GetIncidentIoSeverities).Methods(http.MethodGet)

	rtr.HandleFunc("/opsgenie", api.CreateOpsgenieIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/opsgenie/{uuid}", api.UpdateOpsgenieIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/opsgenie/{uuid}", api.GetOpsgenieIntegration).Methods(http.MethodGet)

	rtr.HandleFunc("/jenkins", api.CreateJenkinsIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/jenkins/{uuid}", api.UpdateJenkinsIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/jenkins/{uuid}", api.GetJenkinsIntegration).Methods(http.MethodGet)

	rtr.HandleFunc("/tableau", api.CreateTableauIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/tableau/{uuid}", api.UpdateTableauIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/tableau/{uuid}", api.GetTableauIntegration).Methods(http.MethodGet)

	rtr.HandleFunc("/sql/{type}/{uuid}", api.UpdateSQLIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/sql/{type}", api.CreateSQLIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/sql/{uuid}", api.GetSQLIntegration).Methods(http.MethodGet)

	rtr.HandleFunc("/hubspot/authorize", api.AuthorizeHubspot).
		Queries("code", "{code}").
		Queries("state", "{state}").
		Methods(http.MethodGet)
	rtr.HandleFunc("/hubspot/authorization/{uuid}", api.UpdateHubspotAuthorization).Methods(http.MethodPost)
	rtr.HandleFunc("/hubspot/authorization/{uuid}", api.GetHubspotAuthorization).Methods(http.MethodGet)
	rtr.HandleFunc("/hubspot", api.CreateHubspotIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/hubspot/{uuid}", api.UpdateHubspotIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/hubspot/{uuid}", api.GetHubspotIntegration).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/hubspot/authorizations", api.GetAllHubspotAuthorizations).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/hubspot/authorization", api.CreateHubspotAuthorization).Methods(http.MethodPost)

	rtr.HandleFunc("/google-sheets", api.CreateGoogleSheetsIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/google-sheets/{uuid}", api.UpdateGoogleSheetsIntegration).Methods(http.MethodPost)
	rtr.HandleFunc("/google-sheets/{uuid}", api.GetGoogleSheetsIntegration).Methods(http.MethodGet)

	// The following endpoints are for the internal use only and should be
	// protected on the server level from any outside of the network use
	rtr.HandleFunc("/private/integration/{uuid}/make/{status}", api.ChangeIntegrationStatus).Methods(http.MethodPost)
	rtr.HandleFunc("/private/api/{uuid}", api.GetApiIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/email/{uuid}", api.GetEmailIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/sms/{uuid}", api.GetSmsIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/slack/{uuid}", api.GetSlackIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/jira/{uuid}", api.GetJiraIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/incidentio/{uuid}", api.GetIncidentIoIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/opsgenie/{uuid}", api.GetOpsgenieIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/tableau/{uuid}", api.GetTableauIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/hubspot/{uuid}", api.GetHubspotIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/salesforce/{uuid}", api.GetSalesforceIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/google-sheets/{uuid}", api.GetGoogleSheetsIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/jenkins/{uuid}", api.GetJenkinsIntegrationPrivate).Methods(http.MethodGet)
	rtr.HandleFunc("/private/sql/{uuid}", api.GetSQLIntegrationPrivate).Methods(http.MethodGet)

	http.Handle("/", rtr)

	server := &http.Server{
		Addr:    config.Cfg().Listen,
		Handler: nil,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(ctx)
	}()

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func NewServer(listen string) *Server {
	s := &Server{
		Listen: listen,
	}

	return s
}
