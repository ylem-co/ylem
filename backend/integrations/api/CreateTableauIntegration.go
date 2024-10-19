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

type HttpCreateTableauIntegration struct {
	Name             string `json:"name" valid:"type(string),required"`
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4,required"`
	Server           string `json:"server" valid:"type(string),required"`
	Username         string `json:"username" valid:"type(string),required"`
	Password         string `json:"password" valid:"type(string),required"`
	Sitename         string `json:"site_name" valid:"type(string),required"`
	ProjectName      string `json:"project_name" valid:"type(string),required"`
	DatasourceName   string `json:"datasource_name" valid:"type(string),required"`
	Mode             string `json:"mode" valid:"type(string),in(overwrite|append),required"`
}

func CreateTableauIntegration(w http.ResponseWriter, r *http.Request) {
	log.Trace("Creating Tableau Integration")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("authorization"))
	if user == nil {
		log.Debug("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	log.Trace("Decode input JSON")
	w.Header().Set("Content-Type", "application/json")
	var request HttpCreateTableauIntegration

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

	log.Trace("Creating the Integration")
	username := kms.NewOpenSecretBox([]byte(request.Username))
	err = encryptSensitiveData(w, r, request.OrganizationUuid, &username)
	if err != nil {
		log.Error(err)
		return
	}

	password := kms.NewOpenSecretBox([]byte(request.Password))
	err = encryptSensitiveData(w, r, request.OrganizationUuid, &password)
	if err != nil {
		log.Error(err)
		return
	}

	entity := entities.Tableau{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			Value:            request.Server,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		Username:       &username,
		Password:       &password,
		Sitename:       request.Sitename,
		ProjectName:    request.ProjectName,
		DatasourceName: request.DatasourceName,
		Mode:           request.Mode,
	}

	err = repositories.CreateTableauIntegration(db, &entity)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.NotifyServiceIntegrationsChanged(db, entity.Integration.OrganizationUuid)

	entity.Username = nil
	entity.Password = nil

	jsonResponse, err := json.Marshal(entity)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	log.Trace("New Tableau Integration was created")
}
