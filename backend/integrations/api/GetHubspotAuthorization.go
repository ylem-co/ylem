package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

func GetHubspotAuthorization(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")

	db := helpers.DbConn()
	defer db.Close()

	authorization, err := repositories.FindHubspotAuthorizationByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if authorization == nil {
		log.Infof("Hubspot authorization %s was not found", uuid)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		authorization.OrganizationUuid,
		services.PermissionActionRead,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]interface{}{
		"model": authorization,
	})
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
