package task

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateRunPipelineTx(tx *sql.Tx, rw *HttpApiNewRunPipeline) (int, error) {
	newUuid := uuid.NewString()

	var item RunPipeline;
	if rw != nil {
		item.PipelineUuid = rw.PipelineUuid
	}

	q := `INSERT INTO run_pipelines
	        (uuid, pipeline_uuid)
	        VALUES (?, ?)
	        `

	result, err := tx.Exec(q, newUuid, item.PipelineUuid)
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

func GetRunPipeline(db *sql.DB, id int64) (*RunPipeline, error) {
	q := `SELECT 
				id,
				uuid,
				pipeline_uuid,
				is_active,
				created_at,
				IF (updated_at IS NULL, "", updated_at)
              FROM run_pipelines
              WHERE id = ? AND is_active = 1
              `

	var item RunPipeline
	err := db.QueryRow(q, id).Scan(&item.Id, &item.Uuid, &item.PipelineUuid, &item.IsActive, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &item, nil
}

func UpdateRunPipelineTx(tx *sql.Tx, id int64, item *HttpApiUpdatedRunPipeline) error {
	q := `UPDATE run_pipelines
        SET pipeline_uuid = ?
        WHERE id = ?
        `

	_, err := tx.Exec(q, item.PipelineUuid, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	return nil
}

func CloneRunPipelineTx(tx *sql.Tx, impl *RunPipeline) (int, error) {
	req := &HttpApiNewRunPipeline{}
	req.PipelineUuid = impl.PipelineUuid

	return CreateRunPipelineTx(tx, req)
}
