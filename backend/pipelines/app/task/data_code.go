package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateCodeTx(tx *sql.Tx, code *HttpApiNewCode) (int, error) {
	newUuid := uuid.NewString()

	var item Code
	if code != nil {
		item.Code = code.Code
	}

	q := `INSERT INTO codes
	        (uuid, code, type)
	        VALUES (?, ?, ?)
	        `

	result, err := tx.Exec(q, newUuid, item.Code, item.Type)
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

func GetCode(db *sql.DB, id int64) (*Code, error) {
	q := `SELECT 
				id,
				uuid,
				code,
				type,
				is_active,
				created_at,
				IF (updated_at IS NULL, "", updated_at)
              FROM codes
              WHERE id = ? AND is_active = 1
              `

	var item Code
	err := db.QueryRow(q, id).Scan(&item.Id, &item.Uuid, &item.Code, &item.Type, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateCodeTx(tx *sql.Tx, id int64, item *HttpApiUpdatedCode) error {
	q := `UPDATE codes
        SET code = ?, type = ?
        WHERE id = ?
        `

	_, err := tx.Exec(q, item.Code, item.Type, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneCodeTx(tx *sql.Tx, impl *Code) (int, error) {
	req := &HttpApiNewCode{}
	req.Code = impl.Code

	return CreateCodeTx(tx, req)
}
