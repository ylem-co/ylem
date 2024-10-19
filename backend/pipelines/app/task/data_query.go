package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateQueryTx(tx *sql.Tx, reqQuery *HttpApiNewQuery) (int, error) {
	newUuid := uuid.NewString()

	var item Query
	if reqQuery != nil {
		item.SQLQuery = reqQuery.SQLQuery
		item.SourceUuid = reqQuery.SourceUuid
	}

	Query := `INSERT INTO queries
	        (uuid, sql_query, source_uuid)
	        VALUES (?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUuid, item.SQLQuery, item.SourceUuid)
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

func GetQuery(db *sql.DB, id int64) (*Query, error) {
	Req := `SELECT 
				id, 
				uuid, 
				sql_query,
				source_uuid,
				is_active, 
				created_at, 
				IF (updated_at IS NULL,"", updated_at)
              FROM queries
              WHERE id = ? AND is_active = 1
              `

	stmt, err := db.Prepare(Req)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer stmt.Close()

	var item Query

	err = stmt.QueryRow(id).Scan(&item.Id, &item.Uuid, &item.SQLQuery, &item.SourceUuid, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateQueryTx(tx *sql.Tx, id int64, item *HttpApiUpdatedQuery) error {
	Query := `UPDATE queries
        SET 
        	sql_query = ?,
        	source_uuid = ?
        WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.SQLQuery, item.SourceUuid, id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneQueryTx(tx *sql.Tx, impl *Query) (int, error) {
	req := &HttpApiNewQuery{}
	req.SQLQuery = impl.SQLQuery
	req.SourceUuid = impl.SourceUuid

	return CreateQueryTx(tx, req)
}
