package task

import (
	"fmt"
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateAggregatorTx(tx *sql.Tx, reqAggregator *HttpApiNewAggregator, wfUuid string) (int, error) {
	var err error
	newUuid := uuid.NewString()

	var item Aggregator
	if reqAggregator != nil {
		item.Expression = reqAggregator.Expression
		item.VariableName = reqAggregator.VariableName
	}

	if item.VariableName == "" {
		item.VariableName, err = generateAggregatorVariableName(tx, wfUuid)
		if err != nil {
			return 0, err
		}
	}

	Query := `INSERT INTO aggregators
	        (uuid, expression, variable_name)
	        VALUES (?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUuid, item.Expression, item.VariableName)
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

func GetAggregator(db *sql.DB, id int64) (*Aggregator, error) {
	Query := `SELECT 
				id, 
				uuid, 
				expression, 
				variable_name,
				is_active, 
				created_at, 
				IF (updated_at IS NULL,"", updated_at)
              FROM aggregators
              WHERE id = ? AND is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer stmt.Close()

	var item Aggregator

	err = stmt.QueryRow(id).Scan(&item.Id, &item.Uuid, &item.Expression, &item.VariableName, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateAggregatorTx(tx *sql.Tx, id int64, item *HttpApiUpdatedAggregator) error {
	Query := `UPDATE aggregators
		SET expression = ?, variable_name = ?
		WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Expression, item.VariableName, id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneAggregatorTx(tx *sql.Tx, impl *Aggregator, wfUuid string) (int, error) {
	req := &HttpApiNewAggregator{}
	req.Expression = impl.Expression

	return CreateAggregatorTx(tx, req, wfUuid)
}

func getAggregatorCount(tx *sql.Tx, wfUuid string) (int64, error) {
	q := `SELECT COUNT(*) 
		FROM 
			aggregators a 
			INNER JOIN tasks t ON t.implementation_id = a.id AND t.type = ?
			INNER JOIN pipelines w ON w.id = t.pipeline_id
		WHERE w.uuid = ?`

	var cnt int64
	err := tx.QueryRow(q, TaskTypeAggregator, wfUuid).Scan(&cnt)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func generateAggregatorVariableName(tx *sql.Tx, wfUuid string) (string, error) {
	aggCount, err := getAggregatorCount(tx, wfUuid)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("aggregated_%d", aggCount+1), nil
}
