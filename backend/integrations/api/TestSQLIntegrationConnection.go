package api

import (
	"net/http"
	"encoding/json"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/services"
	"ylem_integrations/services/sql"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type HttpTestSQLIntegration struct {
	OrganizationUuid string  `json:"organization_uuid" valid:"uuidv4"`
	Type             string  `json:"type" valid:"type(string)"`
	Host             string  `json:"host" valid:"type(string)"`
	Port             int     `json:"port" valid:"type(int)"`
	User             string  `json:"user" valid:"type(string)"`
	Password         *string `json:"password" valid:"type(*string),optional"`
	Database         string  `json:"database" valid:"type(string)"`
	ConnectionType   string  `json:"connection_type" valid:"type(string)"`
	SshHost          string  `json:"ssh_host" valid:"type(string)"`
	SshPort          int     `json:"ssh_port" valid:"type(int)"`
	SshUser          string  `json:"ssh_user" valid:"type(string)"`
	SslEnabled       bool    `json:"ssl_enabled" valid:"type(bool)"`
}

type HttpTestBigQuerySQLIntegration struct {
	OrganizationUuid string  `json:"organization_uuid" valid:"type(string)"`
	Type             string  `json:"type" valid:"type(string)"`
	ProjectId        *string `json:"project_id" valid:"type(*string),optional"`
	Credentials      string  `json:"credentials" valid:"type(string)"`
}

func TestSQLIntegrationConnection(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	integrationType := mux.Vars(r)["type"]
	if !assertTypeSupported(integrationType, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if integrationType == entities.SQLIntegrationTypeGoogleBigQuery {
		testBigQuerySQLIntegrationConnection(user, w, r)
	} else {
		testGenericSQLIntegrationConnection(user, w, r)
	}

}

func testGenericSQLIntegrationConnection(user *services.AuthenticationData, w http.ResponseWriter, r *http.Request) {
	var request HttpTestSQLIntegration

	decodeJsonErr := helpers.DecodeJSONBody(w, r, &request)
	if decodeJsonErr != nil {
		log.Infof("Cannot decode JSON input body")
		rp, _ := json.Marshal(decodeJsonErr.Msg)
		w.WriteHeader(decodeJsonErr.Status)

		_, err := w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		request.OrganizationUuid,
		services.PermissionActionCreate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		log.Error("Operation is not permitted")
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	if !assertConnectionTypeSupported(request.ConnectionType, w) {
		log.Error("Connection type is not supported")
		return
	}

	password := ""
	if request.Password != nil {
		password = *request.Password
	}
	testErr := sql.TestSQLIntegrationConnection(
		request.Type,
		request.ConnectionType == entities.SQLIntegrationConnectionTypeSsh,
		sql.DefaultSQLIntegrationConnectionConfiguration{
			Host:       request.Host,
			Port:       uint16(request.Port),
			User:       request.User,
			Password:   password,
			Database:   request.Database,
			SshHost:    request.SshHost,
			SshPort:    uint16(request.SshPort),
			SshUser:    request.SshUser,
			SslEnabled: request.SslEnabled,
		},
	)
	if testErr != nil {
		helpers.HttpReturnErrorBadRequest(w, testErr.Error(), &[]string{})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func testBigQuerySQLIntegrationConnection(user *services.AuthenticationData, w http.ResponseWriter, r *http.Request) {
	var request HttpTestBigQuerySQLIntegration

	decodeJsonErr := helpers.DecodeJSONBody(w, r, &request)
	if decodeJsonErr != nil {
		log.Info("Cannot decode JSON input body")
		rp, _ := json.Marshal(decodeJsonErr.Msg)
		w.WriteHeader(decodeJsonErr.Status)
		
		_, err := w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		request.OrganizationUuid,
		services.PermissionActionCreate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	testErr := sql.TestSQLIntegrationConnection(
		request.Type,
		false,
		sql.DefaultSQLIntegrationConnectionConfiguration{
			ProjectId:   request.ProjectId,
			Credentials: request.Credentials,
		},
	)

	if testErr != nil {
		helpers.HttpReturnErrorBadRequest(w, testErr.Error(), &[]string{})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
