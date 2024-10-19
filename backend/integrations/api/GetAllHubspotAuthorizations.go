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

func GetAllHubspotAuthorizations(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	OrganizationUuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(OrganizationUuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		OrganizationUuid,
		services.PermissionActionReadList,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	Collection, err := repositories.FindAllHubspotAuthorizationsForOrganization(db, OrganizationUuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(
		map[string][]entities.HubspotAuthorization{"items": Collection.Items},
	)
	
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
