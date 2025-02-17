package api

import (
	"time"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type HttpUpdateWhatsAppIntegration struct {
	Name   string `json:"name" valid:"type(string)"`
	Number string `json:"number" valid:"type(string)"`
	ContentSid string `json:"content_sid" valid:"type(string)"`
}

func UpdateWhatsAppIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpUpdateWhatsAppIntegration

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

	if !entities.IsMobilePhoneValid(request.Number) {
		helpers.HttpReturnErrorBadRequest(w, "Invalid mobile phone", &[]string{"number"})

		return
	}

	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.WhatsApp
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindWhatsAppIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
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
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	entity.Integration.Name = request.Name
	entity.Integration.Value = request.Number
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	entity.ContentSid = request.ContentSid

	tx, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}
	defer tx.Commit() //nolint:all

	err = repositories.UpdateWhatsAppIntegrationTx(tx, entity)
	if err != nil {
		_ = tx.Rollback() 
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
