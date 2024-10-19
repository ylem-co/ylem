package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateMergeTx(tx *sql.Tx, merge *HttpApiNewMerge) (int, error) {
	newUuid := uuid.NewString()

	var item Merge
	if merge != nil {
		item.FieldNames = merge.FieldNames
	}

	q := `INSERT INTO merges
	        (uuid, field_names)
	        VALUES (?, ?)
	        `

	result, err := tx.Exec(q, newUuid, item.FieldNames)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return -1, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return -1, err
	}

	return int(insertID), nil
}

func GetMerge(db *sql.DB, id int64) (*Merge, error) {
	q := `SELECT 
				id,
				uuid,
				field_names,
				is_active,
				created_at,
				IF (updated_at IS NULL, "", updated_at)
              FROM merges
              WHERE id = ? AND is_active = 1
              `

	var item Merge
	err := db.QueryRow(q, id).Scan(&item.Id, &item.Uuid, &item.FieldNames, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateMergeTx(tx *sql.Tx, id int64, item *HttpApiUpdatedMerge) error {
	q := `UPDATE merges
        SET field_names = ?
        WHERE id = ?
        `

	_, err := tx.Exec(q, item.FieldNames, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneMergeTx(tx *sql.Tx, impl *Merge) (int, error) {
	req := &HttpApiNewMerge{}
	req.FieldNames = impl.FieldNames

	return CreateMergeTx(tx, req)
}
