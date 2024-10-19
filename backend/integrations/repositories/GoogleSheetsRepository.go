package repositories

import (
	"database/sql"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

func CreateGoogleSheetsIntegration(db *sql.DB, gsIntegration *entities.GoogleSheets) error {
	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Error(err)

		return err
	}

	gsIntegration.Integration.Uuid = uuid.NewString()
	gsIntegration.Integration.Status = entities.IntegrationStatusOnline
	gsIntegration.Integration.Type = entities.IntegrationTypeGoogleSheets
	gsIntegration.Integration.IoType = entities.IntegrationIoTypeWrite

	{
		q := `INSERT INTO integrations 
        (uuid, creator_uuid, organization_uuid, name, status, value, type, io_type, user_updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(q)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(
			gsIntegration.Integration.Uuid,
			gsIntegration.Integration.CreatorUuid,
			gsIntegration.Integration.OrganizationUuid,
			gsIntegration.Integration.Name,
			gsIntegration.Integration.Status,
			gsIntegration.Integration.Value,
			gsIntegration.Integration.Type,
			gsIntegration.Integration.IoType,
			gsIntegration.Integration.UserUpdatedAt,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Error(err)
			return err
		}

		gsIntegration.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO google_sheets 
        			(
						integration_id,
						spreadsheet_id,
						sheet_id,
						mode,
						credentials,
						write_header
					)
        VALUES (?, ?, ?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			gsIntegration.Integration.Id,
			gsIntegration.SpreadsheetId,
			gsIntegration.SheetId,
			gsIntegration.Mode,
			gsIntegration.Credentials.EncryptedValue,
			gsIntegration.WriteHeader,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Error(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func UpdateGoogleSheetsIntegration(db *sql.DB, gsIntegration *entities.GoogleSheets) error {
	q := `UPDATE 
	google_sheets gs
	INNER JOIN integrations d ON d.id = gs.integration_id
       SET 
	   	d.name = ?,
		gs.spreadsheet_id = ?,
		gs.sheet_id = ?,
		gs.mode = ?,
		gs.credentials = ?,
		gs.write_header = ?,
		d.status = ?,
		d.user_updated_at = ?
       WHERE d.id = ?
       `

	_, err := db.Exec(
		q,
		gsIntegration.Integration.Name,
		gsIntegration.SpreadsheetId,
		gsIntegration.SheetId,
		gsIntegration.Mode,
		gsIntegration.Credentials.EncryptedValue,
		gsIntegration.WriteHeader,
		gsIntegration.Integration.Status,
		gsIntegration.Integration.UserUpdatedAt,
		gsIntegration.Integration.Id,
	)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func FindGoogleSheetsIntegration(db *sql.DB, Uuid string) (*entities.GoogleSheets, error) {
	Query := `SELECT
	d.id as integration_id,
	d.uuid,
	d.creator_uuid,
	d.organization_uuid,
	d.status,
	d.type,
	d.io_type,
	d.name,
	d.value,
	
	gs.spreadsheet_id,
	gs.sheet_id,
	gs.mode,
	gs.credentials,
	gs.write_header,

	d.user_updated_at
FROM
	google_sheets gs
	INNER JOIN integrations d ON d.id = gs.integration_id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var gsIntegration entities.GoogleSheets

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err)

		return nil, err
	}
	defer stmt.Close()

	var rawCredentials []byte

	err = stmt.QueryRow(Uuid).Scan(
		&gsIntegration.Integration.Id,
		&gsIntegration.Integration.Uuid,
		&gsIntegration.Integration.CreatorUuid,
		&gsIntegration.Integration.OrganizationUuid,
		&gsIntegration.Integration.Status,
		&gsIntegration.Integration.Type,
		&gsIntegration.Integration.IoType,
		&gsIntegration.Integration.Name,
		&gsIntegration.Integration.Value,

		&gsIntegration.SpreadsheetId,
		&gsIntegration.SheetId,
		&gsIntegration.Mode,
		&rawCredentials,
		&gsIntegration.WriteHeader,

		&gsIntegration.Integration.UserUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)

		return nil, err
	}

	if len(rawCredentials) > 0 {
		secretBox := kms.NewSealedSecretBox(rawCredentials)
		gsIntegration.Credentials = &secretBox
	}

	return &gsIntegration, nil
}
