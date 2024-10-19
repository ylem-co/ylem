package entities

import "ylem_integrations/services/aws/kms"

const (
	IntegrationTypeGoogleSheets = "google-sheets"

	GoogleSheetsModeOverwrite = "overwrite"
	GoogleSheetsModeAppend    = "append"
)

type GoogleSheets struct {
	Id            int64          `json:"-"`
	Integration   Integration    `json:"integration"`
	SpreadsheetId string         `json:"spreadsheet_id"`
	SheetId       int64          `json:"sheet_id"`
	Mode          string         `json:"mode"`
	Credentials   *kms.SecretBox `json:"credentials"`
	WriteHeader   bool           `json:"write_header"`
}
