package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateFilterTx(tx *sql.Tx, filter *HttpApiNewFilter) (int, error) {
	newUuid := uuid.NewString()

	var item Filter
	if filter != nil {
		item.Expression = filter.Expression
	}

	q := `INSERT INTO filters
	        (uuid, expression)
	        VALUES (?, ?)
	        `

	result, err := tx.Exec(q, newUuid, item.Expression)
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

func GetFilter(db *sql.DB, id int64) (*Filter, error) {
	q := `SELECT 
				id,
				uuid,
				expression,
				is_active,
				created_at,
				IF (updated_at IS NULL, "", updated_at)
              FROM filters
              WHERE id = ? AND is_active = 1
              `

	var item Filter
	err := db.QueryRow(q, id).Scan(&item.Id, &item.Uuid, &item.Expression, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateFilterTx(tx *sql.Tx, id int64, item *HttpApiUpdatedFilter) error {
	q := `UPDATE filters
        SET expression = ?
        WHERE id = ?
        `

	_, err := tx.Exec(q, item.Expression, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneFilterTx(tx *sql.Tx, impl *Filter) (int, error) {
	req := &HttpApiNewFilter{}
	req.Expression = impl.Expression

	return CreateFilterTx(tx, req)
}
