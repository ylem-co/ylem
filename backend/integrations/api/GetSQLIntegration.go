package api

import (
	"encoding/json"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	validation "github.com/go-ozzo/ozzo-validation"
	log "github.com/sirupsen/logrus"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
)

func GetSQLIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	uuid := mux.Vars(r)["uuid"]
	err := validation.Validate(uuid, validation.Required, is.UUIDv4)
	if err != nil {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.SQLIntegration
	db := helpers.DbConn()
	defer db.Close()

	entity, err = repositories.FindSQLIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

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

	if err := entity.Decrypt(r.Context()); err != nil {
		log.Errorf("could not decrypt source: %s", err.Error())

		return
	}

	entity.MaskHosts()

	jsonResponse, err := json.Marshal(entity)

	if err != nil {
		log.Errorf("could not marshal an entity: %s", err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
