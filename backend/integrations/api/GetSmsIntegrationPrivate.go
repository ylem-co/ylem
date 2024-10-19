package api

import (
	"encoding/json"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type GetSmsIntegrationPrivateResponse struct {
	Integration IntegrationPrivate `json:"integration"`
	IsConfirmed bool               `json:"is_confirmed"`
	RequestedAt time.Time          `json:"requested_at"`
}

func GetSmsIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.Sms
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindSmsIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	response := GetEmailIntegrationPrivateResponse{
		Integration: IntegrationPrivate(entity.Integration),
		IsConfirmed: entity.IsConfirmed,
		RequestedAt: entity.RequestedAt,
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
