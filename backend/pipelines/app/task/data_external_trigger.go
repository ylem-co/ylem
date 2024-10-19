package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateExternalTriggerTx(tx *sql.Tx) (int, error) {
	newUuid := uuid.NewString()

	q := `INSERT INTO external_triggers
	        (uuid, test_data)
	        VALUES (?, ?)
	        `

	result, err := tx.Exec(q, newUuid, "")
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


func UpdateExternalTriggerTx(tx *sql.Tx, id int64, item *HttpApiUpdatedExternalTrigger) error {
	q := `UPDATE external_triggers
        SET test_data = ?
        WHERE id = ?
        `

	_, err := tx.Exec(q, item.TestData, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	return nil
}

func GetExternalTrigger(db *sql.DB, id int64) (*ExternalTrigger, error) {
	q := `SELECT 
				id,
				uuid,
       			test_data,
				is_active,
				created_at,
				IF (updated_at IS NULL, "", updated_at)
              FROM external_triggers
              WHERE id = ? AND is_active = 1
              `

	var item ExternalTrigger
	err := db.QueryRow(q, id).Scan(&item.Id, &item.Uuid, &item.TestData, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func CloneExternalTriggerTx(tx *sql.Tx) (int, error) {
	return CreateExternalTriggerTx(tx)
}
