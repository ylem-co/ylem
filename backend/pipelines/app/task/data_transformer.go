package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateTransformerTx(tx *sql.Tx, reqTransformer *HttpApiNewTransformer) (int, error) {
	newUuid := uuid.NewString()

	var item Transformer
	if reqTransformer != nil {
		item.Type = reqTransformer.Type
		item.JsonQueryExpression = reqTransformer.JsonQueryExpression
		item.Delimiter = reqTransformer.Delimiter
		item.DecodeFormat = reqTransformer.DecodeFormat
		item.EncodeFormat = reqTransformer.EncodeFormat
		item.CastToType = reqTransformer.CastToType
	}

	Query := `INSERT INTO transformers
	        (uuid, type, json_query_expression, delimiter, decode_format, encode_format, cast_to_type)
	        VALUES (?, ?, ?, ?, ?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newUuid, item.Type, item.JsonQueryExpression, item.Delimiter, item.DecodeFormat, item.EncodeFormat, item.CastToType)
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

func GetTransformer(db *sql.DB, id int64) (*Transformer, error) {
	Req := `SELECT 
				id, 
				uuid, 
				type,
				json_query_expression,
				delimiter,
				decode_format,
				encode_format,
				cast_to_type,
				is_active, 
				created_at, 
				IF (updated_at IS NULL,"", updated_at)
              FROM transformers
              WHERE id = ? AND is_active = 1
              `

	stmt, err := db.Prepare(Req)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var item Transformer

	err = stmt.QueryRow(id).Scan(&item.Id, &item.Uuid, &item.Type, &item.JsonQueryExpression, &item.Delimiter, &item.DecodeFormat, &item.EncodeFormat, &item.CastToType, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &item, nil
}

func UpdateTransformerTx(tx *sql.Tx, id int64, item *HttpApiUpdatedTransformer) error {
	Query := `UPDATE transformers
        SET 
        	type = ?,
        	json_query_expression = ?,
        	delimiter = ?,
        	decode_format = ?,
        	encode_format = ?,
        	cast_to_type = ?
        WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Type, item.JsonQueryExpression, item.Delimiter, item.DecodeFormat, item.EncodeFormat, item.CastToType, id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneTransformerTx(tx *sql.Tx, impl *Transformer) (int, error) {
	req := &HttpApiNewTransformer{}
	req.CastToType = impl.CastToType
	req.DecodeFormat = impl.DecodeFormat
	req.Delimiter = impl.Delimiter
	req.EncodeFormat = impl.EncodeFormat
	req.JsonQueryExpression = impl.JsonQueryExpression
	req.Type = impl.Type

	return CreateTransformerTx(tx, req)
}
