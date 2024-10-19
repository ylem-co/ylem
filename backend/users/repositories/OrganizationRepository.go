package repositories

import (
	"database/sql"
	"ylem_users/entities"
	"ylem_users/helpers"

	log "github.com/sirupsen/logrus"
)

func DoesOrganizationExist(db *sql.DB, organizationName string, organizationUuid string) bool {
	Query := ""
	var rows *sql.Rows
	if organizationUuid == "" {
		Query = `SELECT COUNT(*)
              FROM organizations
              WHERE name = ?
              `
		rows, _ = db.Query(Query, organizationName)
	} else {
		Query = `SELECT COUNT(*)
              FROM organizations
              WHERE name = ? AND uuid != ?
              `
		rows, _ = db.Query(Query, organizationName, organizationUuid)
	}

	return helpers.NumRows(rows) > 0
}

func GetOrganization(db *sql.DB)  (entities.Organization, bool) {
	Query := `SELECT
				id, 
				uuid,
                name, 
                is_data_source_created, 
                is_destination_created,
                is_pipeline_created,
                data_key
        	FROM organizations
            WHERE is_active = 1
            `
	var org entities.Organization
	stmt, err := db.Prepare(Query)
	if err != nil {
		return org, false
	}

	err = stmt.QueryRow().Scan(
		&org.Id,
		&org.Uuid,
		&org.Name,
		&org.IsDataSourceCreated,
		&org.IsDestinationCreated,
		&org.IsPipelineCreated,
		&org.DataKey,
	)

	if err != nil {
		log.Error(err)
		return org, false
	} else {
		return org, true
	}
}

func GetOrganizationByUuid(db *sql.DB, uuid string) (entities.Organization, bool) {
	Query := `SELECT 
                id, 
				uuid,
                name, 
                is_data_source_created, 
                is_destination_created,
                is_pipeline_created,
                data_key
              FROM organizations
              WHERE uuid = ?
              `
	var org entities.Organization
	stmt, err := db.Prepare(Query)
	if err != nil {
		return org, false
	}

	err = stmt.QueryRow(uuid).Scan(
		&org.Id,
		&org.Uuid,
		&org.Name,
		&org.IsDataSourceCreated,
		&org.IsDestinationCreated,
		&org.IsPipelineCreated,
		&org.DataKey,
	)

	if err != nil {
		log.Error(err)
		return org, false
	} else {
		return org, true
	}
}

func GetOrganizationByUserUuid(db *sql.DB, uuid string) (entities.Organization, bool) {
	Query := `SELECT 
                o.name, 
                o.uuid,
                o.is_data_source_created,
                o.is_destination_created,
                o.is_pipeline_created,
                o.data_key
              FROM organizations o
              JOIN users u ON o.id = u.organization_id
              WHERE u.uuid = ?
              `
	var org entities.Organization
	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err)
		return org, false
	}
	err = stmt.QueryRow(uuid).Scan(&org.Name, &org.Uuid, &org.IsDataSourceCreated, &org.IsDestinationCreated, &org.IsPipelineCreated, &org.DataKey)

	if err != nil {
		log.Error(err)
		return org, false
	} else {
		return org, true
	}
}

func UpdateConnections(db *sql.DB, org entities.Organization) bool {
	updateQuery := `UPDATE organizations 
            SET 
                is_data_source_created = ?,
                is_destination_created = ?,
                is_pipeline_created = ?
            WHERE id = ?
            `

	updateStatement, err := db.Prepare(updateQuery)
	if err != nil {
		log.Error(err)
		return false
	}
	defer updateStatement.Close()

	_, err = updateStatement.Exec(org.IsDataSourceCreated, org.IsDestinationCreated, org.IsPipelineCreated, org.Id)
	if err != nil {
		log.Error(err)
		return false
	}

	return true
}
