package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

type HttpUpdateHubspotIntegration struct {
	Name              string `json:"name" valid:"type(string)"`
	AuthorizationUuid string `json:"authorization_uuid" valid:"uuidv4"`
	PipelineCode      string `json:"pipeline_code" valid:"type(string)"`
	PipelineStageCode string `json:"pipeline_stage_code" valid:"type(string)"`
	OwnerCode         string `json:"owner_code" valid:"type(string)"`
}

func UpdateHubspotIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Updating Hubspot Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Debugf("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Tracef("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpUpdateHubspotIntegration

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

	var entity *entities.Hubspot
	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Find Hubspot Integration")
	var err error
	entity, err = repositories.FindHubspotIntegration(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Infof("Hubspot Integration with uuid %s not found", uuid)
		helpers.HttpReturnErrorNotFound(w)

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
		log.Debugf(
			"User %s can't perform the operation %s in %s",
			user.Uuid,
			services.PermissionActionUpdate,
			services.PermissionResourceTypeIntegration,
		)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	log.Tracef("Find Hubspot Authorization")
	authorization, err := repositories.FindHubspotAuthorizationByUuid(db, request.AuthorizationUuid)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if authorization == nil {
		log.Infof("Hubspot authorization with uuid %s not found", request.AuthorizationUuid)
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.PipelineCode
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	entity.PipelineStageCode = request.PipelineStageCode
	entity.OwnerCode = request.OwnerCode
	entity.HubspotAuthorization = *authorization

	err = repositories.UpdateHubspotIntegration(db, entity)
	if err != nil {
		log.Error(err.Error())
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
