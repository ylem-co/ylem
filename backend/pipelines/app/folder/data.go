package folder

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

func CreateFolder(db *sql.DB, folder *Folder) (int, error) {
	folder.Uuid = uuid.NewString()
	folder.IsActive = FolderIsActive

	Query := `INSERT INTO folders
	        (uuid, organization_uuid, name, type, parent_id, is_active)
	        VALUES (?, ?, ?, ?, ?, ?)
	        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	defer stmt.Close()

	var result sql.Result
	if folder.ParentId != 0 {
		result, err = stmt.Exec(folder.Uuid, folder.OrganizationUuid, folder.Name, folder.Type, folder.ParentId, folder.IsActive)
	} else {
		result, err = stmt.Exec(folder.Uuid, folder.OrganizationUuid, folder.Name, folder.Type, nil, folder.IsActive)
	}

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

func GetFolderByUuid(db *sql.DB, uuid string) (*Folder, error) {
	Query := `SELECT 
				f1.id, 
				f1.uuid, 
				f1.name,
				f1.type, 
				f1.organization_uuid,
				f1.created_at, 
				IF (f1.updated_at IS NULL,"", f1.updated_at),
				f2.id,
				f2.uuid
              FROM folders f1
              LEFT JOIN folders f2 ON f1.parent_id = f2.id
              WHERE f1.uuid = ? AND f1.is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var folder Folder
	var parentId sql.NullInt64
	var parentUuid sql.NullString
	err = stmt.QueryRow(uuid).Scan(&folder.Id, &folder.Uuid, &folder.Name, &folder.Type, &folder.OrganizationUuid, &folder.CreatedAt, &folder.UpdatedAt, &parentId, &parentUuid)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if parentUuid.Valid {
		folder.ParentUuid = parentUuid.String
	}

	if parentId.Valid {
		folder.ParentId = parentId.Int64
	}

	return &folder, nil
}

func GetFoldersByOrgNameTypeAndParentId(db *sql.DB, name string, pId int64, fType string, orgId string) (*Folder, error) {
	var Query string
	if pId != 0 {
		Query = `SELECT 
				f1.id, 
				f1.uuid, 
				f1.name, 
				f1.type,
				f1.organization_uuid,
				f1.created_at, 
				IF (f1.updated_at IS NULL,"", f1.updated_at),
				f2.id,
				f2.uuid
              FROM folders f1
              LEFT JOIN folders f2 ON f1.parent_id = f2.id
              WHERE f1.name = ? AND f1.parent_id = ? AND f1.type = ? AND f1.organization_uuid = ? AND f1.is_active = 1
              `
    } else {
    	Query = `SELECT 
				f1.id, 
				f1.uuid, 
				f1.name,
				f1.type,
				f1.organization_uuid,
				f1.created_at, 
				IF (f1.updated_at IS NULL,"", f1.updated_at),
				f2.id,
				f2.uuid
              FROM folders f1
              LEFT JOIN folders f2 ON f1.parent_id = f2.id
              WHERE f1.name = ? AND f1.type = ? AND f1.organization_uuid = ? AND f1.parent_id IS NULL AND f1.is_active = 1
              `
    }

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var folder Folder
	var parentId sql.NullInt64
	var parentUuid sql.NullString

	if pId != 0 {
		err = stmt.QueryRow(name, pId, fType, orgId).Scan(&folder.Id, &folder.Uuid, &folder.Name, &folder.Type, &folder.OrganizationUuid, &folder.CreatedAt, &folder.UpdatedAt, &parentId, &parentUuid)
	} else {
		err = stmt.QueryRow(name, fType, orgId).Scan(&folder.Id, &folder.Uuid, &folder.Name, &folder.Type, &folder.OrganizationUuid, &folder.CreatedAt, &folder.UpdatedAt, &parentId, &parentUuid)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if parentUuid.Valid {
		folder.ParentUuid = parentUuid.String
	}

	if parentId.Valid {
		folder.ParentId = parentId.Int64
	}

	return &folder, nil
}

func GetFoldersByOrganizationUuidAndParentId(db *sql.DB, uuid string, parentId int64) (*Folders, error) {
	var Query string
	if parentId != 0 {
		Query = `SELECT 
				f1.id, 
				f1.uuid, 
				f1.name, 
				f1.type,
				f1.organization_uuid,
				f1.created_at, 
				IF (f1.updated_at IS NULL,"", f1.updated_at),
				f2.id,
				f2.uuid
              FROM folders f1
              LEFT JOIN folders f2 ON f1.parent_id = f2.id
              WHERE f1.organization_uuid = ? AND f1.is_active = 1 AND f1.parent_id = ?
              ORDER BY f1.updated_at DESC`
	} else {
		Query = `SELECT 
				f1.id, 
				f1.uuid, 
				f1.name,
				f1.type, 
				f1.organization_uuid,
				f1.created_at, 
				IF (f1.updated_at IS NULL,"", f1.updated_at),
				f2.id,
				f2.uuid
              FROM folders f1
              LEFT JOIN folders f2 ON f1.parent_id = f2.id
              WHERE f1.organization_uuid = ? AND f1.is_active = 1 AND f1.parent_id IS NULL
              ORDER BY f1.updated_at DESC`
	}

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	if parentId != 0 {
		rows, err = stmt.Query(uuid, parentId)
	} else {
		rows, err = stmt.Query(uuid)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var folders Folders

	for rows.Next() {
		var folder Folder
		var parentId sql.NullInt64
		var parentUuid sql.NullString
		err := rows.Scan(&folder.Id, &folder.Uuid, &folder.Name, &folder.Type, &folder.OrganizationUuid, &folder.CreatedAt, &folder.UpdatedAt, &parentId, &parentUuid)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if parentUuid.Valid {
			folder.ParentUuid = parentUuid.String
		}

		if parentId.Valid {
			folder.ParentId = parentId.Int64
		}

		folders.Items = append(folders.Items, folder)
	}

	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return &folders, err
	}

	return &folders, err
}

func DeleteFolder(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	defer tx.Rollback() //nolint:all

	// Delete folder
	Query := `UPDATE folders
        SET is_active = 0
        WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	// Delete pipelines
	Query = `UPDATE pipelines
        SET is_active = 0
        WHERE folder_id = ?
        `

	stmt, err = tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	// Delete tasks
	Query = `UPDATE tasks
        SET is_active = 0
        WHERE pipeline_id IN (
        	SELECT id FROM pipelines
        	WHERE folder_id = ?
        )
        `

	stmt, err = tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	// Delete tasks_triggers
	Query = `UPDATE task_triggers
        SET is_active = 0
        WHERE pipeline_id IN (
        	SELECT id FROM pipelines
        	WHERE folder_id = ?
        )
        `

	stmt, err = tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return err
	}

	return nil
}

func UpdateFolder(db *sql.DB, folder *Folder) error {
	Query := `UPDATE folders
        SET name = ?, 
        	parent_id = ?
        WHERE uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	if folder.ParentId != 0 {
		_, err = stmt.Exec(folder.Name, folder.ParentId, folder.Uuid)
	} else {
		_, err = stmt.Exec(folder.Name, nil, folder.Uuid)
	}

	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return err
	}

	return nil
}
