package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ylem_integrations/config"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	"ylem_integrations/services/aws/kms"
)

func GetJiraAuthorization(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	db := helpers.DbConn()
	defer db.Close()

	jiraAuthorization, err := repositories.FindJiraAuthorizationByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if jiraAuthorization == nil {
		log.Infof("Jira authorization %s was not found", uuid)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		jiraAuthorization.OrganizationUuid,
		services.PermissionActionRead,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	err = decryptJiraAccessToken(w, r, jiraAuthorization)
	if err !=nil {
		log.Error(err.Error())

		return
	}

	resources, err := services.JiraListAvailableResources(string(jiraAuthorization.AccessToken.PlainValue))
	if err !=nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]interface{}{
		"model": jiraAuthorization,
		"resources": resources,
	})
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}

func decryptJiraAccessToken(w http.ResponseWriter, r *http.Request, entity *entities.JiraAuthorization) error {
	encryptionKey, err := services.FetchOrganizationDataKey(entity.OrganizationUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return err
	}

	key, err := kms.DecryptDataKey(r.Context(), config.Cfg().Aws.KmsKeyId, encryptionKey)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return err
	}

	decryptedAccessToken, err := kms.Decrypt(entity.AccessToken.EncryptedValue, key)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return err
	}

	entity.AccessToken.Open(decryptedAccessToken)

	return nil
}
