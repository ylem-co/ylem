package readmodel

import (
	"fmt"
	"time"
	"strconv"
	"database/sql"
	"ylem_statistics/config"
	"ylem_statistics/domain/readmodel/dto"
	"ylem_statistics/services/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PipelineReadModel struct {
	db         *gorm.DB
	statsTable string
}

func (rm *PipelineReadModel) GetPipelineStats(uuid uuid.UUID, start time.Time, end time.Time) (dto.Stats, error) {
	var result dto.Stats

	query := `SELECT
				any(t.OrganizationUuid) AS OrganizationUuid,
				sum(t.SuccessCount) AS SuccessCount,
				sum(t.FailureCount) AS FailureCount,
				toInt32(avg(t.Duration)) AS AverageDuration
			FROM
			(
				SELECT 
						any(organization_uuid) AS OrganizationUuid,
						count(distinct(pipeline_run_uuid)) - sum(is_fatal_failure) AS SuccessCount,
						sum(is_fatal_failure) AS FailureCount,
						sum(duration) as Duration
				FROM ` + rm.statsTable + `
				WHERE
						pipeline_uuid = @pipelineUuid
						AND executed_at BETWEEN @start AND @end
				GROUP BY 
					pipeline_run_uuid
			) t
		`

	rm.db.Raw(
		query,
		sql.Named("pipelineUuid", uuid),
		sql.Named("start", start),
		sql.Named("end", end),
	).Scan(&result)

	if rm.db.Error != nil {
		return result, rm.db.Error
	}

	lastRun, err := rm.GetLastPipelineRun(uuid, end)
	if err != nil {
		return result, err
	}

	result.IsLastRunSuccessful = lastRun.IsSuccessful
	result.LastRunExecutedAt = lastRun.ExecutedAt
	result.LastRunDuration = lastRun.Duration

	return result, nil
}

func (rm *PipelineReadModel) GetPipelineAggregatedStats(pipelineUuid uuid.UUID, start time.Time, period Period, periodCount uint) ([]dto.AggregatedStats, error) {
	var result []dto.AggregatedStats

	query := `SELECT
				any(t.OrganizationUuid) AS OrganizationUuid,
				dateTrunc(@period, RunExecutedAt) AS Start,
				(dateTrunc(@period, RunExecutedAt) + INTERVAL 1 ` + string(period) + ` - INTERVAL 1 SECOND) AS End,
				sum(t.SuccessCount) AS SuccessCount,
				sum(t.FailureCount) AS FailureCount,
				toInt32(avg(t.Duration)) AS AverageDuration
			FROM
			(
				SELECT 
						any(organization_uuid) AS OrganizationUuid,
						max(executed_at) AS RunExecutedAt,
						count(distinct(pipeline_run_uuid)) - sum(is_fatal_failure) AS SuccessCount,
						sum(is_fatal_failure) AS FailureCount,
						sum(duration) as Duration
				FROM ` + rm.statsTable + `
				WHERE
						pipeline_uuid = @pipelineUuid
						AND executed_at BETWEEN toDate(@start) AND (toDate(@start) + INTERVAL ` + strconv.FormatUint(uint64(periodCount), 10) + ` ` + string(period) + ` - INTERVAL 1 SECOND)
				GROUP BY 
					pipeline_run_uuid
			) t
			GROUP BY
				dateTrunc(@period, RunExecutedAt)
		`

	rm.db.Raw(
		query,
		sql.Named("pipelineUuid", pipelineUuid),
		sql.Named("start", start),
		sql.Named("period", period),
	).Scan(&result)

	return result, rm.db.Error
}

func (rm *PipelineReadModel) GetLastPipelineRun(pipelineUuid uuid.UUID, end time.Time) (dto.RunStats, error) {
	var result dto.RunStats

	query := `SELECT
				organization_uuid,
				is_successful AS IsSuccessful,
				is_final_task,
				executed_at AS ExecutedAt,
				SUM(duration) OVER(PARTITION BY pipeline_run_uuid) AS Duration
			FROM ` + rm.statsTable + `
			WHERE
				pipeline_uuid = @pipelineUuid
				AND executed_at <= @end
			ORDER BY 
				executed_at DESC,
				is_final_task DESC
			LIMIT 1
			`

	rm.db.Raw(query,
		sql.Named("pipelineUuid", pipelineUuid),
		sql.Named("end", end),
	).Scan(&result)

	if rm.db.Error != nil {
		return result, rm.db.Error
	}

	return result, nil
}

func (rm *PipelineReadModel) GetPipelineRunResultAvg(pipelineUuid uuid.UUID, period Period, periodCount uint64) (float64, error) {
	var result float64

	query := `SELECT
				ifNotFinite(avg(metric_value), 0)
			FROM ` + rm.statsTable + `
			WHERE
				pipeline_uuid = @pipelineUuid
				AND is_metric_value_set = 1
				AND executed_at BETWEEN (NOW() - INTERVAL @interval) AND NOW()
			`

	tx := rm.db.Raw(
		query,
		sql.Named("pipelineUuid", pipelineUuid),
		sql.Named(
			"interval",
			fmt.Sprintf(
				"%s %s",
				strconv.FormatUint(periodCount, 10),
				string(period),
			),
		),
	).Scan(&result)

	if tx.Error != nil {
		return result, tx.Error
	}

	if rm.db.Error != nil {
		return result, rm.db.Error
	}

	return result, nil
}

func (rm *PipelineReadModel) GetPipelineRunResultQuantile(pipelineUuid uuid.UUID, level float64, period Period, periodCount uint64) (float64, error) {
	var result float64

	query := `SELECT
				ifNotFinite(quantile(@level)(metric_value), 0)
			FROM ` + rm.statsTable + `
			WHERE
				pipeline_uuid = @pipelineUuid
				AND is_metric_value_set = 1
				AND executed_at BETWEEN (NOW() - INTERVAL @interval) AND NOW()
			`

	tx := rm.db.Raw(
		query,
		sql.Named("level", level),
		sql.Named("pipelineUuid", pipelineUuid),
		sql.Named(
			"interval",
			fmt.Sprintf(
				"%s %s",
				strconv.FormatUint(periodCount, 10),
				string(period),
			),
		),
	).Scan(&result)

	if tx.Error != nil {
		return result, tx.Error
	}

	if rm.db.Error != nil {
		return result, rm.db.Error
	}

	return result, nil
}

func (rm *PipelineReadModel) GetPipelineDurationStatsQuantile(pipelineUuid uuid.UUID) (int, error) {
	var result float64

	query := `
			SELECT 
				ifNotFinite(quantile(0.8)(wf_duration), 0)
			FROM (
				SELECT SUM(duration) AS wf_duration 
				FROM task_runs
				WHERE pipeline_run_uuid IN (
					SELECT 
						DISTINCT pipeline_run_uuid 
					FROM task_runs
					WHERE
						pipeline_uuid = @pipelineUuid
					ORDER BY executed_at DESC
					LIMIT @numLastPipelineRuns
				)
				GROUP BY pipeline_run_uuid
			);
			`

	tx := rm.db.Raw(
		query,
		sql.Named("level", 0.8),
		sql.Named("numLastPipelineRuns", 5),
		sql.Named("pipelineUuid", pipelineUuid),
	).Scan(&result)

	if tx.Error != nil {
		return int(result), tx.Error
	}

	if rm.db.Error != nil {
		return int(result), rm.db.Error
	}

	return int(result), nil
}

func (rm *PipelineReadModel) GetMetricValues(pipelineUuid uuid.UUID, dateFrom, dateTo time.Time) ([]dto.MetricValueItem, error) {
	result := make([]dto.MetricValueItem, 0)

	query := `SELECT
				organization_uuid,
				pipeline_uuid,
				metric_value,
				duration,
				executed_at AS ExecutedAt
			FROM ` + rm.statsTable + `
			WHERE
				pipeline_uuid = @pipelineUuid
				AND executed_at BETWEEN @start AND @end
				AND is_metric_value_set = 1
			ORDER BY 
				executed_at ASC
	`
	rm.db.Raw(
		query,
		sql.Named("pipelineUuid", pipelineUuid),
		sql.Named("start", dateFrom),
		sql.Named("end", dateTo),
	).Scan(&result)

	return result, rm.db.Error
}

func (rm *PipelineReadModel) GetLastMetricValues(pipelineUuid uuid.UUID, num int64) ([]dto.MetricValueItem, error) {
	result := make([]dto.MetricValueItem, 0)

	query := `SELECT
				organization_uuid,
				pipeline_uuid,
				metric_value,
				duration,
				executed_at AS ExecutedAt
			FROM ` + rm.statsTable + `
			WHERE
				pipeline_uuid = @pipelineUuid
				AND is_metric_value_set = 1
			ORDER BY 
				executed_at DESC
			LIMIT @numOfRows
	`
	rm.db.Raw(
		query,
		sql.Named("pipelineUuid", pipelineUuid),
		sql.Named("numOfRows", num),
	).Scan(&result)

	return result, rm.db.Error
}

func (rm *PipelineReadModel) init(statsTable string) error {
	var err error
	rm.db, err = db.Instance()
	if err != nil {
		return err
	}
	rm.statsTable = statsTable

	return nil
}

func NewPipelineReadModel() (*PipelineReadModel, error) {
	rm := &PipelineReadModel{}
	err := rm.init(config.Cfg().DB.StatsTable)
	if err != nil {
		return nil, err
	}

	return rm, nil
}
