package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

type HttpUpdateJiraIntegration struct {
	Name              string `json:"name" valid:"type(string)"`
	AuthorizationUuid string `json:"authorization_uuid" valid:"uuidv4"`
	ProjectKey        string `json:"project_key" valid:"type(string)"`
	IssueType         string `json:"issue_type" valid:"type(string)"`
}

func UpdateJiraIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Updating Jira Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Debugf("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Tracef("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpUpdateJiraIntegration

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

	var entity *entities.Jira
	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Find Jira Integration")
	var err error
	entity, err = repositories.FindJiraIntegration(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Infof("Jira Integration with uuid %s not found", uuid)
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

	log.Tracef("Find Jira Authorization")
	authorization, err := repositories.FindJiraAuthorizationByUuid(db, request.AuthorizationUuid)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if authorization == nil {
		log.Infof("Jira authorization with uuid %s not found", request.AuthorizationUuid)
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.ProjectKey
	entity.IssueType = request.IssueType
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	entity.JiraAuthorization = *authorization

	err = repositories.UpdateJiraIntegration(db, entity)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
