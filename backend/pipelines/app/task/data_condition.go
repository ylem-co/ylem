package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateConditionTx(tx *sql.Tx, reqCondition *HttpApiNewCondition) (int, error) {
	newUuid := uuid.NewString()

	var item Condition
	if reqCondition != nil {
		item.Expression = reqCondition.Expression
	}

	Query := `INSERT INTO conditions
	        (uuid, expression)
	        VALUES (?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUuid, item.Expression)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return -1, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return -1, err
	}

	return int(insertID), nil
}

func GetCondition(db *sql.DB, id int64) (*Condition, error) {
	Query := `SELECT 
				id, 
				uuid, 
				expression, 
				is_active, 
				created_at, 
				IF (updated_at IS NULL,"", updated_at)
              FROM conditions
              WHERE id = ? AND is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var item Condition

	err = stmt.QueryRow(id).Scan(&item.Id, &item.Uuid, &item.Expression, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &item, nil
}

func UpdateConditionTx(tx *sql.Tx, id int64, item *HttpApiUpdatedCondition) error {
	Query := `UPDATE conditions
        SET expression = ?
        WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Expression, id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneConditionTx(tx *sql.Tx, impl *Condition) (int, error) {
	req := &HttpApiNewCondition{}
	req.Expression = impl.Expression

	return CreateConditionTx(tx, req)
}
