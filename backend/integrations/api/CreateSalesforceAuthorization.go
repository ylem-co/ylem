package api

import (
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

func CreateSalesforceAuthorization(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Creating Salesforce Authorization")

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
		services.PermissionActionCreate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	uniqueCode, err := password.Generate(64, 20, 0, false, true)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	auth := entities.SalesforceAuthorization{
		Name:             "New Salesforce Authorization",
		CreatorUuid:      user.Uuid,
		OrganizationUuid: OrganizationUuid,
		State:            uniqueCode,
		IsActive:         false,
	}

	err = repositories.CreateSalesforceAuthorization(db, &auth)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	createAndWriteSalesforceGrantLink(w, auth)
}
