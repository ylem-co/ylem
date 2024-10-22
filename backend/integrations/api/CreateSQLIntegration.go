package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	"ylem_integrations/services/aws/kms"
	"ylem_integrations/services/es"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"


)

type HttpCreateSQLIntegration struct {
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
	EsVersion        *uint8  `json:"es_version" valid:"type(*uint8),optional"`
}

type HttpCreateBigQuerySQLIntegration struct {
	Name             string  `json:"name" valid:"type(string)"`
	Type             string  `json:"type" valid:"type(string)"`
	ProjectId        *string `json:"project_id" valid:"type(*string),optional"`
	Credentials      string  `json:"credentials" valid:"type(string)"`
	OrganizationUuid string  `json:"organization_uuid" valid:"uuidv4"`
}

func CreateSQLIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Creating SQL Integration")

	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Infof("failed to authenticate a user")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	sqlType := mux.Vars(r)["type"]
	if !assertTypeSupported(sqlType, w) {
		log.Tracef("SQL type is not supported")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	db := helpers.DbConn()
	var entity *entities.SQLIntegration

	if sqlType == entities.SQLIntegrationTypeGoogleBigQuery {
		entity = createBigQuerySQLIntegration(db, user, w, r)
	} else if sqlType == entities.SQLIntegrationTypeElasticSearch {
		entity = createElasticsearchSQLIntegration(db, user, w, r)
	} else {
		entity = createGenericSQLIntegration(db, user, w, r)
	}

	if entity == nil { // response is already written in the functions above
		db.Close()
		return
	}

	go testSQLConnection(db, entity)
}

func createElasticsearchSQLIntegration(db *sql.DB, user *services.AuthenticationData, w http.ResponseWriter, r *http.Request) *entities.SQLIntegration {
	log.Tracef("Creating Elastic Search Integration")
	ctx := r.Context()

	var request HttpCreateSQLIntegration

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

	cnt, err := repositories.GetCurrentIntegrationCount(db, uuid.MustParse(request.OrganizationUuid))
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return nil
	}
	canPerformOperation, err := services.ValidateBilledPermissions(
		user.Uuid,
		request.OrganizationUuid,
		services.PermissionActionCreate,
		services.PermissionResourceTypeIntegration,
		"",
		cnt,
	)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return nil
	}

	if !canPerformOperation {
		log.Error("Cannot create any new integration. Limit is reached")
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)

		return nil
	}

	dataKey, err := kms.IssueDataKeyWithContext(context.TODO())

	if err != nil {
		log.Errorf("could not issue a data key: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	password := ""
	if request.Password != nil {
		password = *request.Password
	}

	entity := &entities.SQLIntegration{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			Value:            request.Type,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		Type:             request.Type,
		Host:             kms.NewOpenSecretBox([]byte(request.Host)),
		Port:             request.Port,
		User:             request.User,
		Password:         kms.NewOpenSecretBox([]byte(password)),
		Database:         request.Database,
		ConnectionType:   request.ConnectionType,
		SshHost:          kms.NewOpenSecretBox([]byte(request.SshHost)),
		SshPort:          request.SshPort,
		SshUser:          request.SshUser,
		SslEnabled:       request.SslEnabled,
		DataKey:          kms.NewSealedSecretBox(dataKey),
	}

	esClient, err := es.NewConnection(
		ctx,
		fmt.Sprintf("https://%s:%d", request.Host, request.Port),
		request.User,
		password,
		nil,
	)
	if err != nil {
		log.Errorf("could not create es client: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}
	esVersion, err := esClient.Version()
	if err != nil {
		log.Errorf("could not create es client: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}
	entity.EsVersion = &esVersion

	if err := entity.Encrypt(ctx); err != nil {
		log.Errorf("could not encrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	_, err = repositories.CreateSQLIntegration(db, entity)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return nil
	}

	updateSQLConnection(db, request.OrganizationUuid, true)

	if err := entity.Decrypt(ctx); err != nil {
		log.Errorf("could not decrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	w.WriteHeader(http.StatusCreated)

	return entity
}

func createGenericSQLIntegration(db *sql.DB, user *services.AuthenticationData, w http.ResponseWriter, r *http.Request) *entities.SQLIntegration {
	log.Tracef("Creating Generic SQL Integration")
	ctx := r.Context()

	var request HttpCreateSQLIntegration

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

	cnt, err := repositories.GetCurrentIntegrationCount(db, uuid.MustParse(request.OrganizationUuid))
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return nil
	}
	canPerformOperation, err := services.ValidateBilledPermissions(
		user.Uuid,
		request.OrganizationUuid,
		services.PermissionActionCreate,
		services.PermissionResourceTypeIntegration,
		"",
		cnt,
	)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return nil
	}

	if !canPerformOperation {
		log.Error("Cannot create any new integration. Limit is reached")
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)

		return nil
	}

	if !assertConnectionTypeSupported(request.ConnectionType, w) {
		log.Error("Connection type is not supported")
		return nil
	}

	dataKey, err := kms.IssueDataKeyWithContext(context.TODO())

	if err != nil {
		log.Errorf("could not issue a data key: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	entity := &entities.SQLIntegration{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			Value:            request.Type,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		Type:             request.Type,
		Host:             kms.NewOpenSecretBox([]byte(request.Host)),
		Port:             request.Port,
		User:             request.User,
		Password:         kms.NewOpenSecretBox([]byte(*request.Password)),
		Database:         request.Database,
		ConnectionType:   request.ConnectionType,
		SshHost:          kms.NewOpenSecretBox([]byte(request.SshHost)),
		SshPort:          request.SshPort,
		SshUser:          request.SshUser,
		SslEnabled:       request.SslEnabled,
		DataKey:          kms.NewSealedSecretBox(dataKey),
	}

	isTrialHost := entities.IsTrialHost(entity.Host.PlainValue)
	if isTrialHost {
		entity.IsTrial = 1
	} else {
		entity.IsTrial = 0
	}

	if err := entity.Encrypt(ctx); err != nil {
		log.Errorf("could not encrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	ent, err := repositories.CreateSQLIntegration(db, entity)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return nil
	}

	if !isTrialHost {
		updateSQLConnection(db, request.OrganizationUuid, true)
	}

	if err := entity.Decrypt(ctx); err != nil {
		log.Errorf("could not decrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	entity.Integration.Uuid = ent.Integration.Uuid
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return nil
	}

	return entity
}

func createBigQuerySQLIntegration(db *sql.DB, user *services.AuthenticationData, w http.ResponseWriter, r *http.Request) *entities.SQLIntegration {
	log.Tracef("Creating BigQuery SQL Integration")

	var request HttpCreateBigQuerySQLIntegration
	ctx := r.Context()

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

	cnt, err := repositories.GetCurrentIntegrationCount(db, uuid.MustParse(request.OrganizationUuid))
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return nil
	}
	canPerformOperation, err := services.ValidateBilledPermissions(
		user.Uuid,
		request.OrganizationUuid,
		services.PermissionActionCreate,
		services.PermissionResourceTypeIntegration,
		"",
		cnt,
	)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return nil
	}

	if !canPerformOperation {
		log.Error("Cannot create any new integration. Limit is reached")
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)

		return nil
	}

	dataKey, err := kms.IssueDataKeyWithContext(context.TODO())

	if err != nil {
		log.Errorf("could not issue a data key: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	entity := &entities.SQLIntegration{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			Value:            entities.SQLIntegrationTypeGoogleBigQuery,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		Host:             kms.NewOpenSecretBox([]byte("")),
		SshHost:          kms.NewOpenSecretBox([]byte("")),
		Type:             entities.SQLIntegrationTypeGoogleBigQuery,
		ConnectionType:   entities.SQLIntegrationConnectionTypeDirect,
		ProjectId:        request.ProjectId,
		Credentials:      kms.NewOpenSecretBox([]byte(request.Credentials)),
		DataKey:          kms.NewSealedSecretBox(dataKey),
	}

	if err := entity.Encrypt(ctx); err != nil {
		log.Errorf("could not encrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	_, err = repositories.CreateSQLIntegration(db, entity)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return nil
	}

	updateSQLConnection(db, request.OrganizationUuid, true)

	if err := entity.Decrypt(ctx); err != nil {
		log.Errorf("could not decrypt integration: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	w.WriteHeader(http.StatusCreated)

	return entity
}
