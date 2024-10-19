package readmodel

import (
	"time"
	"database/sql"
	"strconv"
	"ylem_statistics/config"
	"ylem_statistics/domain/readmodel/dto"
	"ylem_statistics/services/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRunReadModel struct {
	db         *gorm.DB
	statsTable string
}

func (rm *TaskRunReadModel) init(statsTable string) error {
	var err error
	rm.db, err = db.Instance()
	if err != nil {
		return err
	}
	rm.statsTable = statsTable

	return nil
}

func (rm *TaskRunReadModel) GetTaskStats(uuid uuid.UUID, start time.Time, end time.Time) (dto.Stats, error) {
	var result dto.Stats

	query := `SELECT 
				any(organization_uuid) AS OrganizationUuid,
				SUM(is_successful) AS SuccessCount,
				SUM(NOT is_successful) AS FailureCount,
				toUInt32(AVG(duration)) AS AverageDuration
			FROM 
				` + rm.statsTable + `
			WHERE 
				task_uuid = @taskUuid 
				AND executed_at BETWEEN @start AND @end
			`

	rm.db.Raw(
		query,
		sql.Named("taskUuid", uuid),
		sql.Named("start", start),
		sql.Named("end", end),
	).Scan(&result)

	if rm.db.Error != nil {
		return result, rm.db.Error
	}

	lastRun, err := rm.GetLastTaskRun(uuid, end)
	if err != nil {
		return result, err
	}

	result.IsLastRunSuccessful = lastRun.IsSuccessful
	result.LastRunExecutedAt = lastRun.ExecutedAt
	result.LastRunDuration = lastRun.Duration

	return result, nil
}

func (rm *TaskRunReadModel) GetTaskAggregatedStats(uuid uuid.UUID, start time.Time, period Period, periodCount uint) ([]dto.AggregatedStats, error) {
	var result []dto.AggregatedStats

	query := `SELECT 
				any(organization_uuid) AS OrganizationUuid,
				dateTrunc(@period, executed_at) AS Start,
				(dateTrunc(@period, executed_at) + INTERVAL 1 ` + string(period) + ` - INTERVAL 1 SECOND) AS End,
				SUM(is_successful) AS SuccessCount,
				SUM(NOT is_successful) AS FailureCount,
				toUInt32(AVG(duration)) AS AverageDuration
			FROM 
				` + rm.statsTable + `
			WHERE 
				task_uuid = @taskUuid
				AND executed_at BETWEEN toDate(@start) AND (toDate(@start) + INTERVAL ` + strconv.FormatUint(uint64(periodCount), 10) + ` ` + string(period) + ` - INTERVAL 1 SECOND)
			GROUP BY 
				dateTrunc(@period, executed_at)
			`

	rm.db.Raw(
		query,
		sql.Named("taskUuid", uuid),
		sql.Named("start", start),
		sql.Named("period", period),
	).Scan(&result)

	return result, rm.db.Error
}

func (rm *TaskRunReadModel) GetLastTaskRun(uuid uuid.UUID, end time.Time) (dto.RunStats, error) {
	var result dto.RunStats

	query := `SELECT 
				organization_uuid,
				is_successful AS IsSuccessful,
				executed_at AS ExecutedAt,
				duration AS Duration
			FROM 
				` + rm.statsTable + `
			WHERE 
				task_uuid = @taskUuid 
				AND executed_at <= @end
			ORDER BY executed_at DESC
			LIMIT 1
			`

	rm.db.Raw(query,
		sql.Named("taskUuid", uuid),
		sql.Named("end", end),
	).Scan(&result)

	if rm.db.Error != nil {
		return result, rm.db.Error
	}

	return result, nil
}

func (rm *TaskRunReadModel) GetLastRunsLogGroupedByPipeline(pipelineUuid uuid.UUID, dateFrom, dateTo time.Time) ([]dto.RunStatLogs, error) {
	result := make([]dto.RunStatLogs, 0)

	query := `SELECT
				organization_uuid AS OrganizationUuid,
				pipeline_run_uuid AS PipelineRunUuid,
				task_uuid AS TaskUuid,
				task_type AS TaskType,
				is_successful AS IsSuccessful,
				duration AS Duration,
				output AS Output,
				metric_value AS MetricValue,
				executed_at AS ExecutedAt
			FROM ` + rm.statsTable + `
			WHERE
				pipeline_run_uuid IN (
					SELECT
						pipeline_run_uuid
					FROM ` + rm.statsTable + `
					WHERE
						pipeline_uuid = @pipelineUuid
						AND executed_at BETWEEN @start AND @end
					ORDER BY 
						executed_at DESC
				)
			ORDER BY 
				executed_at DESC
	`
	rm.db.Raw(
		query,
		sql.Named("pipelineUuid", pipelineUuid),
		sql.Named("start", dateFrom),
		sql.Named("end", dateTo),
	).Scan(&result)

	return result, rm.db.Error
}

func (rm *TaskRunReadModel) GetSlowTaskRuns(organizationUuid uuid.UUID, dateFrom, dateTo time.Time, threshold int64, taskType string) ([]dto.TaskRunLogs, error) {
	result := make([]dto.TaskRunLogs, 0)

	query := ""

	if (taskType == "all") {
		query = `SELECT
				organization_uuid AS OrganizationUuid,
				pipeline_run_uuid AS PipelineRunUuid,
				task_uuid AS TaskUuid,
				task_type AS TaskType,
				pipeline_uuid AS PipelineUuid,
				is_successful AS IsSuccessful,
				duration AS Duration,
				output AS Output,
				metric_value AS MetricValue,
				executed_at AS ExecutedAt
			FROM ` + rm.statsTable + `
			WHERE
				organization_uuid = @organizationUuid
				AND executed_at BETWEEN @start AND @end
				AND duration >= @threshold
			ORDER BY 
				executed_at DESC
		`
		rm.db.Raw(
			query,
			sql.Named("organizationUuid", organizationUuid),
			sql.Named("start", dateFrom),
			sql.Named("end", dateTo),
			sql.Named("threshold", threshold),
		).Scan(&result)
	} else {
		query = `SELECT
				organization_uuid AS OrganizationUuid,
				pipeline_run_uuid AS PipelineRunUuid,
				task_uuid AS TaskUuid,
				task_type AS TaskType,
				pipeline_uuid AS PipelineUuid,
				is_successful AS IsSuccessful,
				duration AS Duration,
				output AS Output,
				metric_value AS MetricValue,
				executed_at AS ExecutedAt
			FROM ` + rm.statsTable + `
			WHERE
				organization_uuid = @organizationUuid
				AND executed_at BETWEEN @start AND @end
				AND duration >= @threshold
				AND task_type = @taskType
			ORDER BY 
				executed_at DESC
		`
		rm.db.Raw(
			query,
			sql.Named("organizationUuid", organizationUuid),
			sql.Named("start", dateFrom),
			sql.Named("end", dateTo),
			sql.Named("threshold", threshold),
			sql.Named("taskType", taskType),
		).Scan(&result)
	}

	return result, rm.db.Error
}

func NewTaskRunReadModel() (*TaskRunReadModel, error) {
	rm := &TaskRunReadModel{}
	err := rm.init(config.Cfg().DB.StatsTable)
	if err != nil {
		return nil, err
	}

	return rm, nil
}
