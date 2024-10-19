package api

import (
	"encoding/json"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type HttpCreateApiIntegration struct {
	Name                  string  `json:"name" valid:"type(string)"`
	Url                   string  `json:"url" valid:"type(string)"`
	Method                string  `json:"method" valid:"type(string),optional"`
	OrganizationUuid      string  `json:"organization_uuid" valid:"uuidv4"`
	AuthType              string  `json:"auth_type" valid:"type(string)"`
	AuthBearerToken       *string `json:"auth_bearer_token" valid:"type(*string),optional"`
	AuthBasicUserName     *string `json:"auth_basic_user_name" valid:"type(*string),optional"`
	AuthBasicUserPassword *string `json:"auth_basic_password" valid:"type(*string),optional"`
	AuthHeaderName        *string `json:"auth_header_name" valid:"type(*string),optional"`
	AuthHeaderValue       *string `json:"auth_header_value" valid:"type(*string),optional"`
}

func CreateApiIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpCreateApiIntegration

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

	db := helpers.DbConn()
	defer db.Close()

	cnt, err := repositories.GetCurrentIntegrationCount(db, uuid.MustParse(request.OrganizationUuid))
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
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
		return
	}

	if !canPerformOperation {
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)

		return
	}

	existingEntity, err := repositories.FindIntegrationInOrganizationByValue(db, request.Url, request.OrganizationUuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if existingEntity != nil {
		helpers.HttpReturnErrorBadRequest(w, "Such API endpoint already exists in the organization", &[]string{})

		return
	}

	method := request.Method
	if method == "" {
		method = http.MethodPost
	}

	entity := entities.Api{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			Value:            request.Url,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		Method: method,
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
		&entity,
	)
	if err != nil {
		log.Println(err.Error())

		return
	}

	err = repositories.CreateApiIntegration(db, &entity)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.NotifyServiceIntegrationsChanged(db, entity.Integration.OrganizationUuid)

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
