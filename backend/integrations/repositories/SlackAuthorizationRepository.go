package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
	"ylem_integrations/entities"
)

func CreateSlackAuthorization(db *sql.DB, Authorization *entities.SlackAuthorization) error {
	Query := `INSERT INTO slack_authorizations 
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

func FindAllSlackAuthorizationsForOrganization(db *sql.DB, OrganizationUuid string) (entities.SlackAuthorizationCollection, error) {
	collection := entities.SlackAuthorizationCollection{
		Items: []entities.SlackAuthorization{},
	}
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, scopes, bot_user_id 
        FROM slack_authorizations WHERE organization_uuid = ?
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
		entity := entities.SlackAuthorization{}
		err := rows.Scan(
			&entity.Id,
			&entity.Uuid,
			&entity.Name,
			&entity.CreatorUuid,
			&entity.OrganizationUuid,
			&entity.State,
			&entity.IsActive,
			&entity.AccessToken,
			&entity.Scopes,
			&entity.BotUserId,
		)

		if err != nil {
			log.Println(err.Error())

			continue
		}

		collection.Items = append(collection.Items, entity)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err.Error())

		return collection, err
	}

	return collection, nil
}

func FindSlackAuthorizationByState(db *sql.DB, State string) (*entities.SlackAuthorization, error) {
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, scopes, bot_user_id 
        FROM slack_authorizations WHERE state = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	defer stmt.Close()

	var entity entities.SlackAuthorization

	err = stmt.QueryRow(State).Scan(
		&entity.Id,
		&entity.Uuid,
		&entity.Name,
		&entity.CreatorUuid,
		&entity.OrganizationUuid,
		&entity.State,
		&entity.IsActive,
		&entity.AccessToken,
		&entity.Scopes,
		&entity.BotUserId,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}


	return &entity, nil
}

func FindSlackAuthorizationByUuid(db *sql.DB, Uuid string) (*entities.SlackAuthorization, error) {
	Query := `SELECT id, uuid, name, creator_uuid, organization_uuid, state, is_active, access_token, scopes, bot_user_id 
        FROM slack_authorizations WHERE uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}
	defer stmt.Close()

	var entity entities.SlackAuthorization

	err = stmt.QueryRow(Uuid).Scan(
		&entity.Id,
		&entity.Uuid,
		&entity.Name,
		&entity.CreatorUuid,
		&entity.OrganizationUuid,
		&entity.State,
		&entity.IsActive,
		&entity.AccessToken,
		&entity.Scopes,
		&entity.BotUserId,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}


	return &entity, nil
}

func UpdateSlackAuthorization(db *sql.DB, Authorization *entities.SlackAuthorization) error {
	Query := `UPDATE 
	slack_authorizations a
       SET a.name = ?,
       a.state = ?,
       a.is_active = ?,
       a.access_token = ?,
       a.scopes = ?,
       a.bot_user_id = ?
    
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
		Authorization.AccessToken,
		Authorization.Scopes,
		Authorization.BotUserId,
		Authorization.Id,
	)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
