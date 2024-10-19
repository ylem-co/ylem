package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateForEachTx(tx *sql.Tx, reqForEach *HttpApiNewForEach) (int, error) {
	newUuid := uuid.NewString()

	Query := `INSERT INTO for_eaches
	        (uuid)
	        VALUES (?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUuid)
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

func GetForEach(db *sql.DB, id int64) (*ForEach, error) {
	Query := `SELECT 
				id, 
				uuid,
				is_active, 
				created_at, 
				IF (updated_at IS NULL,"", updated_at)
              FROM for_eaches
              WHERE id = ? AND is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var item ForEach

	err = stmt.QueryRow(id).Scan(&item.Id, &item.Uuid, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &item, nil
}

func CloneForEachTx(tx *sql.Tx, impl *ForEach) (int, error) {
	return CreateForEachTx(tx, &HttpApiNewForEach{})
}
