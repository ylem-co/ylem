package api

import (
	"encoding/json"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	"ylem_integrations/services/aws/kms"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type HttpCreateJenkinsIntegration struct {
	Name             string `json:"name" valid:"type(string),required"`
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4,required"`
	BaseUrl          string `json:"base_url" valid:"type(string),required"`
	Token            string `json:"token" valid:"type(string),required"`
	ProjectName      string `json:"project_name" valid:"type(string),required"`
}

func CreateJenkinsIntegration(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Creating Jenkins Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("authorization"))
	if user == nil {
		log.Debugf("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Tracef("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpCreateJenkinsIntegration

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
		log.Debugf(
			"User %s can't perform the operation %s in %s",
			user.Uuid,
			services.PermissionActionCreate,
			services.PermissionResourceTypeIntegration,
		)
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)

		return
	}

	log.Tracef("Creating the Integration")
	secretToken := kms.NewOpenSecretBox([]byte(request.Token))
	err = encryptSensitiveData(w, r, request.OrganizationUuid, &secretToken)
	if err != nil {
		return
	}

	entity := entities.Jenkins{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			Value:            request.ProjectName,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		BaseUrl: request.BaseUrl,
		Token:   &secretToken,
	}

	err = repositories.CreateJenkinsIntegration(db, &entity)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.NotifyServiceIntegrationsChanged(db, entity.Integration.OrganizationUuid)

	entity.Token = nil

	jsonResponse, err := json.Marshal(entity)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	log.Tracef("New jenkins Integration was created")
}
