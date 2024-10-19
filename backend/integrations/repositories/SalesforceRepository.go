package repositories

import (
	"database/sql"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"

	log "github.com/sirupsen/logrus"
	"github.com/google/uuid"
)

func CreateSalesforceIntegration(db *sql.DB, SalesforceIntegration *entities.Salesforce) error {
	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Println(err.Error())

		return err
	}

	SalesforceIntegration.Integration.Uuid = uuid.NewString()
	SalesforceIntegration.Integration.Status = entities.IntegrationStatusOnline
	SalesforceIntegration.Integration.Type = entities.IntegrationTypeSalesforce
	SalesforceIntegration.Integration.IoType = entities.IntegrationIoTypeWrite

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
			SalesforceIntegration.Integration.Uuid,
			SalesforceIntegration.Integration.CreatorUuid,
			SalesforceIntegration.Integration.OrganizationUuid,
			SalesforceIntegration.Integration.Name,
			SalesforceIntegration.Integration.Status,
			SalesforceIntegration.Integration.Value,
			SalesforceIntegration.Integration.Type,
			SalesforceIntegration.Integration.IoType,
			SalesforceIntegration.Integration.UserUpdatedAt,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return err
		}

		SalesforceIntegration.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := `INSERT INTO salesforces 
        (integration_id, salesforce_authorization_id) 
        VALUES (?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())

			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(
			SalesforceIntegration.Integration.Id,
			SalesforceIntegration.SalesforceAuthorization.Id,
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

func UpdateSalesforceIntegration(db *sql.DB, SalesforceIntegration *entities.Salesforce) error {
	Query := `UPDATE 
	salesforces j
	INNER JOIN integrations d ON d.id = j.integration_id
       SET d.name = ?,
       d.value = ?,
       j.salesforce_authorization_id = ?,
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
		SalesforceIntegration.Integration.Name,
		SalesforceIntegration.Integration.Value,
		SalesforceIntegration.SalesforceAuthorization.Id,
		SalesforceIntegration.Integration.Status,
		SalesforceIntegration.Integration.UserUpdatedAt,
		SalesforceIntegration.Integration.Id,
	)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func FindSalesforceIntegration(db *sql.DB, Uuid string) (*entities.Salesforce, error) {
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
	a.id,
	a.uuid,
    a.name,
    a.domain,
    a.access_token,
    a.refresh_token,
    a.organization_uuid,
    a.is_active
FROM
	salesforces j
	INNER JOIN integrations d ON d.id = j.integration_id
	INNER JOIN salesforce_authorizations a ON j.salesforce_authorization_id = a.id
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var SalesforceIntegration entities.Salesforce

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	var (
		accessToken        []byte
		refreshToken       []byte
		sealedAccessToken  kms.SecretBox
		sealedRefreshToken kms.SecretBox
	)

	err = stmt.QueryRow(Uuid).Scan(
		&SalesforceIntegration.Integration.Id,
		&SalesforceIntegration.Integration.Uuid,
		&SalesforceIntegration.Integration.CreatorUuid,
		&SalesforceIntegration.Integration.OrganizationUuid,
		&SalesforceIntegration.Integration.Status,
		&SalesforceIntegration.Integration.Type,
		&SalesforceIntegration.Integration.IoType,
		&SalesforceIntegration.Integration.Name,
		&SalesforceIntegration.Integration.Value,
		&SalesforceIntegration.Integration.UserUpdatedAt,
		&SalesforceIntegration.Id,
		&SalesforceIntegration.SalesforceAuthorization.Id,
		&SalesforceIntegration.SalesforceAuthorization.Uuid,
		&SalesforceIntegration.SalesforceAuthorization.Name,
		&SalesforceIntegration.SalesforceAuthorization.Domain,
		&accessToken,
		&refreshToken,
		&SalesforceIntegration.SalesforceAuthorization.OrganizationUuid,
		&SalesforceIntegration.SalesforceAuthorization.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	if len(accessToken) > 0 {
		sealedAccessToken = kms.NewSealedSecretBox(accessToken)
	}
	SalesforceIntegration.SalesforceAuthorization.AccessToken = &sealedAccessToken

	if len(refreshToken) > 0 {
		sealedRefreshToken = kms.NewSealedSecretBox(refreshToken)
	}
	SalesforceIntegration.SalesforceAuthorization.RefreshToken = &sealedRefreshToken

	return &SalesforceIntegration, nil
}
