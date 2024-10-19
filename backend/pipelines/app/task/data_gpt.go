package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateGptTx(tx *sql.Tx, gpt *HttpApiNewGpt) (int, error) {
	newUuid := uuid.NewString()

	var item Gpt
	if gpt != nil {
		item.Prompt = gpt.Prompt
	}

	q := `INSERT INTO gpts
	        (uuid, prompt)
	        VALUES (?, ?)
	        `

	result, err := tx.Exec(q, newUuid, item.Prompt)
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

func GetGpt(db *sql.DB, id int64) (*Gpt, error) {
	q := `SELECT 
				id,
				uuid,
				prompt,
				is_active,
				created_at,
				IF (updated_at IS NULL, "", updated_at)
              FROM gpts
              WHERE id = ? AND is_active = 1
              `

	var item Gpt
	err := db.QueryRow(q, id).Scan(&item.Id, &item.Uuid, &item.Prompt, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateGptTx(tx *sql.Tx, id int64, item *HttpApiUpdatedGpt) error {
	q := `UPDATE gpts
        SET prompt = ?
        WHERE id = ?
        `

	_, err := tx.Exec(q, item.Prompt, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneGptTx(tx *sql.Tx, impl *Gpt) (int, error) {
	req := &HttpApiNewGpt{}
	req.Prompt = impl.Prompt

	return CreateGptTx(tx, req)
}
