package repositories

import (
	"time"
	"database/sql"
	"github.com/google/uuid"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"

	log "github.com/sirupsen/logrus"
)

func CreateHubspotIntegration(db *sql.DB, entity *entities.Hubspot) error {
	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Println(err.Error())

		return err
	}

	entity.Integration.Uuid = uuid.NewString()
	entity.Integration.Status = entities.IntegrationStatusOnline
	entity.Integration.Type = entities.IntegrationTypeHubspot
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
		Query := `INSERT INTO hubspots 
        (integration_id, hubspot_authorization_id, pipeline_stage_code, owner_code) 
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
			entity.Integration.Id,
			entity.HubspotAuthorization.Id,
			entity.PipelineStageCode,
			entity.OwnerCode,
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

func UpdateHubspotIntegration(db *sql.DB, entity *entities.Hubspot) error {
	Query := `UPDATE 
	hubspots h
	INNER JOIN integrations d ON d.id = h.integration_id
       SET d.name = ?,
       d.value = ?,
       h.hubspot_authorization_id = ?,
       h.pipeline_stage_code = ?,
       h.owner_code = ?,
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
		entity.HubspotAuthorization.Id,
		entity.PipelineStageCode,
		entity.OwnerCode,
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

func FindHubspotIntegration(db *sql.DB, Uuid string) (*entities.Hubspot, error) {
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
	h.id,
	h.pipeline_stage_code,
	h.owner_code,
	a.id,
	a.uuid,
    a.name,
    a.access_token,
    a.access_token_expires_at,
    a.refresh_token,
    a.organization_uuid,
    a.is_active
FROM
	hubspots h
	INNER JOIN integrations d ON d.id = h.integration_id
	INNER JOIN hubspot_authorizations a ON h.hubspot_authorization_id = a.id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var entity entities.Hubspot

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	var (
		accessToken           []byte
		accessTokenExpiresAt  *string
		refreshToken          []byte
		sealedAccessTokenBox  kms.SecretBox
		sealedRefreshTokenBox kms.SecretBox
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
		&entity.PipelineStageCode,
		&entity.OwnerCode,
		&entity.HubspotAuthorization.Id,
		&entity.HubspotAuthorization.Uuid,
		&entity.HubspotAuthorization.Name,
		&accessToken,
		&accessTokenExpiresAt,
		&refreshToken,
		&entity.HubspotAuthorization.OrganizationUuid,
		&entity.HubspotAuthorization.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	if accessTokenExpiresAt != nil {
		entity.HubspotAuthorization.AccessTokenExpiresAt, _ = time.Parse(time.RFC3339, *accessTokenExpiresAt)
	}

	if len(accessToken) > 0 {
		sealedAccessTokenBox = kms.NewSealedSecretBox(accessToken)
	}
	entity.HubspotAuthorization.AccessToken = &sealedAccessTokenBox

	if len(refreshToken) > 0 {
		sealedRefreshTokenBox = kms.NewSealedSecretBox(refreshToken)
	}
	entity.HubspotAuthorization.RefreshToken = &sealedRefreshTokenBox

	return &entity, nil
}
