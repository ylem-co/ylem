package api

import (
	"encoding/json"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type GetGoogleSheetsIntegrationPrivateResponse struct {
	Integration   IntegrationPrivate `json:"integration"`
	DataKey       []byte             `json:"data_key"`
	Mode          string             `json:"mode"`
	SpreadsheetId string             `json:"spreadsheet_id"`
	SheetId       int64              `json:"sheet_id"`
	Credentials   []byte             `json:"credentials"`
	WriteHeader   bool               `json:"write_header"`
}

func GetGoogleSheetsIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
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

	dataKey, err := decryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.Credentials)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(GetGoogleSheetsIntegrationPrivateResponse{
		Integration:   IntegrationPrivate(entity.Integration),
		DataKey:       dataKey,
		Mode:          entity.Mode,
		SpreadsheetId: entity.SpreadsheetId,
		SheetId:       entity.SheetId,
		Credentials:   entity.Credentials.EncryptedValue,
		WriteHeader:   entity.WriteHeader,
	})

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
