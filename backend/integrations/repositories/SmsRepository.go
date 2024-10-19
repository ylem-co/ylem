package repositories

import (
	"database/sql"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateSmsIntegration(tx *sql.Tx, SmsIntegration *entities.Sms) error {
	SmsIntegration.Integration.Uuid = uuid.NewString()
	SmsIntegration.Integration.Status = entities.IntegrationStatusNew
	SmsIntegration.Integration.Type = entities.IntegrationTypeSms
	SmsIntegration.Integration.IoType = entities.IntegrationIoTypeWrite

	{
		Query := `INSERT INTO integrations 
        (uuid, creator_uuid, organization_uuid, name, status, value, type, io_type, user_updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			log.Error(err.Error())

			return err
		}
		defer stmt.Close()

		result, err := stmt.Exec(
			SmsIntegration.Integration.Uuid,
			SmsIntegration.Integration.CreatorUuid,
			SmsIntegration.Integration.OrganizationUuid,
			SmsIntegration.Integration.Name,
			SmsIntegration.Integration.Status,
			SmsIntegration.Integration.Value,
			SmsIntegration.Integration.Type,
			SmsIntegration.Integration.IoType,
			SmsIntegration.Integration.UserUpdatedAt,
		)

		if err != nil {
			log.Error(err.Error())

			return err
		}

		SmsIntegration.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO smses 
        (integration_id, code, is_confirmed, requested_at) 
        VALUES (?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			log.Error(err.Error())

			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			SmsIntegration.Integration.Id,
			SmsIntegration.Code,
			SmsIntegration.IsConfirmed,
			SmsIntegration.RequestedAt.Format(helpers.DB_TIME_TIMESTAMP),
		)

		if err != nil {
			log.Error(err.Error())

			return err
		}
	}

	return nil
}

func FindSmsIntegration(db *sql.DB, Uuid string) (*entities.Sms, error) {
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
	s.is_confirmed,
	s.requested_at,
	s.code,
	d.user_updated_at
FROM
	smses s
	INNER JOIN integrations d ON d.id = s.integration_id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var SmsIntegration entities.Sms

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(Uuid).Scan(
		&SmsIntegration.Integration.Id,
		&SmsIntegration.Integration.Uuid,
		&SmsIntegration.Integration.CreatorUuid,
		&SmsIntegration.Integration.OrganizationUuid,
		&SmsIntegration.Integration.Status,
		&SmsIntegration.Integration.Type,
		&SmsIntegration.Integration.IoType,
		&SmsIntegration.Integration.Name,
		&SmsIntegration.Integration.Value,
		&SmsIntegration.IsConfirmed,
		&SmsIntegration.RequestedAt,
		&SmsIntegration.Code,
		&SmsIntegration.Integration.UserUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &SmsIntegration, nil
}

func FindSmsIntegrationByCode(db *sql.DB, Code string, UserUuid string) (*entities.Sms, error) {
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
	s.is_confirmed,
	s.requested_at,
	s.code,
	d.user_updated_at
FROM
	smses s
	INNER JOIN integrations d ON d.id = s.integration_id
WHERE 
	s.code = ?
  	AND d.creator_uuid = ?
	AND d.deleted_at IS NULL`
	var SmsIntegration entities.Sms

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(Code, UserUuid).Scan(
		&SmsIntegration.Integration.Id,
		&SmsIntegration.Integration.Uuid,
		&SmsIntegration.Integration.CreatorUuid,
		&SmsIntegration.Integration.OrganizationUuid,
		&SmsIntegration.Integration.Status,
		&SmsIntegration.Integration.Type,
		&SmsIntegration.Integration.IoType,
		&SmsIntegration.Integration.Name,
		&SmsIntegration.Integration.Value,
		&SmsIntegration.IsConfirmed,
		&SmsIntegration.RequestedAt,
		&SmsIntegration.Code,
		&SmsIntegration.Integration.UserUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &SmsIntegration, nil
}

func UpdateSmsIntegration(db *sql.DB, SmsIntegration *entities.Sms) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = UpdateSmsIntegrationTx(tx, SmsIntegration)

	if err != nil {
		_ = tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return err
}

func UpdateSmsIntegrationTx(tx *sql.Tx, SmsIntegration *entities.Sms) error {
	Query := `UPDATE 
	smses s
	INNER JOIN integrations d ON d.id = s.integration_id
       SET d.name = ?,
       d.value = ?,
       s.is_confirmed = ?,
       s.requested_at = ?,
       s.code = ?,
       d.status = ?,
	   d.user_updated_at = ?
    
       WHERE d.id = ?
       `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		SmsIntegration.Integration.Name,
		SmsIntegration.Integration.Value,
		SmsIntegration.IsConfirmed,
		SmsIntegration.RequestedAt.Format(helpers.DB_TIME_TIMESTAMP),
		SmsIntegration.Code,
		SmsIntegration.Integration.Status,
		SmsIntegration.Integration.UserUpdatedAt,
		SmsIntegration.Integration.Id,
	)

	if err != nil {
		log.Error(err.Error())

		return err
	}

	return nil
}
