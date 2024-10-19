package api

import (
	"encoding/json"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	"github.com/google/uuid"
	"github.com/slack-go/slack"

	log "github.com/sirupsen/logrus"
)

type HttpCreateSlackIntegration struct {
	Name              string `json:"name" valid:"type(string)"`
	AuthorizationUuid string `json:"authorization_uuid" valid:"uuidv4"`
	Channel           string `json:"channel" valid:"type(string)"`
}

func CreateSlackIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpCreateSlackIntegration

	w.Header().Set("Content-Type", "application/json")

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

	Authorization, err := repositories.FindSlackAuthorizationByUuid(db, request.AuthorizationUuid)
	if err != nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	if Authorization == nil {
		log.Println("Slack authorization not found. UUID: " + request.AuthorizationUuid)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	cnt, err := repositories.GetCurrentIntegrationCount(db, uuid.MustParse(Authorization.OrganizationUuid))
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	canPerformOperation, err := services.ValidateBilledPermissions(
		user.Uuid,
		Authorization.OrganizationUuid,
		services.PermissionActionCreate,
		services.PermissionResourceTypeIntegration,
		"",
		cnt,
	)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if !canPerformOperation {
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)

		return
	}

	if Authorization.AccessToken == nil || !Authorization.IsActive {
		helpers.HttpReturnErrorBadRequest(
			w,
			"Slack authorization is not active. Please check if it's up to date",
			nil,
		)

		return
	}

	existingEntity, err := repositories.FindIntegrationInOrganizationByValue(db, request.Channel, Authorization.OrganizationUuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if existingEntity != nil {
		helpers.HttpReturnErrorBadRequest(w, "Such Slack channel/group already exists in the organization", &[]string{})

		return
	}

	SlackClient := slack.New(*Authorization.AccessToken)
	SlackChannelId, err := services.GetSlackChannelIdFromName(SlackClient, request.Channel)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if SlackChannelId == nil {
		helpers.HttpReturnErrorBadRequest(
			w,
			"Channel/group with such name was not found in Slack",
			nil,
		)

		return
	}

	err = services.JoinSlackChannel(SlackClient, *SlackChannelId)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	entity := entities.Slack{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: Authorization.OrganizationUuid,
			Name:             request.Name,
			Value:            request.Channel,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		SlackAuthorization: *Authorization,
		SlackChannelId:     SlackChannelId,
	}

	err = repositories.CreateSlackIntegration(db, &entity)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.NotifyServiceIntegrationsChanged(db, entity.Integration.OrganizationUuid)

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
