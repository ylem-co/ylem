package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	salesforceclient "github.com/ylem-co/salesforce-client"
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

type GetSalesforceIntegrationPrivateResponse struct {
	Integration IntegrationPrivate `json:"integration"`
	DataKey     []byte             `json:"data_key"`
	AccessToken []byte             `json:"access_token"`
	Domain      string             `json:"domain"`
}

func GetSalesforceIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Getting a Salesforce Integration private")

	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	log.Tracef("Find a Integration")
	var entity *entities.Salesforce
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindSalesforceIntegration(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Infof("Salesforce Integration with uuid %s not found", uuid)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	// ------------------------------------------------------------------------------------------

	_, err = decryptSensitiveData(
		w,
		r,
		entity.Integration.OrganizationUuid,
		entity.SalesforceAuthorization.RefreshToken,
	)
	if err != nil {
		return
	}

	token := &oauth2.Token{
		AccessToken:  "",
		RefreshToken: string(entity.SalesforceAuthorization.RefreshToken.PlainValue),
		Expiry:       time.UnixMicro(0),
	}

	client, err := salesforceclient.CreateInstance(r.Context(), *entity.SalesforceAuthorization.Domain, token)
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
	accessTokenSealedBox, _, err := encryptSalesforceTokens(w, r, &entity.SalesforceAuthorization, cToken)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	entity.SalesforceAuthorization.AccessToken = accessTokenSealedBox

	err = repositories.UpdateSalesforceAuthorization(db, &entity.SalesforceAuthorization)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	// ------------------------------------------------------------------------------------------

	dataKey, err := services.FetchOrganizationDataKey(entity.Integration.OrganizationUuid)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	response := GetSalesforceIntegrationPrivateResponse{
		Integration: IntegrationPrivate(entity.Integration),
		DataKey:     dataKey,
		AccessToken: entity.SalesforceAuthorization.AccessToken.EncryptedValue,
		Domain:      *entity.SalesforceAuthorization.Domain,
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
