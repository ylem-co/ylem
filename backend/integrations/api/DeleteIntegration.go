package api

import (
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

func DeleteIntegration(w http.ResponseWriter, r *http.Request) {
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

	db := helpers.DbConn()
	defer db.Close()

	entity, err := repositories.FindIntegration(db, uuid)

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
		entity.OrganizationUuid,
		services.PermissionActionDelete,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	err = repositories.DeleteIntegration(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.NotifyServiceIntegrationsChanged(db, entity.OrganizationUuid)

	w.WriteHeader(http.StatusNoContent)
}
