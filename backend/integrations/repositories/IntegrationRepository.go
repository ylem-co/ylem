package repositories

import (
	"database/sql"
	"log"
	"ylem_integrations/entities"

	"github.com/google/uuid"
)

func DeleteIntegration(db *sql.DB, Uuid string) error {
	Query := `UPDATE integrations
              SET deleted_at = NOW()
              WHERE uuid = ? AND deleted_at IS NULL
              `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		Uuid,
	)

	return err
}

func FindIntegration(db *sql.DB, Uuid string) (*entities.Integration, error) {
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
	d.user_updated_at
FROM
	integrations d 
WHERE 
	d.uuid = ?
	AND d.deleted_at IS NULL`
	var Integration entities.Integration

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(Uuid).Scan(
		&Integration.Id,
		&Integration.Uuid,
		&Integration.CreatorUuid,
		&Integration.OrganizationUuid,
		&Integration.Status,
		&Integration.Type,
		&Integration.IoType,
		&Integration.Name,
		&Integration.Value,
		&Integration.UserUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &Integration, nil
}

func FindAllIntegrationsBelongToOrganization(db *sql.DB, OrganizationUUid string, IoType string) (entities.IntegrationCollection, error) {
	var Query string

	if IoType == entities.IntegrationIoTypeAll {
		Query = `SELECT 
			d.id,
			d.uuid,
			d.creator_uuid,
			d.organization_uuid,
			d.status,
			d.type,
			d.io_type,
			d.name,
			d.value,
			d.user_updated_at,
			CASE
			    WHEN d.status = 'offline' THEN 0
		    	ELSE 1
			END AS sort_column
		FROM integrations d 
		WHERE organization_uuid = ? AND d.deleted_at IS NULL
		ORDER BY sort_column ASC, user_updated_at DESC`
	} else if IoType == entities.IntegrationIoTypeSQL {
		Query = `SELECT 
			d.id,
			d.uuid,
			d.creator_uuid,
			d.organization_uuid,
			d.status,
			d.type,
			d.io_type,
			d.name,
			d.value,
			d.user_updated_at,
			CASE
			    WHEN d.status = 'offline' THEN 0
		    	ELSE 1
			END AS sort_column
		FROM integrations d 
		WHERE organization_uuid = ? 
			AND d.deleted_at IS NULL
			AND d.type = "sql"
		ORDER BY sort_column ASC, user_updated_at DESC`
	} else {
		Query = `SELECT 
			d.id,
			d.uuid,
			d.creator_uuid,
			d.organization_uuid,
			d.status,
			d.type,
			d.io_type,
			d.name,
			d.value,
			d.user_updated_at,
			CASE
			    WHEN d.status = 'offline' THEN 0
		    	ELSE 1
			END AS sort_column
		FROM integrations d 
		WHERE organization_uuid = ? 
			AND io_type = ?
			AND d.deleted_at IS NULL
		ORDER BY sort_column ASC, user_updated_at DESC`
	}

	collection := entities.IntegrationCollection{
		Items: []entities.Integration{},
	}
	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return collection, err
	}

	defer stmt.Close()

	var rows *sql.Rows
	if IoType == entities.IntegrationIoTypeAll || IoType == entities.IntegrationIoTypeSQL  {
		rows, err = stmt.Query(OrganizationUUid)
	} else {
		rows, err = stmt.Query(OrganizationUUid, IoType)
	}

	if err != nil {
		log.Println(err.Error())

		return collection, err
	}
	defer rows.Close()

	var ignoreColumn int8
	for rows.Next() {
		entity := entities.Integration{}
		err := rows.Scan(
			&entity.Id,
			&entity.Uuid,
			&entity.CreatorUuid,
			&entity.OrganizationUuid,
			&entity.Status,
			&entity.Type,
			&entity.IoType,
			&entity.Name,
			&entity.Value,
			&entity.UserUpdatedAt,
			&ignoreColumn,
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
	}

	return collection, nil
}

func ChangeIntegrationStatus(db *sql.DB, Integration *entities.Integration, Status string) error {
	Query := `UPDATE integrations d
       SET d.status = ?
       WHERE d.id = ?
       `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	defer stmt.Close()

	Integration.Status = Status

	_, err = stmt.Exec(
		Integration.Status,
		Integration.Id,
	)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func FindIntegrationInOrganizationByValue(db *sql.DB, URL string, OrganizationUUID string) (*entities.Integration, error) {
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
    d.user_updated_at
FROM
	integrations d
WHERE 
    d.organization_uuid = ?
	AND d.value = ?
	AND d.deleted_at IS NULL`
	var Integration entities.Integration

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(OrganizationUUID, URL).Scan(
		&Integration.Id,
		&Integration.Uuid,
		&Integration.CreatorUuid,
		&Integration.OrganizationUuid,
		&Integration.Status,
		&Integration.Type,
		&Integration.IoType,
		&Integration.Name,
		&Integration.Value,
		&Integration.UserUpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	return &Integration, nil
}

func GetCurrentIntegrationCount(db *sql.DB, orgUuid uuid.UUID) (int64, error) {
	q := `SELECT COUNT(*) FROM integrations WHERE organization_uuid = ? AND deleted_at IS NULL`
	row := db.QueryRow(q, orgUuid)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var cnt int64
	err := row.Scan(&cnt)

	return cnt, err
}
