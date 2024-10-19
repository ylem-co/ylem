package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

type HttpConfirmSmsIntegration struct {
	Code  string `json:"code" valid:"type(string)"`
}

func ConfirmSmsIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpConfirmSmsIntegration

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

	if entity == nil || entity.IsConfirmed {
		helpers.HttpReturnErrorForbidden(w)

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
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	if request.Code != entity.Code {
		helpers.HttpReturnErrorBadRequest(
			w,
			"Invalid SMS code",
			nil,
		)

		return
	}

	entity.IsConfirmed = true
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	entity.Integration.Status = entities.IntegrationStatusOnline

	err = repositories.UpdateSmsIntegration(db, entity)
	if err != nil {
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
