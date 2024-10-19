package task

import (
	"strings"
	"database/sql"
	"ylem_pipelines/app/pipeline/run"
	"ylem_pipelines/helpers"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateTaskWithImplementation(db *sql.DB, task *Task, reqTask HttpApiNewTask) (int, string, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, "", err
	}

	implementationId, err := CreateImplementation(tx, reqTask, task.PipelineUuid)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return -1, task.Uuid, err
	}

	tId, tUuid, err := CreateTaskTx(tx, task, implementationId)
	if err != nil {
		_ = tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			return -1, task.Uuid, err
		}
	}

	return tId, tUuid, err
}

func CreateTaskTx(tx *sql.Tx, task *Task, implementationId int) (int, string, error) {
	task.Uuid = uuid.NewString()

	Query := `INSERT INTO tasks
	        (uuid, pipeline_id, name, type, implementation_id)
	        VALUES (?, ?, ?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		log.Error(err)
		return -1, task.Uuid, err
	}
	defer stmt.Close()

	var implId *int
	if implementationId != -1 {
		implId = &implementationId
	}

	result, err := stmt.Exec(task.Uuid, task.PipelineId, task.Name, task.Type, implId)
	if err != nil {
		log.Error(err)
		return -1, task.Uuid, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, task.Uuid, err
	}

	return int(insertID), task.Uuid, nil
}

func GetTaskByUuid(db *sql.DB, uuid string) (*Task, error) {
	Query := `SELECT 
				t.id, 
				t.uuid, 
				t.name,
				t.severity, 
				t.pipeline_id,
				w.uuid as pipeline_uuid, 
				t.type,
				t.implementation_id,
				t.is_active, 
				t.created_at, 
				IF (t.updated_at IS NULL,"", t.updated_at)
              FROM tasks t
              JOIN pipelines w ON w.id = t.pipeline_id
              WHERE t.uuid = ? AND t.is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer stmt.Close()

	var task Task

	err = stmt.QueryRow(uuid).Scan(&task.Id, &task.Uuid, &task.Name, &task.Severity, &task.PipelineId, &task.PipelineUuid, &task.Type, &task.ImplementationId, &task.IsActive, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	implementation, err := GetImplementation(db, task)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	task.Implementation = implementation

	return &task, nil
}

func GetTaskById(db *sql.DB, id int64) (*Task, error) {
	result, err := GetTasksById(db, id)
	var t *Task
	if len(result) > 0 {
		t = result[0]
	}

	return t, err
}

func GetTasksById(db *sql.DB, ids ...int64) ([]*Task, error) {
	result := make([]*Task, 0)
	if len(ids) == 0 {
		return result, nil
	}
	Query := `SELECT 
				t.id, 
				t.uuid, 
				t.name,
				t.severity, 
				t.pipeline_id,
				w.uuid as pipeline_uuid,
				w.organization_uuid as organization_uuid,
				t.type,
				t.implementation_id,
				t.is_active, 
				t.created_at, 
				IF (t.updated_at IS NULL,"", t.updated_at)
              FROM tasks t
              JOIN pipelines w ON w.id = t.pipeline_id
              WHERE t.id IN (?` + strings.Repeat(",?", len(ids)-1) + `) AND t.is_active = 1
              `

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer stmt.Close()

	args := make([]interface{}, len(ids))
	for k, v := range ids {
		args[k] = v
	}

	rows, err := stmt.Query(args...)
	if err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		log.Error(err)
		return result, err
	}

	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.Id,
			&task.Uuid,
			&task.Name,
			&task.Severity,
			&task.PipelineId,
			&task.PipelineUuid,
			&task.OrganizationUuid,
			&task.Type,
			&task.ImplementationId,
			&task.IsActive,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			log.Error(err)
			return result, err
		}

		implementation, err := GetImplementation(db, task)
		if err != nil {
			log.Error(err)
			return result, err
		}

		task.Implementation = implementation
		result = append(result, &task)
	}

	return result, nil
}

func GetTasksByPipelineId(db *sql.DB, id int64) (*Tasks, error) {
	Query := `SELECT 
				t.id, 
				t.uuid, 
				t.name, 
				t.severity,
				t.pipeline_id,
				w.uuid as pipeline_uuid, 
				t.type,
				t.implementation_id,
				t.is_active, 
				t.created_at, 
				IF (t.updated_at IS NULL,"", t.updated_at)
              FROM tasks t
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

	var tasks Tasks

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Uuid, &task.Name, &task.Severity, &task.PipelineId, &task.PipelineUuid, &task.Type, &task.ImplementationId, &task.IsActive, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			log.Error(err)
			continue
		}

		implementation, err := GetImplementation(db, task)
		if err != nil {
			log.Error(err)
			continue
		}

		task.Implementation = implementation

		tasks.Items = append(tasks.Items, task)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return &tasks, err
	}

	return &tasks, err
}

func GetTasksByOrganizationUuidAndSearchString(db *sql.DB, uuid string, searchString string) (*SearchedTasks, error) {
	Query := `SELECT 
				t.id, 
				t.uuid, 
				t.name,
				t.pipeline_id,
				w.uuid as pipeline_uuid,
				t.type,
				t.is_active, 
				t.created_at, 
				IF (t.updated_at IS NULL,"", t.updated_at),
				IF (f.id IS NULL, 0, f.id),
				IF (f.uuid IS NULL, "", f.uuid)
              FROM tasks t
              JOIN pipelines w ON w.id = t.pipeline_id
              LEFT JOIN folders f ON w.folder_id = f.id
              WHERE 
              	w.organization_uuid = ?
              	AND w.type = "generic"
              	AND w.is_active = 1 
              	AND t.is_active = 1
              	AND t.name LIKE ?
              ORDER BY t.created_at DESC
              LIMIT 5
              `
	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid, "%" + searchString + "%")

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var tasks SearchedTasks

	for rows.Next() {
		var task SearchedTask
		err := rows.Scan(&task.Id, &task.Uuid, &task.Name, &task.PipelineId, &task.PipelineUuid, &task.Type, &task.IsActive, &task.CreatedAt, &task.UpdatedAt, &task.FolderId, &task.FolderUuid)
		if err != nil {
			log.Error(err)
			continue
		}

		tasks.Items = append(tasks.Items, task)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return &tasks, err
	}

	return &tasks, err
}

func DeleteTask(db *sql.DB, id int64, taskType string, implementationId int64) error {
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	defer tx.Rollback() //nolint:all

	// Delete task
	Query := `UPDATE tasks
        SET is_active = 0
        WHERE id = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	// Delete implementation
	table := getDBTable(taskType)

	Query = `UPDATE ` + table + `
        SET is_active = 0
        WHERE id = ?
        `

	stmt, err = db.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(implementationId)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	// Delete task triggers
	Query = `UPDATE task_triggers
        SET is_active = 0
        WHERE trigger_task_id = ?
        `

	stmt, err = db.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	Query = `UPDATE task_triggers
        SET is_active = 0
        WHERE triggered_task_id = ?
        `

	stmt, err = db.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

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

func UpdateTask(db *sql.DB, task *Task, reqTask HttpApiUpdatedTask) error {
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	defer tx.Rollback() //nolint:all

	Query := `UPDATE tasks
        SET name = ?, 
        	severity = ?
        WHERE uuid = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Name, task.Severity, task.Uuid)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	err = UpdateImplementation(tx, task, reqTask)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}

	return nil
}

func GetInitialTasks(db *sql.DB, pipelineUuid uuid.UUID, config run.PipelineRunConfig) ([]*Task, error) {
	params := make([]interface{}, 0)

	triggerClause, triggerParams := helpers.IdListClause("tt.uuid", config.TaskTriggerIds)
	params = append(params, triggerParams...)
	params = append(params, pipelineUuid.String())

	taskClause, taskParams := helpers.IdListClause("t.uuid", config.TaskIds)
	params = append(params, taskParams...)

	q := `
	SELECT
		t.id
	FROM 
		tasks t
		LEFT JOIN task_triggers tt ON tt.triggered_task_id = t.id AND tt.is_active = 1 ` + triggerClause + `
		INNER JOIN pipelines w ON t.pipeline_id = w.id
	WHERE
		w.uuid = ?
		AND tt.id IS NULL
		AND t.is_active = 1
	` + taskClause

	rows, err := db.Query(q, params...)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var tasks []*Task = make([]*Task, 0)
	for rows.Next() {
		var taskId int64
		err = rows.Scan(&taskId)
		if err != nil {
			return nil, err
		}

		t, err := GetTaskById(db, taskId)
		if err != nil {
			return tasks, err
		}

		if t != nil {
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func CloneTaskTx(tx *sql.Tx, orgUuid string, newWfId int64, newWfUuid string, t Task) (*Task, error) {
	implId, err := CloneImplementation(tx, t.Type, t.Implementation, newWfUuid)
	if err != nil {
		return nil, err
	}

	nt := t
	nt.Id = 0
	nt.Uuid = uuid.NewString()
	nt.PipelineId = newWfId
	nt.PipelineUuid = newWfUuid
	nt.ImplementationId = int64(implId)
	nt.IsActive = t.IsActive
	nt.Name = t.Name
	nt.OrganizationUuid = orgUuid
	nt.Severity = t.Severity
	nt.Type = t.Type

	tId, _, err := CreateTaskTx(tx, &nt, implId)
	if err != nil {
		return nil, err
	}
	nt.Id = int64(tId)

	return &nt, nil
}

func GetExternallyTriggeredPipelineCount(db *sql.DB, orgUuid string) (int, error) {
	q := `SELECT 
			COUNT(DISTINCT(w.id))
		FROM 
			tasks t
		JOIN pipelines w ON w.id = t.pipeline_id
		WHERE 
			w.organization_uuid = ? 
			AND t.type = "external_trigger"
			AND t.is_active = 1
			AND w.is_active = 1
			AND w.type = "generic"
		`

	var cnt int
	row := db.QueryRow(
		q,
		orgUuid,
	)

	err := row.Scan(&cnt)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return cnt, nil
}
