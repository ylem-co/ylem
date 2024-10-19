package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateProcessorTx(tx *sql.Tx, reqProcessor *HttpApiNewProcessor) (int, error) {
	var err error
	newUuid := uuid.NewString()

	var item Processor
	if reqProcessor != nil {
		item.Expression = reqProcessor.Expression
		item.Strategy = reqProcessor.Strategy
	}

	Query := `INSERT INTO processors
	        (uuid, expression, strategy)
	        VALUES (?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUuid, item.Expression, item.Strategy)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return -1, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return -1, err
	}

	return int(insertID), nil
}

func GetProcessor(db *sql.DB, id int64) (*Processor, error) {
	Query := `SELECT 
				id, 
				uuid, 
				expression, 
				strategy,
				is_active, 
				created_at, 
				IF (updated_at IS NULL,"", updated_at)
              FROM processors
              WHERE id = ? AND is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer stmt.Close()

	var item Processor

	err = stmt.QueryRow(id).Scan(&item.Id, &item.Uuid, &item.Expression, &item.Strategy, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateProcessorTx(tx *sql.Tx, id int64, item *HttpApiUpdatedProcessor) error {
	Query := `UPDATE processors
		SET expression = ?, strategy = ?
		WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Expression, item.Strategy, id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneProcessorTx(tx *sql.Tx, impl *Processor) (int, error) {
	req := &HttpApiNewProcessor{}
	req.Expression = impl.Expression
	req.Strategy = impl.Strategy

	return CreateProcessorTx(tx, req)
}
