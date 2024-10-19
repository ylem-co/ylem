package repositories

import (
	"database/sql"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateTableauIntegration(db *sql.DB, entity *entities.Tableau) error {
	log.Tracef("Creating Tableau integration")

	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Error(err)

		return err
	}

	entity.Integration.Uuid = uuid.NewString()
	entity.Integration.Status = entities.IntegrationStatusOnline
	entity.Integration.Type = entities.IntegrationTypeTableau
	entity.Integration.IoType = entities.IntegrationIoTypeWrite

	{
		Query := `INSERT INTO integrations 
        (uuid, creator_uuid, organization_uuid, name, status, value, type, io_type, user_updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

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
			log.Error(err)
			return err
		}

		entity.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO tableau 
        (integration_id, username, password, site_name, project_name, datasource_name, mode) 
        VALUES (?, ?, ?, ?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			entity.Integration.Id,
			entity.Username.EncryptedValue,
			entity.Password.EncryptedValue,
			entity.Sitename,
			entity.ProjectName,
			entity.DatasourceName,
			entity.Mode,
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

func FindTableauIntegration(db *sql.DB, Uuid string) (*entities.Tableau, error) {
	log.Trace("Finding Tableau integration")
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
	t.id,
	t.username,
	t.password,
	t.site_name,
	t.project_name,
	t.datasource_name,
	t.mode
FROM
	tableau t
	INNER JOIN integrations d ON d.id = t.integration_id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var entity entities.Tableau

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err)

		return nil, err
	}
	defer stmt.Close()

	var (
		username          []byte
		usernameSealedBox kms.SecretBox
		password          []byte
		passwordSealedBox kms.SecretBox
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
		&username,
		&password,
		&entity.Sitename,
		&entity.ProjectName,
		&entity.DatasourceName,
		&entity.Mode,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)

		return nil, err
	}

	if len(username) > 0 {
		usernameSealedBox = kms.NewSealedSecretBox(username)
	}
	entity.Username = &usernameSealedBox

	if len(password) > 0 {
		passwordSealedBox = kms.NewSealedSecretBox(password)
	}
	entity.Password = &passwordSealedBox

	return &entity, nil
}

func UpdateTableauIntegration(db *sql.DB, entity *entities.Tableau) error {
	Query := `UPDATE 
	tableau t
	INNER JOIN integrations d ON d.id = t.integration_id
       SET d.name = ?,
       d.value = ?,
       t.mode = ?,
       t.username = ?,
	   t.password = ?,
	   t.site_name = ?,
	   t.project_name = ?,
	   t.datasource_name = ?,
       d.status = ?,
	   d.user_updated_at = ?
       WHERE d.id = ?
       `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err)

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		entity.Integration.Name,
		entity.Integration.Value,
		entity.Mode,
		entity.Username.EncryptedValue,
		entity.Password.EncryptedValue,
		entity.Sitename,
		entity.ProjectName,
		entity.DatasourceName,
		entity.Integration.Status,
		entity.Integration.UserUpdatedAt,
		entity.Integration.Id,
	)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
