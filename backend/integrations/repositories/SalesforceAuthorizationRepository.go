package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"
)

func CreateSalesforceAuthorization(db *sql.DB, Authorization *entities.SalesforceAuthorization) error {
	Query := `INSERT INTO salesforce_authorizations 
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

func FindAllSalesforceAuthorizationsForOrganization(db *sql.DB, OrganizationUuid string) (entities.SalesforceAuthorizationCollection, error) {
	collection := entities.SalesforceAuthorizationCollection{
		Items: []entities.SalesforceAuthorization{},
	}
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, refresh_token, scopes, domain 
        FROM salesforce_authorizations WHERE organization_uuid = ?
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
			entity                entities.SalesforceAuthorization
			accessToken           []byte
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
			&refreshToken,
			&entity.Scopes,
			&entity.Domain,
		)

		if err != nil {
			log.Error(err.Error())

			continue
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

func FindSalesforceAuthorizationByState(db *sql.DB, State string) (*entities.SalesforceAuthorization, error) {
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, refresh_token, scopes, domain 
        FROM salesforce_authorizations WHERE state = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	defer stmt.Close()

	var (
		entity                entities.SalesforceAuthorization
		accessToken           []byte
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
		&refreshToken,
		&entity.Scopes,
		&entity.Domain,
	)

	if err != nil {
		log.Error(err.Error())

		return nil, err
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

func FindSalesforceAuthorizationByUuid(db *sql.DB, Uuid string) (*entities.SalesforceAuthorization, error) {
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, refresh_token, scopes, domain 
        FROM salesforce_authorizations WHERE uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	defer stmt.Close()

	var (
		entity                entities.SalesforceAuthorization
		accessToken           []byte
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
		&refreshToken,
		&entity.Scopes,
		&entity.Domain,
	)

	if err != nil {
		log.Error(err.Error())

		return nil, err
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

func UpdateSalesforceAuthorization(db *sql.DB, Authorization *entities.SalesforceAuthorization) error {
	Query := `UPDATE 
	salesforce_authorizations a
       SET a.name = ?,
       a.state = ?,
       a.is_active = ?,
       a.access_token = ?,
       a.refresh_token = ?,
       a.scopes = ?,
       a.domain = ?
    
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
		Authorization.RefreshToken.EncryptedValue,
		Authorization.Scopes,
		Authorization.Domain,
		Authorization.Id,
	)

	if err != nil {
		log.Error(err.Error())

		return err
	}

	return nil
}
