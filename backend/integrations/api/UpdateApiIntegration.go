package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type HttpUpdateApiIntegration struct {
	Name                  string  `json:"name" valid:"type(string)"`
	Url                   string  `json:"url" valid:"type(string)"`
	Method                string  `json:"method" valid:"type(string),optional"`
	AuthType              string  `json:"auth_type" valid:"type(string)"`
	AuthBearerToken       *string `json:"auth_bearer_token" valid:"type(*string),optional"`
	AuthBasicUserName     *string `json:"auth_basic_user_name" valid:"type(*string),optional"`
	AuthBasicUserPassword *string `json:"auth_basic_password" valid:"type(*string),optional"`
	AuthHeaderName        *string `json:"auth_header_name" valid:"type(*string),optional"`
	AuthHeaderValue       *string `json:"auth_header_value" valid:"type(*string),optional"`
}

func UpdateApiIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpUpdateApiIntegration

	w.Header().Set("Content-Type", "application/json")

	decodeJsonErr := helpers.DecodeJSONBody(w, r, &request)
	if decodeJsonErr != nil {
		rp, _ := json.Marshal(decodeJsonErr.Msg)
		w.WriteHeader(decodeJsonErr.Status)
		
		_, err := w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		return
	}

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
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		entity.Integration.OrganizationUuid,
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	existingEntity, err := repositories.FindIntegrationInOrganizationByValue(db, request.Url, entity.Integration.OrganizationUuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if existingEntity != nil && existingEntity.Uuid != uuid {
		helpers.HttpReturnErrorBadRequest(w, "Such API endpoint already exists in the organization", &[]string{})

		return
	}

	err = helpers.ProcessHttpAuthTypeRequest(
		w,
		helpers.AuthTypeCredentialsToProcess{
			AuthType:      request.AuthType,
			Bearer:        request.AuthBearerToken,
			BasicUsername: request.AuthBasicUserName,
			BasicPassword: request.AuthBasicUserPassword,
			HeaderName:    request.AuthHeaderName,
			HeaderValue:   request.AuthHeaderValue,
		},
		entity,
	)
	if err != nil {
		log.Println(err.Error())

		return
	}

	method := request.Method
	if method == "" {
		method = entity.Method
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.Url
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	entity.Method = method

	err = repositories.UpdateApiIntegration(db, entity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
