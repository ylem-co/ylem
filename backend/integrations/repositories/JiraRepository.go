package repositories

import (
	"database/sql"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"

	log "github.com/sirupsen/logrus"
	"github.com/google/uuid"
)

func CreateJiraIntegration(db *sql.DB, JiraIntegration *entities.Jira) error {
	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Println(err.Error())

		return err
	}

	JiraIntegration.Integration.Uuid = uuid.NewString()
	JiraIntegration.Integration.Status = entities.IntegrationStatusOnline
	JiraIntegration.Integration.Type = entities.IntegrationTypeJira
	JiraIntegration.Integration.IoType = entities.IntegrationIoTypeWrite

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
			JiraIntegration.Integration.Uuid,
			JiraIntegration.Integration.CreatorUuid,
			JiraIntegration.Integration.OrganizationUuid,
			JiraIntegration.Integration.Name,
			JiraIntegration.Integration.Status,
			JiraIntegration.Integration.Value,
			JiraIntegration.Integration.Type,
			JiraIntegration.Integration.IoType,
			JiraIntegration.Integration.UserUpdatedAt,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return err
		}

		JiraIntegration.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO jiras 
        (integration_id, jira_authorization_id, issue_type) 
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
			JiraIntegration.Integration.Id,
			JiraIntegration.JiraAuthorization.Id,
			JiraIntegration.IssueType,
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

func UpdateJiraIntegration(db *sql.DB, JiraIntegration *entities.Jira) error {
	Query := `UPDATE 
	jiras j
	INNER JOIN integrations d ON d.id = j.integration_id
       SET d.name = ?,
       d.value = ?,
       j.jira_authorization_id = ?,
       j.issue_type = ?,
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
		JiraIntegration.Integration.Name,
		JiraIntegration.Integration.Value,
		JiraIntegration.JiraAuthorization.Id,
		JiraIntegration.IssueType,
		JiraIntegration.Integration.Status,
		JiraIntegration.Integration.UserUpdatedAt,
		JiraIntegration.Integration.Id,
	)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func FindJiraIntegration(db *sql.DB, Uuid string) (*entities.Jira, error) {
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
	j.id,
	j.issue_type,
	a.id,
	a.uuid,
    a.name,
    a.cloudid,
    a.access_token,
    a.organization_uuid,
    a.is_active
FROM
	jiras j
	INNER JOIN integrations d ON d.id = j.integration_id
	INNER JOIN jira_authorizations a ON j.jira_authorization_id = a.id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var JiraIntegration entities.Jira

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	var (
		accessToken []byte
		sealedBox   kms.SecretBox
	)

	err = stmt.QueryRow(Uuid).Scan(
		&JiraIntegration.Integration.Id,
		&JiraIntegration.Integration.Uuid,
		&JiraIntegration.Integration.CreatorUuid,
		&JiraIntegration.Integration.OrganizationUuid,
		&JiraIntegration.Integration.Status,
		&JiraIntegration.Integration.Type,
		&JiraIntegration.Integration.IoType,
		&JiraIntegration.Integration.Name,
		&JiraIntegration.Integration.Value,
		&JiraIntegration.Integration.UserUpdatedAt,
		&JiraIntegration.Id,
		&JiraIntegration.IssueType,
		&JiraIntegration.JiraAuthorization.Id,
		&JiraIntegration.JiraAuthorization.Uuid,
		&JiraIntegration.JiraAuthorization.Name,
		&JiraIntegration.JiraAuthorization.Cloudid,
		&accessToken,
		&JiraIntegration.JiraAuthorization.OrganizationUuid,
		&JiraIntegration.JiraAuthorization.IsActive,
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
	JiraIntegration.JiraAuthorization.AccessToken = &sealedBox

	return &JiraIntegration, nil
}
