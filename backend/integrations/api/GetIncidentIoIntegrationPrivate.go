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

type GetIncidentIoIntegrationPrivateResponse struct {
	Integration IntegrationPrivate `json:"integration"`
	DataKey     []byte             `json:"data_key"`
	ApiKey      []byte             `json:"api_key"`
	Mode        string             `json:"mode"`
	Visibility  string             `json:"visibility"`
}

func GetIncidentIoIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.IncidentIo
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindIncidentIoIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	dataKey, err := decryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.ApiKey)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(GetIncidentIoIntegrationPrivateResponse{
		Integration: IntegrationPrivate(entity.Integration),
		DataKey:     dataKey,
		ApiKey:      entity.ApiKey.EncryptedValue,
		Mode:        entity.Integration.Value,
		Visibility:  entity.Visibility,
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
