package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"

	log "github.com/sirupsen/logrus"
)

func CreateIncidentIoIntegration(db *sql.DB, entity *entities.IncidentIo) error {
	log.Tracef("Creating incident io integration")

	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Error(err.Error())

		return err
	}

	entity.Integration.Uuid = uuid.NewString()
	entity.Integration.Status = entities.IntegrationStatusOnline
	entity.Integration.Type = entities.IntegrationTypeIncidentIo
	entity.Integration.IoType = entities.IntegrationIoTypeWrite

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
			entity.Integration.Uuid,
			entity.Integration.CreatorUuid,
			entity.Integration.OrganizationUuid,
			entity.Integration.Name,
			entity.Integration.Status,
			entity.Integration.Value,
			entity.Integration.Type,
			entity.Integration.IoType,
			entity.Integration.UserUpdatedAt,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return err
		}

		entity.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO incident_ios 
        (integration_id, api_key, visibility) 
        VALUES (?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())

			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(
			entity.Integration.Id,
			entity.ApiKey.EncryptedValue,
			entity.Visibility,
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


func FindIncidentIoIntegration(db *sql.DB, Uuid string) (*entities.IncidentIo, error) {
	log.Trace("Finding incident.io integration")
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
	d.user_updated_at,
	i.id,
	i.api_key,
	i.visibility
FROM
	incident_ios i
	INNER JOIN integrations d ON d.id = i.integration_id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var entity entities.IncidentIo

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	defer stmt.Close()

	var (
		apiKey []byte
		sealedBox   kms.SecretBox
	)

	err = stmt.QueryRow(Uuid).Scan(
		&entity.Integration.Id,
		&entity.Integration.Uuid,
		&entity.Integration.CreatorUuid,
		&entity.Integration.OrganizationUuid,
		&entity.Integration.Status,
		&entity.Integration.Type,
		&entity.Integration.IoType,
		&entity.Integration.Name,
		&entity.Integration.Value,
		&entity.Integration.UserUpdatedAt,
		&entity.Id,
		&apiKey,
		&entity.Visibility,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	if len(apiKey) > 0 {
		sealedBox = kms.NewSealedSecretBox(apiKey)
	}
	entity.ApiKey = &sealedBox

	return &entity, nil
}

func UpdateIncidentIoIntegration(db *sql.DB, entity *entities.IncidentIo) error {
	Query := `UPDATE 
	incident_ios i
	INNER JOIN integrations d ON d.id = i.integration_id
       SET d.name = ?,
       d.value = ?,
       i.api_key = ?,
       i.visibility = ?,
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
		entity.Integration.Name,
		entity.Integration.Value,
		entity.ApiKey.EncryptedValue,
		entity.Visibility,
		entity.Integration.Status,
		entity.Integration.UserUpdatedAt,
		entity.Integration.Id,
	)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
