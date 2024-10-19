package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type HttpUpdateSQLIntegration struct {
	Name             string  `json:"name" valid:"type(string)"`
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
	OrganizationUuid string  `json:"organization_uuid" valid:"uuidv4"`
	SslEnabled       bool    `json:"ssl_enabled" valid:"type(bool)"`
}

type HttpUpdateBigQuerySQLIntegration struct {
	Name             string  `json:"name" valid:"type(string)"`
	Type             string  `json:"type" valid:"type(string)"`
	ProjectId        *string `json:"project_id" valid:"type(*string),optional"`
	Credentials      string  `json:"credentials" valid:"type(string)"`
	OrganizationUuid string  `json:"organization_uuid" valid:"uuidv4"`
}

type HttpUpdateElasticSearchSQLIntegration struct {
	Name             string  `json:"name" valid:"type(string)"`
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
	OrganizationUuid string  `json:"organization_uuid" valid:"uuidv4"`
	SslEnabled       bool    `json:"ssl_enabled" valid:"type(bool)"`
	EsVersion        uint8   `json:"es_version" valid:"type(uint8)"`
}

func UpdateSQLIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Updating SQL Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Infof("failed to authenticate a user")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	SQLIntegrationType := mux.Vars(r)["type"]
	if !assertTypeSupported(SQLIntegrationType, w) {
		log.Tracef("SQL type is not supported")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	db := helpers.DbConn()
	var entity *entities.SQLIntegration

	if SQLIntegrationType == entities.SQLIntegrationTypeGoogleBigQuery {
		entity = updateBigQuerySQLIntegration(db, user, w, r)
	} else if SQLIntegrationType == entities.SQLIntegrationTypeElasticSearch {
		entity = updateElasticsearchSQLIntegration(db, user, w, r)
	} else {
		entity = updateGenericSQLIntegration(db, user, w, r)
	}

	if entity == nil { // response is already written in the functions above
		return
	}

	go testSQLConnection(db, entity)
}

func updateElasticsearchSQLIntegration(db *sql.DB, user *services.AuthenticationData, w http.ResponseWriter, r *http.Request) *entities.SQLIntegration {
	log.Tracef("Updating Elastic Search Integration")

	var request HttpUpdateElasticSearchSQLIntegration
	ctx := context.TODO()

	decodeJsonErr := helpers.DecodeJSONBody(w, r, &request)
	if decodeJsonErr != nil {
		log.Infof("Cannot decode JSON input body")
		rp, _ := json.Marshal(decodeJsonErr.Msg)
		w.WriteHeader(decodeJsonErr.Status)
		
		_, err := w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return nil
		}

		return nil
	}

	uuid := mux.Vars(r)["uuid"]
	err := validation.Validate(uuid, validation.Required, is.UUIDv4)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorBadUuidRequest(w)

		return nil
	}

	var entity *entities.SQLIntegration

	entity, err = repositories.FindSQLIntegration(db, uuid)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return nil
	}

	if entity == nil {
		log.Error("Entity is not found")
		helpers.HttpReturnErrorForbidden(w)

		return nil
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		entity.Integration.OrganizationUuid,
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		log.Error("Operation is not permitted")
		helpers.HttpReturnErrorForbidden(w)

		return nil
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.Type
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	entity.Host.SetPlainValue([]byte(request.Host))
	entity.Port = request.Port
	entity.Database = request.Database
	entity.User = request.User
	entity.EsVersion = &request.EsVersion
	
	rewritePassword := false
	if request.Password != nil && len(*request.Password) > 0 {
		entity.Password.SetPlainValue([]byte(*request.Password))
		rewritePassword = true
	}

	if err = entity.Encrypt(ctx); err != nil {
		log.Errorf("could not encrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	err = repositories.UpdateSQLIntegration(db, entity, rewritePassword)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	w.WriteHeader(http.StatusCreated)

	return entity
}

func updateGenericSQLIntegration(db *sql.DB, user *services.AuthenticationData, w http.ResponseWriter, r *http.Request) *entities.SQLIntegration {
	log.Tracef("Updating Generic SQL Integration")

	var request HttpUpdateSQLIntegration
	ctx := context.TODO()

	decodeJsonErr := helpers.DecodeJSONBody(w, r, &request)
	if decodeJsonErr != nil {
		log.Infof("Cannot decode JSON input body")
		rp, _ := json.Marshal(decodeJsonErr.Msg)
		w.WriteHeader(decodeJsonErr.Status)
		
		_, err := w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return nil
		}

		return nil
	}

	uuid := mux.Vars(r)["uuid"]
	err := validation.Validate(uuid, validation.Required, is.UUIDv4)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorBadUuidRequest(w)

		return nil
	}

	var entity *entities.SQLIntegration

	entity, err = repositories.FindSQLIntegration(db, uuid)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return nil
	}

	if entity == nil {
		log.Error("Entity is not found")
		helpers.HttpReturnErrorForbidden(w)

		return nil
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		entity.Integration.OrganizationUuid,
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		log.Error("Operation is not permitted")
		helpers.HttpReturnErrorForbidden(w)

		return nil
	}

	if !assertConnectionTypeSupported(request.ConnectionType, w) {
		log.Error("Connection type is not supported")
		return nil
	}

	if err := entity.Decrypt(r.Context()); err != nil {
		log.Errorf("could not decrypt integration: %s", err.Error())
		helpers.HttpReturnErrorInternal(w)

		return nil
	}

	if entities.IsTrialHost(entity.Host.PlainValue) {
		log.Error("Cannot update trial host")
		helpers.HttpReturnErrorBadRequest(w, "Can't update the demo host", &[]string{})

		return nil
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.Type
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	entity.Port = request.Port
	entity.User = request.User
	entity.Database = request.Database
	entity.Type = request.Type
	entity.ConnectionType = request.ConnectionType
	entity.SshPort = request.SshPort
	entity.SshUser = request.SshUser
	entity.SslEnabled = request.SslEnabled

	entity.Host.SetPlainValue([]byte(request.Host))
	entity.SshHost.SetPlainValue([]byte(request.SshHost))

	rewritePassword := false
	if request.Password != nil && len(*request.Password) > 0 {
		entity.Password.SetPlainValue([]byte(*request.Password))
		rewritePassword = true
	}

	if err = entity.Encrypt(ctx); err != nil {
		log.Errorf("could not encrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	err = repositories.UpdateSQLIntegration(db, entity, rewritePassword)

	if err != nil {
		log.Errorf("could not update integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	w.WriteHeader(http.StatusCreated)
	if err = entity.Decrypt(ctx); err != nil {
		log.Errorf("could not decrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	return entity
}

func updateBigQuerySQLIntegration(db *sql.DB, user *services.AuthenticationData, w http.ResponseWriter, r *http.Request) *entities.SQLIntegration {
	log.Tracef("Updating Big Query SQL Integration")

	var request HttpUpdateBigQuerySQLIntegration
	ctx := context.TODO()

	decodeJsonErr := helpers.DecodeJSONBody(w, r, &request)
	if decodeJsonErr != nil {
		log.Infof("Cannot decode JSON input body")
		rp, _ := json.Marshal(decodeJsonErr.Msg)
		w.WriteHeader(decodeJsonErr.Status)
		
		_, err := w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return nil
		}

		return nil
	}

	uuid := mux.Vars(r)["uuid"]
	err := validation.Validate(uuid, validation.Required, is.UUIDv4)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorBadUuidRequest(w)

		return nil
	}

	var entity *entities.SQLIntegration

	entity, err = repositories.FindSQLIntegration(db, uuid)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return nil
	}

	if entity == nil {
		log.Error("Entity is not found")
		helpers.HttpReturnErrorForbidden(w)

		return nil
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		entity.Integration.OrganizationUuid,
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		log.Error("Operation is not permitted")
		helpers.HttpReturnErrorForbidden(w)

		return nil
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.Type
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	entity.Type = entities.SQLIntegrationTypeGoogleBigQuery
	entity.ConnectionType = entities.SQLIntegrationConnectionTypeDirect
	entity.ProjectId = request.ProjectId

	entity.Credentials.SetPlainValue([]byte(request.Credentials))
	entity.Password.SetPlainValue(nil)

	if err = entity.Encrypt(ctx); err != nil {
		log.Errorf("could not encrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	err = repositories.UpdateSQLIntegration(db, entity, false)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	w.WriteHeader(http.StatusCreated)

	return entity
}
