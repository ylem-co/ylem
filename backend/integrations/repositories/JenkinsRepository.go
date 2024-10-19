package repositories

import (
	"database/sql"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateJenkinsIntegration(db *sql.DB, entity *entities.Jenkins) error {
	log.Tracef("Creating jenkins integration")

	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Error(err.Error())

		return err
	}

	entity.Integration.Uuid = uuid.NewString()
	entity.Integration.Status = entities.IntegrationStatusOnline
	entity.Integration.Type = entities.IntegrationTypeJenkins
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
		Query := `INSERT INTO jenkinses 
        (integration_id, base_url, token) 
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
			entity.BaseUrl,
			entity.Token.EncryptedValue,
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


func FindJenkinsIntegration(db *sql.DB, Uuid string) (*entities.Jenkins, error) {
	log.Trace("Finding jenkins io integration")
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
	i.base_url,
	i.token
FROM
	jenkinses i
	INNER JOIN integrations d ON d.id = i.integration_id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var entity entities.Jenkins

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	defer stmt.Close()

	var (
		token     []byte
		sealedBox kms.SecretBox
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
		&entity.BaseUrl,
		&token,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	if len(token) > 0 {
		sealedBox = kms.NewSealedSecretBox(token)
	}
	entity.Token = &sealedBox

	return &entity, nil
}

func UpdateJenkinsIntegration(db *sql.DB, entity *entities.Jenkins) error {
	Query := `UPDATE 
	jenkinses i
	INNER JOIN integrations d ON d.id = i.integration_id
       SET d.name = ?,
       d.value = ?,
       i.base_url = ?,
       i.token = ?,
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
		entity.BaseUrl,
		entity.Token.EncryptedValue,
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
