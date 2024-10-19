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

type IntegrationPrivate struct {
	Id               int64  `json:"-"`
	Uuid             string `json:"uuid"`
	CreatorUuid      string `json:"creator_uuid"`
	OrganizationUuid string `json:"organization_uuid"`
	Status           string `json:"status"`
	Type             string `json:"type"`
	IoType           string `json:"io_type"`
	Name             string `json:"name"`
	Value            string `json:"value"`
	UserUpdatedAt    string `json:"user_updated_at"`
}

type GetApiIntegrationPrivateResponse struct {
	Id                    int64              `json:"-"`
	Integration           IntegrationPrivate `json:"integration"`
	Method                string             `json:"method"`
	AuthType              string             `json:"auth_type"`
	AuthBearerToken       *string            `json:"auth_bearer_token"`
	AuthBasicUserName     *string            `json:"auth_basic_user_name"`
	AuthBasicUserPassword *string            `json:"auth_basic_user_password"`
	AuthHeaderName        *string            `json:"auth_header_name"`
	AuthHeaderValue       *string            `json:"auth_header_value"`
}

func GetApiIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.Api
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindApiIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	response := GetApiIntegrationPrivateResponse{
		Id:                    entity.Id,
		Integration:           IntegrationPrivate(entity.Integration),
		Method:                entity.Method,
		AuthType:              entity.AuthType,
		AuthBearerToken:       entity.AuthBearerToken,
		AuthBasicUserName:     entity.AuthBasicUserName,
		AuthBasicUserPassword: entity.AuthBasicUserPassword,
		AuthHeaderName:        entity.AuthHeaderName,
		AuthHeaderValue:       entity.AuthHeaderValue,
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
