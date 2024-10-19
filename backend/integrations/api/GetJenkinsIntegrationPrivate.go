package api

import (
	"encoding/json"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type GetJenkinsIntegrationPrivateResponse struct {
	Integration IntegrationPrivate `json:"integration"`
	DataKey     []byte             `json:"data_key"`
	Token       []byte             `json:"token"`
	BaseUrl     string             `json:"base_url"`
}

func GetJenkinsIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.Jenkins
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindJenkinsIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	dataKey, err := decryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.Token) // wtf??
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(GetJenkinsIntegrationPrivateResponse{
		Integration: IntegrationPrivate(entity.Integration),
		DataKey:     dataKey,
		Token:       entity.Token.EncryptedValue,
		BaseUrl:     entity.BaseUrl,
	})

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
