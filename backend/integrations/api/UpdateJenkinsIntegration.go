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

type HttpUpdateJenkinsIntegration struct {
	Name        string  `json:"name" valid:"type(string)"`
	BaseUrl     string  `json:"base_url" valid:"type(string),required"`
	Token       *string `json:"token" valid:"type(*string)"`
	ProjectName string  `json:"project_name" valid:"type(string),required"`
}

func UpdateJenkinsIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Updating Jenkins Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Debugf("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Tracef("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpUpdateJenkinsIntegration

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

	var entity *entities.Jenkins
	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Find Jenkins Integration")
	var err error
	entity, err = repositories.FindJenkinsIntegration(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Infof("Jenkins Integration with uuid %s not found", uuid)
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

	if request.Token != nil && *request.Token != "" {
		entity.Token.PlainValue = []byte(*request.Token)
		err = encryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.Token)

		if err != nil {
			return
		}
	}

	entity.Integration.Name = request.Name
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	entity.Integration.Value = request.ProjectName
	entity.BaseUrl = request.BaseUrl

	err = repositories.UpdateJenkinsIntegration(db, entity)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	entity.Token = nil

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
