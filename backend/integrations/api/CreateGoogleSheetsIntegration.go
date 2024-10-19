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

type HttpCreateGoogleSheetsIntegration struct {
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4"`
	Name             string `json:"name" valid:"type(string)"`
	Mode             string `json:"mode" valid:"type(string),in(overwrite|append)"`
	SpreadsheetId    string `json:"spreadsheet_id" valid:"type(string)"`
	SheetId          int64  `json:"sheet_id"`
	Credentials      string `json:"credentials" valid:"type(string)"`
	WriteHeader      bool   `json:"write_header"`
}

func CreateGoogleSheetsIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpCreateGoogleSheetsIntegration

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
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	credentials := kms.NewOpenSecretBox([]byte(request.Credentials))
	err = encryptSensitiveData(w, r, request.OrganizationUuid, &credentials)
	if err != nil {
		log.Error(err)
		return
	}

	entity := entities.GoogleSheets{
		Integration: entities.Integration{
			CreatorUuid:      user.Uuid,
			OrganizationUuid: request.OrganizationUuid,
			Name:             request.Name,
			UserUpdatedAt:    time.Now().Format(helpers.DB_TIME_TIMESTAMP),
		},
		Mode:          request.Mode,
		SpreadsheetId: request.SpreadsheetId,
		SheetId:       request.SheetId,
		WriteHeader:   request.WriteHeader,
	}

	entity.Credentials = &credentials

	err = repositories.CreateGoogleSheetsIntegration(db, &entity)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.NotifyServiceIntegrationsChanged(db, entity.Integration.OrganizationUuid)

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
