package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	hubspotclient "github.com/ylem-co/hubspot-client"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

type GetHubspotIntegrationPrivateResponse struct {
	Integration       IntegrationPrivate `json:"integration"`
	PipelineStageCode string             `json:"pipeline_stage_code"`
	OwnerCode         string             `json:"owner_code"`
	DataKey           []byte             `json:"data_key"`
	AccessToken       []byte             `json:"access_token"`
}

func GetHubspotIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Getting a private Hubspot Integration")

	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	log.Tracef("Find a Integration")
	var entity *entities.Hubspot
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindHubspotIntegration(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Infof("Hubspot Integration with uuid %s not found", uuid)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	_, err = decryptSensitiveData(
		w,
		r,
		entity.Integration.OrganizationUuid,
		entity.HubspotAuthorization.AccessToken,
	)
	if err != nil {
		return
	}

	_, err = decryptSensitiveData(
		w,
		r,
		entity.Integration.OrganizationUuid,
		entity.HubspotAuthorization.RefreshToken,
	)
	if err != nil {
		return
	}

	token := &oauth2.Token{
		AccessToken:  string(entity.HubspotAuthorization.AccessToken.PlainValue),
		RefreshToken: string(entity.HubspotAuthorization.RefreshToken.PlainValue),
		Expiry:       entity.HubspotAuthorization.AccessTokenExpiresAt,
	}

	client, err := hubspotclient.CreateInstance(r.Context(), token)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	cToken, err := client.Token()
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	// the token was refreshed
	if cToken.AccessToken != token.AccessToken {
		accessTokenSealedBox, _, err := encryptHubspotTokens(w, r, &entity.HubspotAuthorization, cToken)
		if err != nil {
			log.Error(err.Error())
			helpers.HttpReturnErrorInternal(w)

			return
		}

		entity.HubspotAuthorization.AccessToken = accessTokenSealedBox
		entity.HubspotAuthorization.AccessTokenExpiresAt = cToken.Expiry.Add(services.HubspotTokenSubtractTime * time.Minute)

		err = repositories.UpdateHubspotAuthorization(db, &entity.HubspotAuthorization)
		if err != nil {
			log.Error(err.Error())
			helpers.HttpReturnErrorInternal(w)

			return
		}
	}

	dataKey, err := services.FetchOrganizationDataKey(entity.Integration.OrganizationUuid)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	response := GetHubspotIntegrationPrivateResponse{
		Integration:       IntegrationPrivate(entity.Integration),
		PipelineStageCode: entity.PipelineStageCode,
		OwnerCode:         entity.OwnerCode,
		DataKey:           dataKey,
		AccessToken:       entity.HubspotAuthorization.AccessToken.EncryptedValue,
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
