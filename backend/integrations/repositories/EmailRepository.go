package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"

	log "github.com/sirupsen/logrus"
)

func CreateEmailIntegration(db *sql.DB, EmailIntegration *entities.Email) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())

		return err
	}
	defer tx.Rollback() //nolint:all

	EmailIntegration.Integration.Uuid = uuid.NewString()
	EmailIntegration.Integration.Status = entities.IntegrationStatusNew
	EmailIntegration.Integration.Type = entities.IntegrationTypeEmail
	EmailIntegration.Integration.IoType = entities.IntegrationIoTypeWrite
	EmailIntegration.IsConfirmed = true

	{
		Query := `INSERT INTO integrations 
        (uuid, creator_uuid, organization_uuid, name, status, value, type, io_type, user_updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())

			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(
			EmailIntegration.Integration.Uuid,
			EmailIntegration.Integration.CreatorUuid,
			EmailIntegration.Integration.OrganizationUuid,
			EmailIntegration.Integration.Name,
			EmailIntegration.Integration.Status,
			EmailIntegration.Integration.Value,
			EmailIntegration.Integration.Type,
			EmailIntegration.Integration.IoType,
			EmailIntegration.Integration.UserUpdatedAt,
		)

		if err != nil {
			log.Println(err.Error())
			_ = tx.Rollback()

			return err
		}

		EmailIntegration.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO emails 
        (integration_id, code, is_confirmed, requested_at) 
        VALUES (?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())

			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			EmailIntegration.Integration.Id,
			EmailIntegration.Code,
			EmailIntegration.IsConfirmed,
			EmailIntegration.RequestedAt.Format(helpers.DB_TIME_TIMESTAMP),
		)

		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
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

func FindEmailIntegration(db *sql.DB, Uuid string) (*entities.Email, error) {
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
	e.is_confirmed,
	e.requested_at,
	e.code,
	d.user_updated_at
FROM
	emails e
	INNER JOIN integrations d ON d.id = e.integration_id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var EmailIntegration entities.Email

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(Uuid).Scan(
		&EmailIntegration.Integration.Id,
		&EmailIntegration.Integration.Uuid,
		&EmailIntegration.Integration.CreatorUuid,
		&EmailIntegration.Integration.OrganizationUuid,
		&EmailIntegration.Integration.Status,
		&EmailIntegration.Integration.Type,
		&EmailIntegration.Integration.IoType,
		&EmailIntegration.Integration.Name,
		&EmailIntegration.Integration.Value,
		&EmailIntegration.IsConfirmed,
		&EmailIntegration.RequestedAt,
		&EmailIntegration.Code,
		&EmailIntegration.Integration.UserUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &EmailIntegration, nil
}

func UpdateEmailIntegration(db *sql.DB, EmailIntegration *entities.Email) error {
	Query := `UPDATE 
	emails e
	INNER JOIN integrations d ON d.id = e.integration_id
       SET d.name = ?,
       d.value = ?,
       e.is_confirmed = ?,
       e.requested_at = ?,
       e.code = ?,
       d.status = ?,
	   d.user_updated_at = ?
    
       WHERE d.id = ?
       `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		EmailIntegration.Integration.Name,
		EmailIntegration.Integration.Value,
		EmailIntegration.IsConfirmed,
		EmailIntegration.RequestedAt.Format(helpers.DB_TIME_TIMESTAMP),
		EmailIntegration.Code,
		EmailIntegration.Integration.Status,
		EmailIntegration.Integration.UserUpdatedAt,
		EmailIntegration.Integration.Id,
	)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
