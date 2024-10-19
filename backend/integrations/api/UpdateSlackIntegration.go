package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

type HttpUpdateSlackIntegration struct {
	Name              string `json:"name" valid:"type(string)"`
	AuthorizationUuid string `json:"authorization_uuid" valid:"uuidv4"`
	Channel           string `json:"channel" valid:"type(string)"`
}

func UpdateSlackIntegration(w http.ResponseWriter, r *http.Request) {
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

	var request HttpUpdateSlackIntegration
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

	db := helpers.DbConn()
	defer db.Close()

	SlackIntegration, err := repositories.FindSlackIntegration(db, uuid)
	if err != nil || SlackIntegration == nil {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		SlackIntegration.Integration.OrganizationUuid,
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	existingEntity, err := repositories.FindIntegrationInOrganizationByValue(db, request.Channel, SlackIntegration.SlackAuthorization.OrganizationUuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if existingEntity != nil && existingEntity.Uuid != uuid {
		helpers.HttpReturnErrorBadRequest(w, "Such Slack channel already exists in the organization", &[]string{})

		return
	}

	SlackIntegration.Integration.Name = request.Name
	if SlackIntegration.Integration.Value != request.Channel ||
		SlackIntegration.SlackAuthorization.Uuid != request.AuthorizationUuid {

		slackAuth, err := hydrateSlackAuthForUpdateSlackIntegration(
			w,
			db,
			request.AuthorizationUuid,
			SlackIntegration.SlackAuthorization.OrganizationUuid,
		)

		if err != nil {
			return
		}

		if slackAuth == nil {
			helpers.HttpReturnErrorNotFound(w)

			return
		}

		SlackClient := slack.New(*slackAuth.AccessToken)
		SlackChannelId, err := services.GetSlackChannelIdFromName(SlackClient, request.Channel)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)

			return
		}

		if SlackChannelId == nil {
			helpers.HttpReturnErrorBadRequest(
				w,
				"Channel with such name was not found (or it's private) in Slack",
				nil,
			)

			return
		}

		err = services.JoinSlackChannel(SlackClient, *SlackChannelId)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)

			return
		}

		SlackIntegration.Integration.Value = request.Channel
		SlackIntegration.SlackChannelId = SlackChannelId
		SlackIntegration.SlackAuthorization.Id = slackAuth.Id
	}

	SlackIntegration.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	err = repositories.UpdateSlackIntegration(db, SlackIntegration)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(SlackIntegration)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}

func hydrateSlackAuthForUpdateSlackIntegration(w http.ResponseWriter, db *sql.DB, uuid string, orgUuid string) (*entities.SlackAuthorization, error) {
	slackAuth, err := repositories.FindSlackAuthorizationByUuid(db, uuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return nil, err
	}

	if slackAuth == nil {
		log.Infof("Slack auth with uuid %s not found", uuid)
		helpers.HttpReturnErrorNotFound(w)

		return nil, nil
	}

	if slackAuth.OrganizationUuid != orgUuid {
		log.Infof("A user from the org %s tried to set slack auth from org %s", orgUuid, slackAuth.OrganizationUuid)
		helpers.HttpReturnErrorForbidden(w)

		return nil, errors.New("forbidden")
	}

	return slackAuth, nil
}
