package result

import (
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreatePendingTaskRunResult(db *sql.DB, taskId int64, taskUuid uuid.UUID, taskRunUuid uuid.UUID, pipelineRunUuid uuid.UUID) (TaskRunResult, error) {
	trr := TaskRunResult{}
	tx, err := db.Begin()
	if err != nil {
		return trr, err
	}
	trr, err = CreatePendingTaskRunResultTx(tx, taskId, taskUuid, taskRunUuid, pipelineRunUuid)
	if err != nil {
		_ = tx.Rollback()
		return trr, err
	}

	err = tx.Commit()

	return trr, err
}

func PurgeTaskRunResultsTx(tx *sql.Tx, pipelineUuid uuid.UUID) error {
	q := `DELETE
			tre
		FROM
			task_run_errors tre
			INNER JOIN task_run_results trr ON trr.id = tre.task_run_result_id
			INNER JOIN tasks t ON t.id = trr.task_id
			INNER JOIN pipelines w ON w.id = t.pipeline_id
		WHERE w.uuid = ?
		`
	_, err := tx.Exec(q, pipelineUuid)
	if err != nil {
		return err
	}

	q = `DELETE
			trr
		FROM
			task_run_results trr
			INNER JOIN tasks t ON t.id = trr.task_id
			INNER JOIN pipelines w ON w.id = t.pipeline_id
		WHERE w.uuid = ?
		`
	_, err = tx.Exec(q, pipelineUuid)

	return err
}

func CreatePendingTaskRunResultTx(tx *sql.Tx, taskId int64, taskUuid uuid.UUID, taskRunUuid uuid.UUID, pipelineRunUuid uuid.UUID) (TaskRunResult, error) {
	trr := TaskRunResult{
		TaskId:          taskId,
		TaskUuid:        taskUuid,
		TaskRunUuid:     taskRunUuid,
		PipelineRunUuid: pipelineRunUuid,
		State:           StatePending,
	}

	q := `
	INSERT INTO 
		task_run_results(
			state,
			task_id, 
			task_run_uuid,
			pipeline_run_uuid,
			is_successful,
			output
		) 
	VALUES 
		(
			?,
			?,
			?,
			?,
			?,
			?
		)
	ON DUPLICATE KEY UPDATE
		state = ?,
		task_run_uuid = ?,
		pipeline_run_uuid = ? 
	`

	r, err := tx.Exec(
		q,
		StatePending,
		taskId,
		taskRunUuid,
		pipelineRunUuid,
		false,
		"",
		StatePending,
		taskRunUuid,
		pipelineRunUuid,
	)

	if err != nil {
		return trr, err
	}

	trr.Id, err = r.LastInsertId()

	return trr, err
}

func FindTaskRunResults(db *sql.DB, pipelineId int64) ([]TaskRunResult, error) {
	results := []TaskRunResult{}
	q := `
	SELECT
		trr.id,
		trr.state, 
		trr.task_id,
		t.uuid,
		trr.task_run_uuid,
		trr.pipeline_run_uuid,
		trr.is_successful,
		trr.output
	FROM 
		task_run_results trr
		INNER JOIN tasks t ON t.id = trr.task_id AND t.is_active = 1
	WHERE
		t.pipeline_id = ?
	ORDER BY 
		trr.id ASC
	`
	rows, err := db.Query(q, pipelineId)
	if err == sql.ErrNoRows {
		return results, nil
	}

	if err != nil {
		return results, err
	}

	for rows.Next() {
		trr := TaskRunResult{}
		err = rows.Scan(
			&trr.Id,
			&trr.State,
			&trr.TaskId,
			&trr.TaskUuid,
			&trr.TaskRunUuid,
			&trr.PipelineRunUuid,
			&trr.IsSuccessful,
			&trr.Output,
		)

		if err != nil {
			return results, err
		}

		trr.Errors, err = FindTaskRunErrors(db, trr.Id)
		if err != nil {
			return results, err
		}
		results = append(results, trr)
	}

	return results, nil
}

func FindTaskRunErrors(db *sql.DB, taskRunResultId int64) ([]TaskRunError, error) {
	result := make([]TaskRunError, 0)
	q := `
	SELECT
		id,
		task_run_result_id,
		code,
		severity,
		message
	FROM
		task_run_errors
	WHERE
		task_run_result_id = ?
	`
	rows, err := db.Query(q, taskRunResultId)

	if err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		return result, err
	}

	for rows.Next() {
		tre := TaskRunError{}
		err = rows.Scan(
			&tre.Id,
			&tre.TaskRunResultId,
			&tre.Code,
			&tre.Severity,
			&tre.Message,
		)
		if err != nil {
			return result, err
		}

		result = append(result, tre)
	}

	return result, nil
}

func updateTaskRunErrorsTx(tx *sql.Tx, trr *TaskRunResult) error {
	q := `DELETE FROM task_run_errors WHERE task_run_result_id = ?`
	_, err := tx.Exec(q, trr.Id)
	if err != nil {
		return err
	}

	q = `
	INSERT INTO
		task_run_errors(
			task_run_result_id,
			code,
			severity,
			message
		)
	VALUES
		(
			?,
			?,
			?,
			?
		)
	`
	for _, e := range trr.Errors {
		r, err := tx.Exec(
			q,
			trr.Id,
			e.Code,
			e.Severity,
			e.Message,
		)

		if err != nil {
			return err
		}
		e.Id, err = r.LastInsertId()
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateTaskRunResult(db *sql.DB, taskUuid uuid.UUID, trr *TaskRunResult) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	result, err := UpdateTaskRunResultTx(tx, taskUuid, trr)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return result, err
}

func UpdateTaskRunResultTx(tx *sql.Tx, taskUuid uuid.UUID, trr *TaskRunResult) (bool, error) {
	q := `
	UPDATE 
		task_run_results trr
		INNER JOIN tasks t ON t.id = trr.task_id
	SET
		trr.state = ?,
		trr.is_successful = ?,
		trr.output = ?
	WHERE
		t.uuid = ? 
		AND trr.pipeline_run_uuid = ?
	`

	r, err := tx.Exec(
		q, StateExecuted,
		trr.IsSuccessful,
		trr.Output,
		taskUuid.String(),
		trr.PipelineRunUuid,
	)

	if err != nil {
		return false, err
	}

	numRows, err := r.RowsAffected()
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		log.Tracef("No task run results stored for run %s", trr.PipelineRunUuid)
		return false, nil
	}

	log.Tracef("Task run result %s stored", trr.PipelineRunUuid)

	q = `SELECT 
			trr.id 
		FROM 
			task_run_results trr 
			INNER JOIN tasks t ON t.id = trr.task_id 
		WHERE
			t.uuid = ? 
			AND trr.pipeline_run_uuid = ?`

	res, err := tx.Query(
		q,
		taskUuid.String(),
		trr.PipelineRunUuid.String(),
	)

	if err == sql.ErrNoRows {
		return true, nil
	}

	if err != nil {
		log.Error(err)
		return false, nil
	}

	res.Next()
	err = res.Scan(&trr.Id)
	if err != nil {
		log.Error(err)
		return true, nil
	}

	err = res.Close()
	if err != nil {
		log.Error(err)
		return true, nil
	}

	err = updateTaskRunErrorsTx(tx, trr)
	if err != nil {
		return false, err
	}

	return true, nil
}
