package repositories

import (
	"database/sql"
	"ylem_integrations/entities"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

func CreateApiIntegration(db *sql.DB, ApiIntegration *entities.Api) error {
	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Error(err)

		return err
	}

	ApiIntegration.Integration.Uuid = uuid.NewString()
	ApiIntegration.Integration.Status = entities.IntegrationStatusOnline
	ApiIntegration.Integration.Type = entities.IntegrationTypeApi
	ApiIntegration.Integration.IoType = entities.IntegrationIoTypeReadWrite

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
			ApiIntegration.Integration.Uuid,
			ApiIntegration.Integration.CreatorUuid,
			ApiIntegration.Integration.OrganizationUuid,
			ApiIntegration.Integration.Name,
			ApiIntegration.Integration.Status,
			ApiIntegration.Integration.Value,
			ApiIntegration.Integration.Type,
			ApiIntegration.Integration.IoType,
			ApiIntegration.Integration.UserUpdatedAt,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Error(err)
			return err
		}

		ApiIntegration.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO apis 
        (integration_id, method, auth_type, auth_bearer_token, auth_basic_user_name, auth_basic_user_password, auth_header_name, auth_header_value) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			ApiIntegration.Integration.Id,
			ApiIntegration.Method,
			ApiIntegration.AuthType,
			ApiIntegration.AuthBearerToken,
			ApiIntegration.AuthBasicUserName,
			ApiIntegration.AuthBasicUserPassword,
			ApiIntegration.AuthHeaderName,
			ApiIntegration.AuthHeaderValue,
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

func UpdateApiIntegration(db *sql.DB, ApiIntegration *entities.Api) error {
	Query := `UPDATE 
	apis a
	INNER JOIN integrations d ON d.id = a.integration_id
       SET d.name = ?,
       d.value = ?,
	   a.method = ?,
       a.auth_type = ?,
       a.auth_bearer_token = ?,
       a.auth_basic_user_name = ?,
       a.auth_basic_user_password = ?,
       a.auth_header_name = ?,
       a.auth_header_value = ?,
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
		ApiIntegration.Integration.Name,
		ApiIntegration.Integration.Value,
		ApiIntegration.Method,
		ApiIntegration.AuthType,
		ApiIntegration.AuthBearerToken,
		ApiIntegration.AuthBasicUserName,
		ApiIntegration.AuthBasicUserPassword,
		ApiIntegration.AuthHeaderName,
		ApiIntegration.AuthHeaderValue,
		ApiIntegration.Integration.Status,
		ApiIntegration.Integration.UserUpdatedAt,
		ApiIntegration.Integration.Id,
	)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func FindApiIntegration(db *sql.DB, Uuid string) (*entities.Api, error) {
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
	a.method,
	a.auth_type,
	a.auth_bearer_token,
	a.auth_basic_user_name,
	a.auth_basic_user_password,
	a.auth_header_name,
	a.auth_header_value,
	d.user_updated_at
FROM
	apis a
	INNER JOIN integrations d ON d.id = a.integration_id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var ApiIntegration entities.Api

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(Uuid).Scan(
		&ApiIntegration.Integration.Id,
		&ApiIntegration.Integration.Uuid,
		&ApiIntegration.Integration.CreatorUuid,
		&ApiIntegration.Integration.OrganizationUuid,
		&ApiIntegration.Integration.Status,
		&ApiIntegration.Integration.Type,
		&ApiIntegration.Integration.IoType,
		&ApiIntegration.Integration.Name,
		&ApiIntegration.Integration.Value,
		&ApiIntegration.Method,
		&ApiIntegration.AuthType,
		&ApiIntegration.AuthBearerToken,
		&ApiIntegration.AuthBasicUserName,
		&ApiIntegration.AuthBasicUserPassword,
		&ApiIntegration.AuthHeaderName,
		&ApiIntegration.AuthHeaderValue,
		&ApiIntegration.Integration.UserUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)

		return nil, err
	}

	return &ApiIntegration, nil
}
