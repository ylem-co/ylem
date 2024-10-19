package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go/client"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

func ResendConfirmationSms(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	w.Header().Set("Content-Type", "application/json")

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
		services.PermissionActionUpdate, // update?
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	if !entity.CanResendSms() {
		helpers.HttpReturnErrorBadRequest(
			w,
			"Can't resend a sms <@todo: when next time can do>",
			nil,
		)

		return
	}

	entity.Code = helpers.CreateRandomNumericString(6)
	entity.RequestedAt = time.Now()
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	tx, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}
	defer tx.Commit() //nolint:all

	err = repositories.UpdateSmsIntegrationTx(tx, entity)
	if err != nil {
		_ = tx.Rollback()
		helpers.HttpReturnErrorInternal(w)

		return
	}

	err = services.SendPhoneNumberVerificationSms(entity.Integration.Value, entity.Code)
	if err != nil {
		log.Error(err.Error())

		switch e := err.(type) {
		case *client.TwilioRestError:
			if e.Status == http.StatusBadRequest {
				helpers.HttpReturnErrorBadRequest(w, "Invalid mobile phone", &[]string{"number"})

				_ = tx.Rollback()
			} else {
				helpers.HttpReturnErrorServiceUnavailable(w)
			}
		default:
			helpers.HttpReturnErrorInternal(w)
		}

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
