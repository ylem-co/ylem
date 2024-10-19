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

type HttpCreateHubspotIntegration struct {
	Name              string `json:"name" valid:"type(string)"`
	AuthorizationUuid string `json:"authorization_uuid" valid:"uuidv4"`
	PipelineCode      string `json:"pipeline_code" valid:"type(string)"`
	PipelineStageCode string `json:"pipeline_stage_code" valid:"type(string)"`
	OwnerCode         string `json:"owner_code" valid:"type(string)"`
}

func CreateHubspotIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Creating Hubspot Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("authorization"))
	if user == nil {
		log.Debugf("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Tracef("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpCreateHubspotIntegration

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

	log.Tracef("Find an authorization")
	authorization, err := repositories.FindHubspotAuthorizationByUuid(db, request.AuthorizationUuid)
	if err != nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	if authorization == nil {
		log.Info("Hubspot authorization not found. UUID: " + request.AuthorizationUuid)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	cnt, err := repositories.GetCurrentIntegrationCount(db, uuid.MustParse(authorization.OrganizationUuid))
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	canPerformOperation, err := services.ValidateBilledPermissions(
		user.Uuid,
		authorization.OrganizationUuid,
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
		log.Debugf(
			"User %s can't perform the operation %s in %s",
			user.Uuid,
			services.PermissionActionCreate,
			services.PermissionResourceTypeIntegration,
		)
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)

		return
	}

	log.Tracef("Creating the Integration")
	entity := entities.Hubspot{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: authorization.OrganizationUuid,
			Name:             request.Name,
			Value:            request.PipelineCode,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		HubspotAuthorization: entities.HubspotAuthorization{
			Id:       authorization.Id,
			Name:     authorization.Name,
			Uuid:     authorization.Uuid,
			IsActive: authorization.IsActive,
		},
		PipelineStageCode: request.PipelineStageCode,
		OwnerCode:         request.OwnerCode,
	}

	err = repositories.CreateHubspotIntegration(db, &entity)

	if err != nil {
		log.Error(err.Error())
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

	log.Tracef("New Hubspot integration was created")
}
