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
	log "github.com/sirupsen/logrus"
)

type HttpCreateEmailIntegration struct {
	Name             string `json:"name" valid:"type(string)"`
	Email            string `json:"email" valid:"email"`
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4"`
}

func CreateEmailIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Creating an email Integration")

	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		log.Infof("failed to authenticate a user")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpCreateEmailIntegration

	w.Header().Set("Content-Type", "application/json")

	decodeJsonErr := helpers.DecodeJSONBody(w, r, &request)
	if decodeJsonErr != nil {
		log.Infof("failed to decode a json")

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
		log.Infof("User %s can't create email Integration on behalf of %s organization", user.Uuid, request.OrganizationUuid)
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)

		return
	}

	existingEntity, err := repositories.FindIntegrationInOrganizationByValue(db, request.Email, request.OrganizationUuid)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if existingEntity != nil {
		helpers.HttpReturnErrorBadRequest(w, "Such E-mail already exists in the organization", &[]string{})

		return
	}

	entity := entities.Email{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			Value:            request.Email,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		Code:        helpers.CreateRandomNumericString(6),
		IsConfirmed: true,
		RequestedAt: time.Now(),
	}

	err = repositories.CreateEmailIntegration(db, &entity)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.NotifyServiceIntegrationsChanged(db, entity.Integration.OrganizationUuid)

	/*_, err = services.SendEmailConfirmationEmail(entity.Integration.Value, entity.Integration.Uuid, entity.Code)
	if err != nil {
		log.Errorf(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}*/

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	log.Tracef("New email integration was created")
}
