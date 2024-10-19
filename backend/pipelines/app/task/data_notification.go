package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateNotificationTx(tx *sql.Tx, reqNotification *HttpApiNewNotification) (int, error) {
	newUuid := uuid.NewString()

	var item Notification
	if reqNotification != nil {
		item.Type = reqNotification.Type
		item.Body = reqNotification.Body
		item.AttachedFileName = reqNotification.AttachedFileName
		item.DestinationUuid = reqNotification.DestinationUuid
	}

	Query := `INSERT INTO notifications
	        (uuid, type, body, attached_file_name, destination_uuid)
	        VALUES (?, ?, ?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUuid, item.Type, item.Body, item.AttachedFileName, item.DestinationUuid)
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

func GetNotification(db *sql.DB, id int64) (*Notification, error) {
	Req := `SELECT 
				id, 
				uuid, 
				type,
				body,
				attached_file_name,
				destination_uuid,
				is_active, 
				created_at, 
				IF (updated_at IS NULL,"", updated_at)
              FROM notifications
              WHERE id = ? AND is_active = 1
              `

	stmt, err := db.Prepare(Req)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer stmt.Close()

	var item Notification

	err = stmt.QueryRow(id).Scan(&item.Id, &item.Uuid, &item.Type, &item.Body, &item.AttachedFileName, &item.DestinationUuid, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateNotificationTx(tx *sql.Tx, id int64, item *HttpApiUpdatedNotification) error {
	Query := `UPDATE notifications
        SET 
        	type = ?,
        	body = ?,
        	attached_file_name = ?,
        	destination_uuid = ?
        WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Type, item.Body, item.AttachedFileName, item.DestinationUuid, id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneNotificationTx(tx *sql.Tx, impl *Notification) (int, error) {
	req := &HttpApiNewNotification{}
	req.AttachedFileName = impl.AttachedFileName
	req.Body = impl.Body
	req.DestinationUuid = impl.DestinationUuid
	req.Type = impl.Type

	return CreateNotificationTx(tx, req)
}
