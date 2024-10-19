package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"
)

func CreateJiraAuthorization(db *sql.DB, Authorization *entities.JiraAuthorization) error {
	Query := `INSERT INTO jira_authorizations 
        (uuid, name, creator_uuid, organization_uuid, state, is_active) 
        VALUES (?, ?, ?, ?, ?, ?)
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

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
		log.Println(err.Error())

		return err
	}

	Authorization.Id, _ = result.LastInsertId()

	return nil
}

func FindAllJiraAuthorizationsForOrganization(db *sql.DB, OrganizationUuid string) (entities.JiraAuthorizationCollection, error) {
	collection := entities.JiraAuthorizationCollection{
		Items: []entities.JiraAuthorization{},
	}
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, cloudid, scopes 
        FROM jira_authorizations WHERE organization_uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return collection, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(OrganizationUuid)
	if err != nil {
		log.Println(err.Error())

		return collection, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			entity      entities.JiraAuthorization
			accessToken []byte
			sealedBox   kms.SecretBox
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
			&entity.Cloudid,
			&entity.Scopes,
		)

		if err != nil {
			log.Println(err.Error())

			continue
		}

		if len(accessToken) > 0 {
			sealedBox = kms.NewSealedSecretBox(accessToken)
		}
		entity.AccessToken = &sealedBox

		collection.Items = append(collection.Items, entity)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err.Error())

		return collection, err
	}

	return collection, nil
}

func FindJiraAuthorizationByState(db *sql.DB, State string) (*entities.JiraAuthorization, error) {
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, cloudid, scopes 
        FROM jira_authorizations WHERE state = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	defer stmt.Close()

	var (
		entity      entities.JiraAuthorization
		accessToken []byte
		sealedBox   kms.SecretBox
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
		&entity.Cloudid,
		&entity.Scopes,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	if len(accessToken) > 0 {
		sealedBox = kms.NewSealedSecretBox(accessToken)
	}
	entity.AccessToken = &sealedBox

	return &entity, nil
}

func FindJiraAuthorizationByUuid(db *sql.DB, Uuid string) (*entities.JiraAuthorization, error) {
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, cloudid, scopes 
        FROM jira_authorizations WHERE uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	defer stmt.Close()

	var (
		entity      entities.JiraAuthorization
		accessToken []byte
		sealedBox   kms.SecretBox
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
		&entity.Cloudid,
		&entity.Scopes,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	if len(accessToken) > 0 {
		sealedBox = kms.NewSealedSecretBox(accessToken)
	}
	entity.AccessToken = &sealedBox

	return &entity, nil
}

func UpdateJiraAuthorization(db *sql.DB, Authorization *entities.JiraAuthorization) error {
	Query := `UPDATE 
	jira_authorizations a
       SET a.name = ?,
       a.state = ?,
       a.is_active = ?,
       a.access_token = ?,
       a.cloudid = ?,
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
		Authorization.Cloudid,
		Authorization.Scopes,
		Authorization.Id,
	)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
