package repositories

import (
	"database/sql"
	"ylem_integrations/entities"

	log "github.com/sirupsen/logrus"
	"github.com/google/uuid"
)

func CreateSlackIntegration(db *sql.DB, SlackIntegration *entities.Slack) error {
	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Println(err.Error())

		return err
	}

	SlackIntegration.Integration.Uuid = uuid.NewString()
	SlackIntegration.Integration.Status = entities.IntegrationStatusOnline
	SlackIntegration.Integration.Type = entities.IntegrationTypeSlack
	SlackIntegration.Integration.IoType = entities.IntegrationIoTypeWrite

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
			SlackIntegration.Integration.Uuid,
			SlackIntegration.Integration.CreatorUuid,
			SlackIntegration.Integration.OrganizationUuid,
			SlackIntegration.Integration.Name,
			SlackIntegration.Integration.Status,
			SlackIntegration.Integration.Value,
			SlackIntegration.Integration.Type,
			SlackIntegration.Integration.IoType,
			SlackIntegration.Integration.UserUpdatedAt,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return err
		}

		SlackIntegration.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO slacks 
        (integration_id, slack_authorization_id, slack_channel_id) 
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
			SlackIntegration.Integration.Id,
			SlackIntegration.SlackAuthorization.Id,
			SlackIntegration.SlackChannelId,
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

func UpdateSlackIntegration(db *sql.DB, SlackIntegration *entities.Slack) error {
	Query := `UPDATE 
	slacks s
	INNER JOIN integrations d ON d.id = s.integration_id
       SET d.name = ?,
       d.value = ?,
       s.slack_authorization_id = ?,
       s.slack_channel_id = ?,
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
		SlackIntegration.Integration.Name,
		SlackIntegration.Integration.Value,
		SlackIntegration.SlackAuthorization.Id,
		SlackIntegration.SlackChannelId,
		SlackIntegration.Integration.Status,
		SlackIntegration.Integration.UserUpdatedAt,
		SlackIntegration.Integration.Id,
	)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func FindSlackIntegration(db *sql.DB, Uuid string) (*entities.Slack, error) {
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
	s.id,
	s.slack_channel_id,
	a.id,
	a.uuid,
    a.name,
    a.access_token,
    a.organization_uuid
FROM
	slacks s
	INNER JOIN integrations d ON d.id = s.integration_id
	INNER JOIN slack_authorizations a ON s.slack_authorization_id = a.id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var SlackIntegration entities.Slack

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(Uuid).Scan(
		&SlackIntegration.Integration.Id,
		&SlackIntegration.Integration.Uuid,
		&SlackIntegration.Integration.CreatorUuid,
		&SlackIntegration.Integration.OrganizationUuid,
		&SlackIntegration.Integration.Status,
		&SlackIntegration.Integration.Type,
		&SlackIntegration.Integration.IoType,
		&SlackIntegration.Integration.Name,
		&SlackIntegration.Integration.Value,
		&SlackIntegration.Integration.UserUpdatedAt,
		&SlackIntegration.Id,
		&SlackIntegration.SlackChannelId,
		&SlackIntegration.SlackAuthorization.Id,
		&SlackIntegration.SlackAuthorization.Uuid,
		&SlackIntegration.SlackAuthorization.Name,
		&SlackIntegration.SlackAuthorization.AccessToken,
		&SlackIntegration.SlackAuthorization.OrganizationUuid,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &SlackIntegration, nil
}
