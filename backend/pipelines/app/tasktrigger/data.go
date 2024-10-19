package tasktrigger

import (
	"fmt"
	"database/sql"
	"ylem_pipelines/app/pipeline/run"
	"ylem_pipelines/helpers"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateTaskTrigger(db *sql.DB, taskTrigger *TaskTrigger) (int, string, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, "", err
	}

	ttId, ttUuid, err := CreateTaskTriggerTx(tx, taskTrigger)
	if err != nil {
		_ = tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			return 0, "", err
		}
	}

	return ttId, ttUuid, err
}

func CreateTaskTriggerTx(tx *sql.Tx, taskTrigger *TaskTrigger) (int, string, error) {
	if taskTrigger.Uuid == "" {
		taskTrigger.Uuid = uuid.NewString()
	}

	Query := `INSERT INTO task_triggers
	        (uuid, pipeline_id, trigger_task_id, triggered_task_id, trigger_type, schedule)
	        VALUES (?, ?, ?, ?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		log.Error(err)
		return -1, taskTrigger.Uuid, err
	}
	defer stmt.Close()

	var triggerTaskId interface{} = taskTrigger.TriggerTaskId
	if triggerTaskId.(int64) == 0 {
		triggerTaskId = nil
	}

	result, err := stmt.Exec(taskTrigger.Uuid, taskTrigger.PipelineId, triggerTaskId, taskTrigger.TriggeredTaskId, taskTrigger.TriggerType, taskTrigger.Schedule)
	if err != nil {
		log.Error(err)
		return -1, taskTrigger.Uuid, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, taskTrigger.Uuid, err
	}
	return int(insertID), taskTrigger.Uuid, nil
}

func GetTaskTriggerByUuid(db *sql.DB, uuid string) (*TaskTrigger, error) {
	Query := `SELECT 
				t.id, 
				t.uuid, 
				t.pipeline_id,
				w.uuid as pipeline_uuid,
				IFNULL(t.trigger_task_id, 0), 
				t.triggered_task_id, 
				IFNULL(ts1.uuid, "") as trigger_task_uuid,
				ts2.uuid as triggered_task_uuid,
				t.trigger_type,
				t.schedule, 
				t.created_at, 
				IF (t.updated_at IS NULL,"", t.updated_at)
              FROM task_triggers t
              LEFT JOIN tasks ts1 ON ts1.id = t.trigger_task_id
              JOIN tasks ts2 ON ts2.id = t.triggered_task_id
              JOIN pipelines w ON w.id = t.pipeline_id
              WHERE t.uuid = ? AND t.is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer stmt.Close()

	var taskTrigger TaskTrigger

	err = stmt.QueryRow(uuid).Scan(&taskTrigger.Id, &taskTrigger.Uuid, &taskTrigger.PipelineId, &taskTrigger.PipelineUuid, &taskTrigger.TriggerTaskId, &taskTrigger.TriggeredTaskId, &taskTrigger.TriggerTaskUuid, &taskTrigger.TriggeredTaskUuid, &taskTrigger.TriggerType, &taskTrigger.Schedule, &taskTrigger.CreatedAt, &taskTrigger.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &taskTrigger, nil
}

func GetTaskTriggersByPipelineId(db *sql.DB, id int64) (*TaskTriggers, error) {
	Query := `SELECT 
				t.id, 
				t.uuid, 
				t.pipeline_id,
				w.uuid as pipeline_uuid,
				IFNULL(t.trigger_task_id, 0),
				t.triggered_task_id, 
				IFNULL(ts1.uuid, "") as trigger_task_uuid,
				ts2.uuid as triggered_task_uuid,
				t.trigger_type,
				t.schedule, 
				t.created_at, 
				IF (t.updated_at IS NULL,"", t.updated_at)
              FROM task_triggers t
              LEFT JOIN tasks ts1 ON ts1.id = t.trigger_task_id
              JOIN tasks ts2 ON ts2.id = t.triggered_task_id
              JOIN pipelines w ON w.id = t.pipeline_id
              WHERE t.pipeline_id = ? AND t.is_active = 1
              ORDER BY t.created_at DESC
              `
	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var taskTriggers TaskTriggers

	for rows.Next() {
		var taskTrigger TaskTrigger
		err := rows.Scan(&taskTrigger.Id, &taskTrigger.Uuid, &taskTrigger.PipelineId, &taskTrigger.PipelineUuid, &taskTrigger.TriggerTaskId, &taskTrigger.TriggeredTaskId, &taskTrigger.TriggerTaskUuid, &taskTrigger.TriggeredTaskUuid, &taskTrigger.TriggerType, &taskTrigger.Schedule, &taskTrigger.CreatedAt, &taskTrigger.UpdatedAt)
		if err != nil {
			log.Error(err)
			continue
		}

		taskTriggers.Items = append(taskTriggers.Items, taskTrigger)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return &taskTriggers, err
	}

	return &taskTriggers, err
}

func DeleteTaskTrigger(db *sql.DB, uuid string) error {
	tt, err := GetTaskTriggerByUuid(db, uuid)
	if err != nil {
		return err
	}
	if tt == nil {
		return fmt.Errorf("task trigger %s not found", uuid)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	Query := `UPDATE task_triggers
        SET is_active = 0
        WHERE uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uuid)
	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}

	return nil
}

func UpdateTaskTrigger(db *sql.DB, taskTrigger *TaskTrigger) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	Query := `UPDATE task_triggers
        SET trigger_type = ?, 
        schedule = ?
        WHERE uuid = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskTrigger.TriggerType, taskTrigger.Schedule, taskTrigger.Uuid)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}

	return nil
}

func GetTriggeredTaskIds(db *sql.DB, triggerTaskUuid uuid.UUID, triggerType string, config run.PipelineRunConfig) ([]int64, error) {
	q := `SELECT tt.triggered_task_id
		FROM 
			task_triggers tt
			INNER JOIN tasks trigger_tasks ON trigger_tasks.id = tt.trigger_task_id
			INNER JOIN tasks triggered_tasks ON triggered_tasks.id = tt.triggered_task_id
		WHERE
			trigger_tasks.uuid = ?
			AND tt.is_active = 1
			AND triggered_tasks.is_active = 1
	`

	params := []interface{}{
		triggerTaskUuid,
	}

	if triggerType != "" {
		q = q + " AND tt.trigger_type = ?"
		params = append(params, triggerType)
	}

	triggerClause, triggerParams := helpers.IdListClause("tt.uuid", config.TaskTriggerIds)
	q += ` ` + triggerClause
	params = append(params, triggerParams...)

	taskClause, taskParams := helpers.IdListClause("triggered_tasks.uuid", config.TaskIds)
	q += ` ` + taskClause
	params = append(params, taskParams...)

	result := make([]int64, 0)
	rows, err := db.Query(q, params...)

	if err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return result, err
		}

		result = append(result, id)
	}

	err = rows.Err()
	if err != nil {
		log.Error(err)
		return result, err
	}

	return result, nil
}

// Returns number of inputs of a task
func GetInputCount(db *sql.DB, taskId int64) (int64, error) {
	q := `SELECT COUNT(*) FROM task_triggers WHERE triggered_task_id = ? AND is_active = 1`
	var cnt int64
	err := db.QueryRow(q, taskId).Scan(&cnt)

	return cnt, err
}
