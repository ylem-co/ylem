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

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type HttpUpdateGoogleSheetsIntegration struct {
	Name          string `json:"name" valid:"type(string)"`
	Mode          string `json:"mode" valid:"type(string),in(overwrite|append)"`
	SpreadsheetId string `json:"spreadsheet_id" valid:"type(string)"`
	SheetId       int64  `json:"sheet_id"`
	Credentials   string `json:"credentials" valid:"type(string)"`
	WriteHeader   bool   `json:"write_header"`
}

func UpdateGoogleSheetsIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpUpdateGoogleSheetsIntegration

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

	var entity *entities.GoogleSheets
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindGoogleSheetsIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
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

	credentials := kms.NewOpenSecretBox([]byte(request.Credentials))
	err = encryptSensitiveData(w, r, entity.Integration.OrganizationUuid, &credentials)
	if err != nil {
		return
	}

	entity.Integration.Name = request.Name
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	entity.Mode = request.Mode
	entity.SpreadsheetId = request.SpreadsheetId
	entity.SheetId = request.SheetId
	entity.Credentials = &credentials
	entity.WriteHeader = request.WriteHeader

	err = repositories.UpdateGoogleSheetsIntegration(db, entity)
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
