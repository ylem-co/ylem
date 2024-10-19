package schedule

import (
	"bytes"
	"strings"
	"time"
	"errors"
	"database/sql"
	"encoding/json"
	"ylem_pipelines/app/pipeline/common"
	"ylem_pipelines/helpers"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func GetLastScheduledRunsForScheduling(db *sql.DB, period time.Duration) ([]ScheduledRun, *sql.Tx, error) {
	q := `SELECT
			w.id,
			MIN(w.schedule),
			MAX(sr.execute_at) AS max_execute_at
		FROM
			pipelines w
			LEFT JOIN scheduled_runs sr ON w.id = sr.pipeline_id
		WHERE
			w.is_active = 1
			AND w.is_paused = 0
			AND w.schedule <> ""
		GROUP BY
			w.id
		HAVING
			max_execute_at IS NULL
			OR max_execute_at < ?
		ORDER BY
			max_execute_at ASC
		LIMIT 100
	`

	t := time.Now().Add(period + (period / 10))

	rows, err := db.Query(q, t)
	if err != nil {
		return nil, nil, err
	}

	result := make([]ScheduledRun, 0)

	for rows.Next() {
		var sr ScheduledRun
		err := rows.Scan(&sr.PipelineId, &sr.Schedule, &sr.ExecuteAt)
		if err != nil {
			log.Error(err)
			continue
		}

		result = append(result, sr)
	}

	log.Debugf("Found %d pipelines for schedule generation", len(result))

	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}
	return result, tx, nil
}

func AddScheduledRunsTx(tx *sql.Tx, srs []ScheduledRun) error {
	batches := helpers.ChunkSlice(len(srs), 5000)

	for _, batch := range batches {
		var srsBatch []ScheduledRun = make([]ScheduledRun, len(batch))
		for i, idx := range batch {
			srsBatch[i] = srs[idx]
		}

		err := doAddScheduledRuns(tx, srsBatch)
		if err != nil {
			return err
		}
	}

	return nil
}

func doAddScheduledRuns(tx *sql.Tx, srs []ScheduledRun) error {
	if len(srs) == 0 {
		return nil
	}

	var q bytes.Buffer
	q.WriteString("INSERT INTO scheduled_runs(pipeline_run_uuid, pipeline_id, input, env_vars, config, execute_at) VALUES ")
	params := make([]interface{}, 0)
	for k, sr := range srs {
		q.WriteString("(?, ?, ?, ?, ?, ?)")
		if k < len(srs)-1 {
			q.WriteString(", ")
		}

		if sr.PipelineRunUuid != uuid.Nil {
			params = append(params, sr.PipelineRunUuid.String())
		} else {
			params = append(params, nil)
		}
		params = append(params, sr.PipelineId)
		params = append(params, sr.Input)

		envVars, err := json.Marshal(sr.EnvVars)
		if err != nil {
			return err
		}
		params = append(params, envVars)

		config, err := json.Marshal(sr.Config)
		if err != nil {
			return err
		}
		params = append(params, config)
		params = append(params, sr.ExecuteAt)
	}

	stmt, err := tx.Prepare(q.String())
	if err != nil {
		return err
	}

	_, err = stmt.Exec(params...)

	if err == nil {
		log.Debugf("Generated %d schedule items", len(srs))
	}

	return err
}

func GetSchedulesForPublishing(db *sql.DB, period time.Duration) ([]ScheduledRun, *sql.Tx, error) {
	result := make([]ScheduledRun, 0)

	// loading IDs
	ids := make([]interface{}, 0)
	q := `SELECT
			id
		FROM
			scheduled_runs
		WHERE
			execute_at < ?
		ORDER BY
			execute_at ASC
		LIMIT 10000
		`

	until := time.Now().Add(period)
	rows, err := db.Query(q, until)
	if err != nil {
		log.Error(err)
		return result, nil, err
	}
	defer rows.Close()

	var id int64
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return result, nil, nil
	}

	// loading scheduled runs in a transaction by ID to avoid gap locks
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	q = `SELECT
			id, 
			pipeline_run_uuid,
			pipeline_id, 
			input,
			env_vars,
			config,
			execute_at
		FROM
			scheduled_runs
		WHERE
		id IN (?` + strings.Repeat(",?", len(ids)-1) + `)
		FOR UPDATE
		SKIP LOCKED
	`

	rows, err = tx.Query(q, ids...)
	if err != nil {
		return nil, tx, err
	}

	defer rows.Close()

	for rows.Next() {
		var sr ScheduledRun
		var envVars, config []byte
		err := rows.Scan(&sr.Id, &sr.PipelineRunUuid, &sr.PipelineId, &sr.Input, &envVars, &config, &sr.ExecuteAt)
		if err != nil {
			log.Errorf("Error reading rows: %s", err)
			return result, tx, err
		}

		if len(envVars) > 0 {
			err = json.Unmarshal(envVars, &sr.EnvVars)
			if err != nil {
				log.Errorf("Env vars unmarshalling error: %s", err)
				continue
			}
		}

		if len(config) > 0 {
			err = json.Unmarshal(config, &sr.Config)
			if err != nil {
				log.Errorf("Pipeline run config unmarshalling error: %s", err)
				continue
			}
		}
		result = append(result, sr)
	}

	return result, tx, nil
}

func Delete(db *sql.DB, id int64) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	result, err := DeleteTx(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return result, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return result, nil
}

func DeleteTx(tx *sql.Tx, id int64) (int64, error) {
	q := "DELETE FROM scheduled_runs WHERE id = ?"
	stmt, err := tx.Prepare(q)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		log.Debugf("Unable to delete scheduled run %d: %s", id, err.Error())
		return 0, err
	}

	log.Debugf("Deleted scheduled run %d", id)

	return result.RowsAffected()
}

func DeleteForPipeline(tx *sql.Tx, pipelineId int64) (int64, error) {
	q := `
		DELETE FROM scheduled_runs 
		WHERE pipeline_id = ?`

	stmt, err := tx.Prepare(q)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(pipelineId)
	if err != nil {
		log.Debugf("Unable to delete scheduled runs for pipeline %d: %s", pipelineId, err.Error())
		return 0, err
	}

	log.Debugf("Deleted scheduled runs for pipeline %d", pipelineId)

	return result.RowsAffected()
}

func GetCurrentPipelineRunCount(db *sql.DB, orgUuid uuid.UUID, wfType string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	result, err := GetCurrentPipelineRunCountTx(tx, orgUuid, wfType)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	return result, err
}

func GetCurrentPipelineRunCountTx(tx *sql.Tx, orgUuid uuid.UUID, wfType string) (int64, error) {
	tableName, err := getRunCountTableName(wfType)
	if err != nil {
		return -1, err
	}

	q := `SELECT 
			run_count 
		FROM 
			 ` + tableName + `
		WHERE 
			organization_uuid = ? 
			AND ` + "`year_month`" + ` = ?`

	var cnt int64
	row := tx.QueryRow(
		q,
		orgUuid.String(),
		currentYearMonth(),
	)

	err = row.Scan(&cnt)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return cnt, nil
}

func IncrementCurrentPipelineRunCount(tx *sql.Tx, orgUuid uuid.UUID, wfType string) error {
	tableName, err := getRunCountTableName(wfType)
	if err != nil {
		return err
	}

	q := `INSERT INTO 
				` + tableName + `(
					organization_uuid,
					` + "`year_month`" + `,
					run_count
				) 
			VALUES(?, ?, ?) 
			ON DUPLICATE KEY
				UPDATE run_count = run_count + 1`

	_, err = tx.Exec(q, orgUuid.String(), currentYearMonth(), 1)

	return err
}

func getRunCountTableName(wfType string) (string, error) {
	switch wfType {
	case common.PipelineTypeGeneric:
		return "pipeline_run_counts_monthly", nil
	case common.PipelineTypeMetric:
		return "metrics_run_counts_monthly", nil
	}

	return "", errors.New("unknown pipeline type")
}

func currentYearMonth() string {
	return time.Now().Format("2006-01")
}
