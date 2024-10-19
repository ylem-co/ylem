package api

import (
	"context"
	"time"
	"fmt"
	"database/sql"
	"encoding/json"
	hubspotclient "github.com/ylem-co/hubspot-client"
	"github.com/ylem-co/salesforce-client"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ylem_integrations/config"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	sqlService "ylem_integrations/services/sql"
	"ylem_integrations/services/aws/kms"
)

type HttpSQLIntegrationType struct {
	Type             string `json:"type" valid:"type(string)"`
}

func createAndWriteSlackGrantLink(w http.ResponseWriter, auth entities.SlackAuthorization) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{
		"url": services.CreateSlackGrantLink(auth.State),
	})
	_, _ = w.Write(jsonResponse)
}

func createAndWriteJiraGrantLink(w http.ResponseWriter, auth entities.JiraAuthorization) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{
		"url": services.CreateJiraCloudGrantLink(auth.State),
	})
	_, _ = w.Write(jsonResponse)
}

func createAndWriteHubspotGrantLink(w http.ResponseWriter, auth entities.HubspotAuthorization) {
	w.Header().Set("Content-Type", "application/json")
	link, err := hubspotclient.CreateGrantLink(auth.State)

	if err != nil {
		log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{
		"url": link,
	})
	
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func createAndWriteSalesforceGrantLink(w http.ResponseWriter, auth entities.SalesforceAuthorization) {
	w.Header().Set("Content-Type", "application/json")
	link, err := salesforceclient.CreateGrantLink(auth.State)

	if err != nil {
		log.Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{
		"url": link,
	})
	
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func encryptSensitiveData(w http.ResponseWriter, r *http.Request, organizationUuid string, plainBox *kms.SecretBox) error {
	encryptionKey, err := services.FetchOrganizationDataKey(organizationUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return err
	}

	key, err := kms.DecryptDataKey(r.Context(), config.Cfg().Aws.KmsKeyId, encryptionKey)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return err
	}

	encryptedApiKey, err := kms.Encrypt(plainBox.PlainValue, key)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return err
	}

	plainBox.SetEncryptedValue(encryptedApiKey).Seal()

	return nil
}

func decryptSensitiveData(w http.ResponseWriter, r *http.Request, organizationUuid string, sealedBox *kms.SecretBox) ([]byte, error) {
	encryptionKey, err := services.FetchOrganizationDataKey(organizationUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, err
	}

	key, err := kms.DecryptDataKey(r.Context(), config.Cfg().Aws.KmsKeyId, encryptionKey)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, err
	}

	decryptedData, err := kms.Decrypt(sealedBox.EncryptedValue, key)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, err
	}

	sealedBox.Open(decryptedData)

	return encryptionKey, nil
}

func assertTypeSupported(integrationType string, w http.ResponseWriter) bool {
	if entities.IsSQLIntegrationTypeSupported(integrationType) {
		return true
	}

	errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "type"})
	w.WriteHeader(http.StatusBadRequest)

	_, err := w.Write(errorJson)
	if err != nil {
		log.Error(err)
		return false
	}

	return false
}

func assertConnectionTypeSupported(connectionType string, w http.ResponseWriter) bool {
	if entities.IsSQLIntegrationConnectionTypeSupported(connectionType) {
		return true
	}

	errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "connection_type"})
	w.WriteHeader(http.StatusBadRequest)
	
	_, err := w.Write(errorJson)
	if err != nil {
		log.Error(err)
		return false
	}

	return false
}

func updateSQLConnection(db *sql.DB, orgUuid string, isIntegrationCreated bool) {
	integrations, _ := repositories.FindAllIntegrationsBelongToOrganization(db, orgUuid, "all")
	if len(integrations.Items) == 0 {
		sqlService.UpdateSQLIntegrationConnection(orgUuid, false)
	} else {
		sqlService.UpdateSQLIntegrationConnection(orgUuid, true)
	}
}

func testSQLConnection(db *sql.DB, entity *entities.SQLIntegration) {
	defer db.Close()
	ctx := context.TODO()

	if err := entity.Decrypt(ctx); err != nil {
		log.Errorf("could not decrypt integration: %s", err.Error())

		return
	}

	password := make([]byte, 0)
	if entity.Password.PlainValue != nil {
		password = entity.Password.PlainValue
	}

	testErr := sqlService.TestSQLIntegrationConnection(
		entity.Type,
		entity.IsSshConnection(),
		sqlService.DefaultSQLIntegrationConnectionConfiguration{
			Host:        string(entity.Host.PlainValue),
			Port:        uint16(entity.Port),
			User:        entity.User,
			Password:    string(password),
			Database:    entity.Database,
			SshHost:     string(entity.SshHost.PlainValue),
			SshPort:     uint16(entity.SshPort),
			SshUser:     entity.SshUser,
			ProjectId:   entity.ProjectId,
			Credentials: string(entity.Credentials.PlainValue),
			SslEnabled:  entity.SslEnabled,
			EsVersion:   entity.EsVersion,
		},
	)

	if err := entity.Encrypt(ctx); err != nil {
		log.Errorf("could not encrypt integration: %s", err.Error())

		return
	}

	if testErr != nil {
		entity.Integration.Status = entities.IntegrationStatusOffline
		entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
		err := repositories.UpdateSQLIntegration(db, entity, false)
		if err != nil {
			log.Error(err.Error())
		}

		return
	}

	entity.Integration.Status = entities.IntegrationStatusOnline
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	err := repositories.UpdateSQLIntegration(db, entity, false)
	if err != nil {
		log.Error(err.Error())

		return
	}
}

func CheckSQLDatabaseExists(w http.ResponseWriter, databases []string, queryDb string) bool {
	found := false
	for _, v := range databases {
		if v == queryDb {
			found = true

			break
		}
	}

	if !found {
		log.Infof("the database is not found")

		helpers.HttpReturnErrorBadRequest(w, fmt.Sprintf("Database %s not found", queryDb), &[]string{})

		return false
	}

	return true
}
