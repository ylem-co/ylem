package api

import (
	"encoding/json"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type GetSlacklIntegrationPrivateResponse struct {
	Integration        IntegrationPrivate `json:"integration"`
	SlackAuthorization interface{}        `json:"authorization"`
	SlackChannelId     string             `json:"slack_channel_id"`
}

func GetSlackIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.Slack
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindSlackIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	response := GetSlacklIntegrationPrivateResponse{
		Integration: IntegrationPrivate(entity.Integration),
		SlackAuthorization: map[string]interface{}{
			"name":         entity.SlackAuthorization.Name,
			"uuid":         entity.SlackAuthorization.Uuid,
			"is_active":    entity.SlackAuthorization.IsActive,
			"access_token": entity.SlackAuthorization.AccessToken,
		},
		SlackChannelId: *entity.SlackChannelId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(response)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
