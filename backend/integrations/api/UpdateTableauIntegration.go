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

type HttpUpdateTableauIntegration struct {
	Name           string  `json:"name" valid:"type(string),required"`
	Server         string  `json:"server" valid:"type(string),required"`
	Username       string  `json:"username" valid:"type(string),required"`
	Password       *string `json:"password" valid:"type(*string)"`
	Sitename       string  `json:"site_name" valid:"type(string),required"`
	ProjectName    string  `json:"project_name" valid:"type(string),required"`
	DatasourceName string  `json:"datasource_name" valid:"type(string),required"`
	Mode           string  `json:"mode" valid:"type(string),in(overwrite|append),required"`
}

func UpdateTableauIntegration(w http.ResponseWriter, r *http.Request) {
	log.Trace("Updating Tableau Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Debug("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Trace("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpUpdateTableauIntegration

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

	var entity *entities.Tableau
	db := helpers.DbConn()
	defer db.Close()

	log.Trace("Find Tableau Integration")
	var err error
	entity, err = repositories.FindTableauIntegration(db, uuid)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Debugf("Tableau Integration with uuid %s not found", uuid)
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

	if request.Password != nil && *request.Password != "" {
		entity.Password.PlainValue = []byte(*request.Password)
		err = encryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.Password)

		if err != nil {
			return
		}
	}

	entity.Username.PlainValue = []byte(request.Username)
	err = encryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.Username)

	if err != nil {
		return
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.Server
	entity.Sitename = request.Sitename
	entity.ProjectName = request.ProjectName
	entity.DatasourceName = request.DatasourceName
	entity.Mode = request.Mode
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	err = repositories.UpdateTableauIntegration(db, entity)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	entity.Username = nil
	entity.Password = nil

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
