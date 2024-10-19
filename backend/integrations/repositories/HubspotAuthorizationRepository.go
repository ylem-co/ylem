package repositories

import (
	"database/sql"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/services/aws/kms"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateHubspotAuthorization(db *sql.DB, Authorization *entities.HubspotAuthorization) error {
	Query := `INSERT INTO hubspot_authorizations 
        (uuid, name, creator_uuid, organization_uuid, state, is_active) 
        VALUES (?, ?, ?, ?, ?, ?)
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err.Error())

		return err
	}
	defer stmt.Close()

	Authorization.Uuid = uuid.NewString()

	result, err := stmt.Exec(
		Authorization.Uuid,
		Authorization.Name,
		Authorization.CreatorUuid,
		Authorization.OrganizationUuid,
		Authorization.State,
		Authorization.IsActive,
	)

	if err != nil {
		log.Error(err.Error())

		return err
	}

	Authorization.Id, _ = result.LastInsertId()

	return nil
}

func FindAllHubspotAuthorizationsForOrganization(db *sql.DB, OrganizationUuid string) (entities.HubspotAuthorizationCollection, error) {
	collection := entities.HubspotAuthorizationCollection{
		Items: []entities.HubspotAuthorization{},
	}
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, access_token_expires_at, refresh_token, scopes 
        FROM hubspot_authorizations WHERE organization_uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err.Error())

		return collection, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(OrganizationUuid)
	if err != nil {
		log.Error(err.Error())

		return collection, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			entity                entities.HubspotAuthorization
			accessToken           []byte
			accessTokenExpiresAt  *string
			refreshToken          []byte
			sealedAccessTokenBox  kms.SecretBox
			sealedRefreshTokenBox kms.SecretBox
		)

		err := rows.Scan(
			&entity.Id,
			&entity.Uuid,
			&entity.Name,
			&entity.CreatorUuid,
			&entity.OrganizationUuid,
			&entity.State,
			&entity.IsActive,
			&accessToken,
			&accessTokenExpiresAt,
			&refreshToken,
			&entity.Scopes,
		)

		if err != nil {
			log.Error(err.Error())

			continue
		}

		if accessTokenExpiresAt != nil {
			entity.AccessTokenExpiresAt, _ = time.Parse(helpers.DB_TIME_TIMESTAMP, *accessTokenExpiresAt)
		}

		if len(accessToken) > 0 {
			sealedAccessTokenBox = kms.NewSealedSecretBox(accessToken)
		}
		entity.AccessToken = &sealedAccessTokenBox

		if len(refreshToken) > 0 {
			sealedRefreshTokenBox = kms.NewSealedSecretBox(refreshToken)
		}
		entity.RefreshToken = &sealedRefreshTokenBox

		collection.Items = append(collection.Items, entity)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err.Error())

		return collection, err
	}

	return collection, nil
}

func FindHubspotAuthorizationByState(db *sql.DB, State string) (*entities.HubspotAuthorization, error) {
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, access_token_expires_at, refresh_token, scopes 
        FROM hubspot_authorizations WHERE state = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	defer stmt.Close()

	var (
		entity                entities.HubspotAuthorization
		accessToken           []byte
		accessTokenExpiresAt  *string
		refreshToken          []byte
		sealedAccessTokenBox  kms.SecretBox
		sealedRefreshTokenBox kms.SecretBox
	)

	err = stmt.QueryRow(State).Scan(
		&entity.Id,
		&entity.Uuid,
		&entity.Name,
		&entity.CreatorUuid,
		&entity.OrganizationUuid,
		&entity.State,
		&entity.IsActive,
		&accessToken,
		&accessTokenExpiresAt,
		&refreshToken,
		&entity.Scopes,
	)

	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	if accessTokenExpiresAt != nil {
		entity.AccessTokenExpiresAt, _ = time.Parse(helpers.DB_TIME_TIMESTAMP, *accessTokenExpiresAt)
	}

	if len(accessToken) > 0 {
		sealedAccessTokenBox = kms.NewSealedSecretBox(accessToken)
	}
	entity.AccessToken = &sealedAccessTokenBox

	if len(refreshToken) > 0 {
		sealedRefreshTokenBox = kms.NewSealedSecretBox(refreshToken)
	}
	entity.RefreshToken = &sealedRefreshTokenBox

	return &entity, nil
}

func FindHubspotAuthorizationByUuid(db *sql.DB, Uuid string) (*entities.HubspotAuthorization, error) {
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, access_token_expires_at, refresh_token, scopes 
        FROM hubspot_authorizations WHERE uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	defer stmt.Close()

	var (
		entity                entities.HubspotAuthorization
		accessToken           []byte
		accessTokenExpiresAt  *string
		refreshToken          []byte
		sealedAccessTokenBox  kms.SecretBox
		sealedRefreshTokenBox kms.SecretBox
	)

	err = stmt.QueryRow(Uuid).Scan(
		&entity.Id,
		&entity.Uuid,
		&entity.Name,
		&entity.CreatorUuid,
		&entity.OrganizationUuid,
		&entity.State,
		&entity.IsActive,
		&accessToken,
		&accessTokenExpiresAt,
		&refreshToken,
		&entity.Scopes,
	)

	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	if accessTokenExpiresAt != nil {
		entity.AccessTokenExpiresAt, err = time.Parse(time.RFC3339, *accessTokenExpiresAt)
		if err != nil {
			log.Error(err.Error())

			return nil, err
		}
	}
	if len(accessToken) > 0 {
		sealedAccessTokenBox = kms.NewSealedSecretBox(accessToken)
	}
	entity.AccessToken = &sealedAccessTokenBox

	if len(refreshToken) > 0 {
		sealedRefreshTokenBox = kms.NewSealedSecretBox(refreshToken)
	}
	entity.RefreshToken = &sealedRefreshTokenBox

	return &entity, nil
}

func UpdateHubspotAuthorization(db *sql.DB, Authorization *entities.HubspotAuthorization) error {
	Query := `UPDATE 
	hubspot_authorizations a
       SET a.name = ?,
       a.state = ?,
       a.is_active = ?,
       a.access_token = ?,
       a.access_token_expires_AT = ?,
       a.refresh_token = ?,
       a.scopes = ?
    
       WHERE a.id = ?
       `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		Authorization.Name,
		Authorization.State,
		Authorization.IsActive,
		Authorization.AccessToken.EncryptedValue,
		Authorization.AccessTokenExpiresAt.Format(helpers.DB_TIME_TIMESTAMP),
		Authorization.RefreshToken.EncryptedValue,
		Authorization.Scopes,
		Authorization.Id,
	)

	if err != nil {
		log.Error(err.Error())

		return err
	}

	return nil
}
