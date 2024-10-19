package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateApiCallTx(tx *sql.Tx, reqApiCall *HttpApiNewApiCall) (int, error) {
	newUuid := uuid.NewString()

	var item ApiCall
	if reqApiCall != nil {
		item.Type = reqApiCall.Type
		item.Payload = reqApiCall.Payload
		item.QueryString = reqApiCall.QueryString
		item.Headers = reqApiCall.Headers
		item.AttachedFileName = reqApiCall.AttachedFileName
		item.DestinationUuid = reqApiCall.DestinationUuid
	}

	Query := `INSERT INTO api_calls
	        (uuid, type, payload, query_string, headers, attached_file_name, destination_uuid)
	        VALUES (?, ?, ?, ?, ?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUuid, item.Type, item.Payload, item.QueryString, item.Headers, item.AttachedFileName, item.DestinationUuid)
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

func GetApiCall(db *sql.DB, id int64) (*ApiCall, error) {
	Req := `SELECT 
				id, 
				uuid, 
				type,
				payload,
				query_string,
				headers,
				attached_file_name,
				destination_uuid,
				is_active, 
				created_at, 
				IF (updated_at IS NULL,"", updated_at)
              FROM api_calls
              WHERE id = ? AND is_active = 1
              `

	stmt, err := db.Prepare(Req)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var item ApiCall

	err = stmt.QueryRow(id).Scan(&item.Id, &item.Uuid, &item.Type, &item.Payload, &item.QueryString, &item.Headers, &item.AttachedFileName, &item.DestinationUuid, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &item, nil
}

func UpdateApiCallTx(tx *sql.Tx, id int64, item *HttpApiUpdatedApiCall) error {
	Query := `UPDATE api_calls
        SET 
        	type = ?,
        	payload = ?,
        	query_string = ?,
        	headers = ?,
        	attached_file_name = ?,
        	destination_uuid = ?
        WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Type, item.Payload, item.QueryString, item.Headers, item.AttachedFileName, item.DestinationUuid, id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneApiCallTx(tx *sql.Tx, impl *ApiCall) (int, error) {
	req := &HttpApiNewApiCall{}
	req.Type = impl.Type
	req.Headers = impl.Headers
	req.Payload = impl.Payload
	req.QueryString = impl.QueryString
	req.AttachedFileName = impl.AttachedFileName
	req.DestinationUuid = impl.DestinationUuid

	return CreateApiCallTx(tx, req)
}
