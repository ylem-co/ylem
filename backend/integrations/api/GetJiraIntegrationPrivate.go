package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

type GetJiraIntegrationPrivateResponse struct {
	Integration IntegrationPrivate `json:"integration"`
	IssueType   string             `json:"issue_type"`
	DataKey     []byte             `json:"data_key"`
	AccessToken []byte             `json:"access_token"`
	CloudId     string             `json:"cloudid"`
}

func GetJiraIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Getting a Jira Integration private")

	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	log.Tracef("Find a Integration")
	var entity *entities.Jira
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindJiraIntegration(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Infof("Jira Integration with uuid %s not found", uuid)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	dataKey, err := services.FetchOrganizationDataKey(entity.Integration.OrganizationUuid)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	response := GetJiraIntegrationPrivateResponse{
		Integration: IntegrationPrivate(entity.Integration),
		IssueType:   entity.IssueType,
		DataKey:     dataKey,
		AccessToken: entity.JiraAuthorization.AccessToken.EncryptedValue,
		CloudId:     *entity.JiraAuthorization.Cloudid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(response)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
