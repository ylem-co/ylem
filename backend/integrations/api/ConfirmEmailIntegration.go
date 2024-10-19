package api

import (
	"time"
	"encoding/json"
	"net/http"
	"ylem_integrations/config"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type HttpConfirmEmailIntegration struct {
	Code  string `json:"code" valid:"type(string)"`
}

func ConfirmEmailIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpConfirmEmailIntegration

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

	db := helpers.DbConn()
	defer db.Close()

	Entity, err := repositories.FindEmailIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if Entity == nil || Entity.IsConfirmed {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		Entity.Integration.OrganizationUuid,
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	if Entity.Code != request.Code {
		helpers.HttpReturnErrorBadRequest(w, "Invalid E-Mail code", nil)

		return
	}

	Entity.IsConfirmed = true
	Entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	Entity.Integration.Status = entities.IntegrationStatusOnline

	err = repositories.UpdateEmailIntegration(db, Entity)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	var config config.Config
	err = envconfig.Process("", &config)
	if err != nil {
		log.Println(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if Entity.Integration.Value == user.Email {
		_ = services.ConfirmUsersEmail(user.Uuid);
	}

	w.WriteHeader(http.StatusNoContent)
}
