package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

func GetJenkinsIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

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

	var err error
	entity, err = repositories.FindJenkinsIntegration(db, uuid)

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
		services.PermissionActionRead,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	entity.Token = nil

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(entity)

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
