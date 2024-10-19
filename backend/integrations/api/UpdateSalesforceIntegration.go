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

type HttpUpdateSalesforceIntegration struct {
	Name              string `json:"name" valid:"type(string)"`
	AuthorizationUuid string `json:"authorization_uuid" valid:"uuidv4"`
}

func UpdateSalesforceIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Updating Salesforce Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Debugf("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Tracef("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpUpdateSalesforceIntegration

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

	var entity *entities.Salesforce
	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Find Salesforce Integration")
	var err error
	entity, err = repositories.FindSalesforceIntegration(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Infof("Salesforce Integration with uuid %s not found", uuid)
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

	log.Tracef("Find Salesforce Authorization")
	authorization, err := repositories.FindSalesforceAuthorizationByUuid(db, request.AuthorizationUuid)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if authorization == nil {
		log.Infof("Salesforce authorization with uuid %s not found", request.AuthorizationUuid)
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	entity.Integration.Name = request.Name
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	entity.SalesforceAuthorization = *authorization

	err = repositories.UpdateSalesforceIntegration(db, entity)
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
