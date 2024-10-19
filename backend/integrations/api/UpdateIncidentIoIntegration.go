package api

import (
	"encoding/json"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HttpUpdateIncidentIoIntegration struct {
	Name       string  `json:"name" valid:"type(string)"`
	ApiKey     *string `json:"api_key" valid:"type(*string)"`
	Mode       string  `json:"mode" valid:"type(string)"`
	Visibility string  `json:"visibility" valid:"type(string)"`
}

func UpdateIncidentIoIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Updating Incident Io Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Debugf("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Tracef("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpUpdateIncidentIoIntegration

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

	var entity *entities.IncidentIo
	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Find IncidentIo Integration")
	var err error
	entity, err = repositories.FindIncidentIoIntegration(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Infof("Incident io Integration with uuid %s not found", uuid)
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

	if request.ApiKey != nil && *request.ApiKey != ""  {
		entity.ApiKey.PlainValue = []byte(*request.ApiKey)
		err = encryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.ApiKey)

		if err != nil {
			return
		}
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.Mode
	entity.Visibility = request.Visibility
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	err = repositories.UpdateIncidentIoIntegration(db, entity)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	entity.ApiKey = nil

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
