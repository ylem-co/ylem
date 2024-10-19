package pipeline

import (
	"strings"
	"database/sql"
	"ylem_pipelines/app/schedule"
	"ylem_pipelines/app/pipeline/common"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

func CreatePipeline(db *sql.DB, pipeline *Pipeline) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	wId, err := CreatePipelineTx(tx, pipeline)
	if err != nil {
		_ = tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			return 0, err
		}
	}

	return wId, err
}

func CreatePipelineTx(tx *sql.Tx, pipeline *Pipeline) (int, error) {
	pipeline.Uuid = uuid.NewString()
	pipeline.IsActive = PipelineIsActive

	Query := `INSERT INTO pipelines
	        (uuid, organization_uuid, name, type, creator_uuid, elements_layout, folder_id, preview, is_active, is_template, schedule)
	        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	defer stmt.Close()

	var result sql.Result
	if pipeline.FolderId != 0 {
		result, err = stmt.Exec(pipeline.Uuid, pipeline.OrganizationUuid, pipeline.Name, pipeline.Type, pipeline.CreatorUuid, pipeline.ElementsLayout, pipeline.FolderId, pipeline.Preview, pipeline.IsActive, pipeline.IsTemplate, pipeline.Schedule)
	} else {
		result, err = stmt.Exec(pipeline.Uuid, pipeline.OrganizationUuid, pipeline.Name, pipeline.Type, pipeline.CreatorUuid, pipeline.ElementsLayout, nil, pipeline.Preview, pipeline.IsActive, pipeline.IsTemplate, pipeline.Schedule)
	}

	if err != nil {
		return -1, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	return int(insertID), nil
}

func GetPipelineByUuid(db *sql.DB, uuid string) (*Pipeline, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer tx.Commit() //nolint:all

	return GetPipelineByUuidTx(tx, uuid)
}

func GetPipelineByUuidTx(tx *sql.Tx, uuid string) (*Pipeline, error) {
	Query := `SELECT 
				w.id, 
				w.uuid, 
				w.name, 
				w.type,
				w.organization_uuid, 
				w.creator_uuid, 
				w.elements_layout, 
				w.preview,
				w.is_paused, 
				w.is_active, 
				w.created_at, 
				IF (w.updated_at IS NULL,"", w.updated_at),
				w.is_template,
				w.schedule,
				f.id,
				f.uuid
              FROM pipelines w
              LEFT JOIN folders f ON w.folder_id = f.id
              WHERE w.uuid = ? AND w.is_active = 1
              `

	stmt, err := tx.Prepare(Query)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	var pipeline Pipeline
	var folderId sql.NullInt64
	var folderUuid sql.NullString
	err = stmt.QueryRow(uuid).Scan(&pipeline.Id, &pipeline.Uuid, &pipeline.Name, &pipeline.Type, &pipeline.OrganizationUuid, &pipeline.CreatorUuid, &pipeline.ElementsLayout, &pipeline.Preview, &pipeline.IsPaused, &pipeline.IsActive, &pipeline.CreatedAt, &pipeline.UpdatedAt, &pipeline.IsTemplate, &pipeline.Schedule, &folderId, &folderUuid)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if folderUuid.Valid {
		pipeline.FolderUuid = folderUuid.String
	}

	if folderId.Valid {
		pipeline.FolderId = folderId.Int64
	}

	return &pipeline, nil
}

func GetPipelineById(db *sql.DB, id int64) (*Pipeline, error) {
	Query := `SELECT 
				w.id, 
				w.uuid, 
				w.name, 
				w.type,
				w.organization_uuid, 
				w.creator_uuid, 
				w.elements_layout, 
				w.preview,
				w.is_paused, 
				w.is_active, 
				w.created_at, 
				IF (w.updated_at IS NULL,"", w.updated_at),
				w.is_template,
				w.schedule,
				f.id,
				f.uuid
              FROM pipelines w
              LEFT JOIN folders f ON w.folder_id = f.id
              WHERE w.id = ? AND w.is_active = 1
              `

	row := db.QueryRow(Query, id)
	if row.Err() == sql.ErrNoRows {
		return nil, nil
	}

	if row.Err() != nil {
		log.Error(row.Err())
		return nil, row.Err()
	}

	var pipeline Pipeline
	var folderId sql.NullInt64
	var folderUuid sql.NullString
	err := row.Scan(&pipeline.Id, &pipeline.Uuid, &pipeline.Name, &pipeline.Type, &pipeline.OrganizationUuid, &pipeline.CreatorUuid, &pipeline.ElementsLayout, &pipeline.Preview, &pipeline.IsPaused, &pipeline.IsActive, &pipeline.CreatedAt, &pipeline.UpdatedAt, &pipeline.IsTemplate, &pipeline.Schedule, &folderId, &folderUuid)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if folderUuid.Valid {
		pipeline.FolderUuid = folderUuid.String
	}

	if folderId.Valid {
		pipeline.FolderId = folderId.Int64
	}

	return &pipeline, nil
}

func GetPipelinesByOrganizationUuid(db *sql.DB, uuid string) (*Pipelines, error) {
	Query := `SELECT 
				w.id, 
				w.uuid, 
				w.name, 
				w.type,
				w.is_paused, 
				w.organization_uuid, 
				w.creator_uuid, 
				w.elements_layout, 
				w.created_at, 
				IF (w.updated_at IS NULL,"", w.updated_at),
				w.is_template,
				w.schedule,
				f.id,
				f.uuid
              FROM pipelines w
              LEFT JOIN folders f ON w.folder_id = f.id
              WHERE 
              	w.organization_uuid = ? 
              	AND w.is_active = 1
              ORDER BY w.updated_at DESC`
	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var pipelines Pipelines

	for rows.Next() {
		var pipeline Pipeline
		var folderId sql.NullInt64
		var folderUuid sql.NullString
		err := rows.Scan(&pipeline.Id, &pipeline.Uuid, &pipeline.Name, &pipeline.Type, &pipeline.IsPaused, &pipeline.OrganizationUuid, &pipeline.CreatorUuid, &pipeline.ElementsLayout, &pipeline.CreatedAt, &pipeline.UpdatedAt, &pipeline.IsTemplate, &pipeline.Schedule, &folderId, &folderUuid)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if folderUuid.Valid {
			pipeline.FolderUuid = folderUuid.String
		}

		if folderId.Valid {
			pipeline.FolderId = folderId.Int64
		}

		pipelines.Items = append(pipelines.Items, pipeline)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return &pipelines, err
	}

	return &pipelines, err
}

func GetPipelinesByOrganizationUuidAndSearchString(db *sql.DB, uuid string, searchString string) (*SearchedPipelines, error) {
	Query := `SELECT 
				w.id, 
				w.uuid, 
				w.name, 
				w.type,
				w.organization_uuid, 
				w.creator_uuid,
				w.created_at, 
				IF (w.updated_at IS NULL,"", w.updated_at),
				w.is_template,
				w.schedule,
				f.id,
				f.uuid
              FROM pipelines w
              LEFT JOIN folders f ON w.folder_id = f.id
              WHERE 
              	w.organization_uuid = ? 
              	AND w.is_active = 1
              	AND w.name LIKE ?
              ORDER BY w.updated_at DESC
              LIMIT 5`
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

	var pipelines SearchedPipelines

	for rows.Next() {
		var pipeline SearchedPipeline
		var folderId sql.NullInt64
		var folderUuid sql.NullString
		err := rows.Scan(&pipeline.Id, &pipeline.Uuid, &pipeline.Name, &pipeline.Type, &pipeline.OrganizationUuid, &pipeline.CreatorUuid, &pipeline.CreatedAt, &pipeline.UpdatedAt, &pipeline.IsTemplate, &pipeline.Schedule, &folderId, &folderUuid)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if folderUuid.Valid {
			pipeline.FolderUuid = folderUuid.String
		}

		if folderId.Valid {
			pipeline.FolderId = folderId.Int64
		}

		pipelines.Items = append(pipelines.Items, pipeline)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return &pipelines, err
	}

	return &pipelines, err
}

func GetPipelinesByFolderIdAndOrganizationUuid(db *sql.DB, uuid string, id int64) (*Pipelines, error) {
	var Query string

	if id != 0 {
		Query = `SELECT 
				w.id, 
				w.uuid, 
				w.name, 
				w.type,
				w.is_paused, 
				w.organization_uuid, 
				w.creator_uuid, 
				w.elements_layout, 
				w.created_at,
				IF (w.updated_at IS NULL,"", w.updated_at),
				w.is_template,
				w.schedule,
              	f.id,
				f.uuid
              FROM pipelines w
              LEFT JOIN folders f ON w.folder_id = f.id
              WHERE 
              	w.folder_id = ? 
              	AND w.organization_uuid = ? 
              	AND w.is_active = 1
              ORDER BY w.updated_at DESC`
	} else {
		Query = `SELECT 
				w.id, 
				w.uuid, 
				w.name, 
				w.type,
				w.is_paused, 
				w.organization_uuid, 
				w.creator_uuid, 
				w.elements_layout, 
				w.created_at, 
				IF (w.updated_at IS NULL,"", w.updated_at),
				w.is_template,
				w.schedule,
              	f.id,
				f.uuid
              FROM pipelines w
              LEFT JOIN folders f ON w.folder_id = f.id
              WHERE 
              	w.folder_id IS NULL 
              	AND w.organization_uuid = ?
              	AND w.is_active = 1
              ORDER BY w.updated_at DESC`
	}
	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	if id != 0 {
		rows, err = stmt.Query(id, uuid)
	} else {
		rows, err = stmt.Query(uuid)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var pipelines Pipelines

	for rows.Next() {
		var pipeline Pipeline
		var folderId sql.NullInt64
		var folderUuid sql.NullString
		err := rows.Scan(&pipeline.Id, &pipeline.Uuid, &pipeline.Name, &pipeline.Type, &pipeline.IsPaused, &pipeline.OrganizationUuid, &pipeline.CreatorUuid, &pipeline.ElementsLayout, &pipeline.CreatedAt, &pipeline.UpdatedAt, &pipeline.IsTemplate, &pipeline.Schedule, &folderId, &folderUuid)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if folderUuid.Valid {
			pipeline.FolderUuid = folderUuid.String
		}

		if folderId.Valid {
			pipeline.FolderId = folderId.Int64
		}

		pipelines.Items = append(pipelines.Items, pipeline)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return &pipelines, err
	}

	return &pipelines, err
}

func DeletePipeline(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback() //nolint:all

	// Delete pipeline
	Query := `UPDATE pipelines
        SET is_active = 0
        WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
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

	// Delete tasks
	Query = `UPDATE tasks
        SET is_active = 0
        WHERE pipeline_id = ?
        `

	stmt, err = tx.Prepare(Query)
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

	// Delete task triggers
	Query = `UPDATE task_triggers
        SET is_active = 0
        WHERE pipeline_id = ?
        `

	stmt, err = tx.Prepare(Query)
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

	_, err = schedule.DeleteForPipeline(tx, id)
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

func TogglePipeline(db *sql.DB, pipeline *Pipeline) error {
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback() //nolint:all

	Query := `UPDATE pipelines
        SET is_paused = !is_paused
        WHERE id = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pipeline.Id)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	if pipeline.IsPaused == 0 {
		_, err = schedule.DeleteForPipeline(tx, pipeline.Id)
		if err != nil && err != sql.ErrNoRows {
			_ = tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}

	return nil
}

func UpdatePipeline(db *sql.DB, pipeline *Pipeline, isScheduleChanged bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = UpdatePipelineTx(tx, pipeline, isScheduleChanged)
	if err != nil {
		_ = tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return err
}

func UpdatePipelineTx(tx *sql.Tx, pipeline *Pipeline, isScheduleChanged bool) error {
	Query := `UPDATE pipelines
        SET name = ?,
        	folder_id = ?,
        	elements_layout = ?,
			schedule = ?
        WHERE uuid = ?
        `

	stmt, err := tx.Prepare(Query)
	if err != nil {
		log.Error(err)
		return err
	}
	defer stmt.Close()

	if pipeline.FolderId != 0 {
		_, err = stmt.Exec(pipeline.Name, pipeline.FolderId, pipeline.ElementsLayout, pipeline.Schedule, pipeline.Uuid)
	} else {
		_, err = stmt.Exec(pipeline.Name, nil, pipeline.ElementsLayout, pipeline.Schedule, pipeline.Uuid)
	}

	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}

	if isScheduleChanged {
		_, err = schedule.DeleteForPipeline(tx, pipeline.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Error(err)
			return err
		}
	}

	return nil
}

func UpdatePipelinePreview(db *sql.DB, pipeline *Pipeline) error {
	Query := `UPDATE pipelines
        SET preview = ?
        WHERE uuid = ?
        `

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pipeline.Preview, pipeline.Uuid)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}

	return nil
}

func UpdateElementsLayoutTx(tx *sql.Tx, pipeline *Pipeline) error {
	Query := `UPDATE pipelines
        SET elements_layout = ?
        WHERE uuid = ?
        `

	_, err := tx.Exec(Query, pipeline.ElementsLayout, pipeline.Uuid)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}

	return nil
}

func FindAllAvailableTemplates(db *sql.DB, orgUuid, creatorUuid string, onlyShared bool) (*Pipelines, error) {
	params := []interface{}{}
	clauses := []string{}

	if orgUuid != "" {
		clauses = append(clauses, "w.organization_uuid = ?")
		params = append(params, orgUuid)
	}

	if creatorUuid != "" {
		clauses = append(clauses, "w.creator_uuid = ?")
		params = append(params, creatorUuid)
	}

	if onlyShared {
		clauses = append(clauses, "sw.id IS NOT NULL")
	}

	clauseStr := ""
	if len(clauses) > 0 {
		clauseStr = " AND " + strings.Join(clauses, " AND ")
	}

	q := `SELECT 
			w.id,
			w.uuid,
			w.name,
			w.type,
			w.organization_uuid,
			w.creator_uuid,
			w.elements_layout,
			w.preview,
			w.is_paused,
			w.is_active,
			w.created_at,
			IF (w.updated_at IS NULL,"", w.updated_at),
			w.is_template,
			w.schedule
		FROM pipelines w
		LEFT JOIN shared_pipelines sw ON w.uuid = sw.pipeline_uuid AND sw.is_active = 1
		WHERE
			w.is_active = 1
			AND w.is_template = 1` + clauseStr

	result := &Pipelines{
		Items: []Pipeline{},
	}

	rows, err := db.Query(q, params...)
	if err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var pipeline Pipeline
		err := rows.Scan(
			&pipeline.Id,
			&pipeline.Uuid,
			&pipeline.Name,
			&pipeline.Type,
			&pipeline.OrganizationUuid,
			&pipeline.CreatorUuid,
			&pipeline.ElementsLayout,
			&pipeline.Preview,
			&pipeline.IsPaused,
			&pipeline.IsActive,
			&pipeline.CreatedAt,
			&pipeline.UpdatedAt,
			&pipeline.IsTemplate,
			&pipeline.Schedule,
		)
		if err != nil {
			return result, err
		}

		result.Items = append(result.Items, pipeline)
	}

	return result, nil
}

func GetCurrentPipelineCount(db *sql.DB, orgUuid, wfType string) (int64, error) {
	q := `SELECT 
			COUNT(*) 
		FROM 
			pipelines 
		WHERE 
			organization_uuid = ? 
			AND type = ?
			AND is_active = 1
			AND is_template = 0
		`

	var cnt int64
	row := db.QueryRow(
		q,
		orgUuid,
		wfType,
	)

	err := row.Scan(&cnt)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return cnt, nil
}

func GetRunsByOrganizationUuid(db *sql.DB, orgUuid, wfType string) (*PipelineRunsPerMonths, error) {
	var Query string

	if wfType == common.PipelineTypeGeneric {
		Query = "SELECT `year_month`, run_count FROM pipeline_run_counts_monthly WHERE organization_uuid = ? ORDER BY id DESC LIMIT 12"
	} else {
		Query = "SELECT `year_month`, run_count FROM metrics_run_counts_monthly WHERE organization_uuid = ? ORDER BY id DESC LIMIT 12"
	}
	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	rows, err = stmt.Query(orgUuid)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var runs PipelineRunsPerMonths

	for rows.Next() {
		var run PipelineRunsPerMonth
		err := rows.Scan(&run.YearMonth, &run.RunCount)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		runs.Items = append(runs.Items, run)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return &runs, err
	}

	return &runs, err
}

