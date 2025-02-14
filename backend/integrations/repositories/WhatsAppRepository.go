package repositories

import (
	"database/sql"
	"ylem_integrations/entities"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateWhatsAppIntegration(tx *sql.Tx, WhatsAppIntegration *entities.WhatsApp) error {
	WhatsAppIntegration.Integration.Uuid = uuid.NewString()
	WhatsAppIntegration.Integration.Status = entities.IntegrationStatusNew
	WhatsAppIntegration.Integration.Type = entities.IntegrationTypeWhatsApp
	WhatsAppIntegration.Integration.IoType = entities.IntegrationIoTypeWrite

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
			WhatsAppIntegration.Integration.Uuid,
			WhatsAppIntegration.Integration.CreatorUuid,
			WhatsAppIntegration.Integration.OrganizationUuid,
			WhatsAppIntegration.Integration.Name,
			WhatsAppIntegration.Integration.Status,
			WhatsAppIntegration.Integration.Value,
			WhatsAppIntegration.Integration.Type,
			WhatsAppIntegration.Integration.IoType,
			WhatsAppIntegration.Integration.UserUpdatedAt,
		)

		if err != nil {
			log.Error(err.Error())

			return err
		}

		WhatsAppIntegration.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO whatsapps 
        (integration_id, content_sid) 
        VALUES (?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			log.Error(err.Error())

			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			WhatsAppIntegration.Integration.Id,
			WhatsAppIntegration.ContentSid,
		)

		if err != nil {
			log.Error(err.Error())

			return err
		}
	}

	return nil
}

func FindWhatsAppIntegration(db *sql.DB, Uuid string) (*entities.WhatsApp, error) {
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
	s.content_sid,
	d.user_updated_at
FROM
	whatsapps s
	INNER JOIN integrations d ON d.id = s.integration_id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var WhatsAppIntegration entities.WhatsApp

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(Uuid).Scan(
		&WhatsAppIntegration.Integration.Id,
		&WhatsAppIntegration.Integration.Uuid,
		&WhatsAppIntegration.Integration.CreatorUuid,
		&WhatsAppIntegration.Integration.OrganizationUuid,
		&WhatsAppIntegration.Integration.Status,
		&WhatsAppIntegration.Integration.Type,
		&WhatsAppIntegration.Integration.IoType,
		&WhatsAppIntegration.Integration.Name,
		&WhatsAppIntegration.Integration.Value,
		&WhatsAppIntegration.ContentSid,
		&WhatsAppIntegration.Integration.UserUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &WhatsAppIntegration, nil
}

func UpdateWhatsAppIntegration(db *sql.DB, WhatsAppIntegration *entities.WhatsApp) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = UpdateWhatsAppIntegrationTx(tx, WhatsAppIntegration)

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

func UpdateWhatsAppIntegrationTx(tx *sql.Tx, WhatsAppIntegration *entities.WhatsApp) error {
	Query := `UPDATE 
	whatsapps s
	INNER JOIN integrations d ON d.id = s.integration_id
       SET d.name = ?,
       d.value = ?,
       s.content_sid = ?,
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
		WhatsAppIntegration.Integration.Name,
		WhatsAppIntegration.Integration.Value,
		WhatsAppIntegration.ContentSid,
		WhatsAppIntegration.Integration.UserUpdatedAt,
		WhatsAppIntegration.Integration.Id,
	)

	if err != nil {
		log.Error(err.Error())

		return err
	}

	return nil
}
