package envvariable

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateEnvVariable(db *sql.DB, envVariable *EnvVariable) (int, error) {
	envVariable.Uuid = uuid.NewString()

	Query := `INSERT INTO env_variables
	        (uuid, organization_uuid, name, value)
	        VALUES (?, ?, ?, ?)
	        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(envVariable.Uuid, envVariable.OrganizationUuid, envVariable.Name, envVariable.Value)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	return int(insertID), nil
}

func UpdateEnvVariable(db *sql.DB, envVariable *EnvVariable) error {
	Query := `UPDATE env_variables
        SET name = ?, 
            value = ?
        WHERE uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(envVariable.Name, envVariable.Value, envVariable.Uuid)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return err
	}

	return nil
}

func DeleteEnvVariable(db *sql.DB, id int64) error {
	Query := `UPDATE env_variables
        SET is_active = 0, name = CONCAT(id, "_", name)
        WHERE id = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return err
	}

	return nil
}

func GetEnvVariableByUuid(db *sql.DB, uuid string) (*EnvVariable, error) {
	Query := `SELECT 
				e.id, 
				e.uuid, 
				e.name,  
				e.value,
       			e.organization_uuid,
				e.created_at, 
				IF (e.updated_at IS NULL,"", e.updated_at)
              FROM env_variables e
              WHERE e.uuid = ? AND e.is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var envVariable EnvVariable

	err = stmt.QueryRow(uuid).Scan(&envVariable.Id, &envVariable.Uuid, &envVariable.Name, &envVariable.Value, &envVariable.OrganizationUuid, &envVariable.CreatedAt, &envVariable.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &envVariable, nil
}

func GetEnvVariablesByOrganizationUuid(db *sql.DB, uuid string) (*EnvVariables, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	result, err := GetEnvVariablesByOrganizationUuidTx(tx, uuid)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetEnvVariablesByOrganizationUuidTx(tx *sql.Tx, uuid string) (*EnvVariables, error) {
	Query := `SELECT e.id, e.uuid, e.name, e.value, e.organization_uuid, e.created_at, IF (e.updated_at IS NULL,"", e.updated_at)
              FROM env_variables e
              WHERE e.organization_uuid = ? AND e.is_active = 1
              ORDER BY e.updated_at DESC`
	rows, err := tx.Query(Query, uuid)
	if err == nil {
		defer rows.Close()
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var envVariables EnvVariables

	for rows.Next() {
		var envVariable EnvVariable
		err := rows.Scan(&envVariable.Id, &envVariable.Uuid, &envVariable.Name, &envVariable.Value, &envVariable.OrganizationUuid, &envVariable.CreatedAt, &envVariable.UpdatedAt)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		envVariables.Items = append(envVariables.Items, envVariable)
	}

	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return &envVariables, err
	}

	return &envVariables, err
}

func GetEnvVariableByNameAndOrganizationUuid(db *sql.DB, name string, uuid string) (*EnvVariable, error) {
	Query := `SELECT 
				e.id, 
				e.uuid, 
				e.name,  
				e.value,
       			e.organization_uuid,
				e.created_at, 
				IF (e.updated_at IS NULL,"", e.updated_at)
              FROM env_variables e
              WHERE e.name = ? AND organization_uuid = ? AND e.is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var envVariable EnvVariable

	err = stmt.QueryRow(name, uuid).Scan(&envVariable.Id, &envVariable.Uuid, &envVariable.Name, &envVariable.Value, &envVariable.OrganizationUuid, &envVariable.CreatedAt, &envVariable.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &envVariable, nil
}
