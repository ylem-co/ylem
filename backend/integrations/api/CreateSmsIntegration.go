package api

import (
	"encoding/json"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"

	"github.com/twilio/twilio-go/client"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type HttpCreateSmsIntegration struct {
	Name             string `json:"name" valid:"type(string)"`
	Number           string `json:"number" valid:"type(string)"`
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4"`
}

func CreateSmsIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpCreateSmsIntegration

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

	db := helpers.DbConn()
	defer db.Close()

	cnt, err := repositories.GetCurrentIntegrationCount(db, uuid.MustParse(request.OrganizationUuid))
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	canPerformOperation, err := services.ValidateBilledPermissions(
		user.Uuid,
		request.OrganizationUuid,
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

	existingEntity, err := repositories.FindIntegrationInOrganizationByValue(db, request.Number, request.OrganizationUuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if existingEntity != nil {
		helpers.HttpReturnErrorBadRequest(w, "Such phone number already exists in the organization", &[]string{})

		return
	}

	entity := entities.Sms{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			Value:            request.Number,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		Code:        helpers.CreateRandomNumericString(6),
		RequestedAt: time.Now(),
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}
	defer tx.Commit() //nolint:all

	err = repositories.CreateSmsIntegration(tx, &entity)

	if err != nil {
		_ = tx.Rollback()
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.NotifyServiceIntegrationsChanged(db, entity.Integration.OrganizationUuid)

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
