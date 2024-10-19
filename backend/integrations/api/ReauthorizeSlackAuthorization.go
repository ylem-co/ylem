package api

import (
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

func ReauthorizeSlackAuthorization(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Reauthorizing Slack Authorization")

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

	slackAuthorization, err := repositories.FindSlackAuthorizationByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if slackAuthorization == nil {
		log.Infof("Slack authorization %s was not found", uuid)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		slackAuthorization.OrganizationUuid,
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	uniqueCode, err := password.Generate(64, 20, 0, false, true)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	slackAuthorization.State = uniqueCode
	slackAuthorization.IsActive = false

	err = repositories.UpdateSlackAuthorization(db, slackAuthorization)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	createAndWriteSlackGrantLink(w, *slackAuthorization)
}
