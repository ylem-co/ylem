package dto

import (
	"time"
	"encoding/json"
	"ylem_statistics/helpers"

	"github.com/google/uuid"
)

type RunStats struct {
	OrganizationUuid uuid.UUID     `json:"-"`
	IsSuccessful     bool          `json:"is_successful"`
	ExecutedAt       time.Time     `json:"executed_at"`
	Duration         time.Duration `json:"duration"`
}

type Stats struct {
	OrganizationUuid    uuid.UUID     `json:"-"`
	SuccessCount        uint          `json:"num_of_successes"`
	FailureCount        uint          `json:"num_of_failures"`
	AverageDuration     uint          `json:"average_duration"`
	IsLastRunSuccessful bool          `json:"is_last_run_successful"`
	LastRunExecutedAt   time.Time     `json:"last_run_executed_at"`
	LastRunDuration     time.Duration `json:"last_run_duration"`
}

type AggregatedStats struct {
	OrganizationUuid uuid.UUID `json:"-"`
	Start            time.Time `json:"date_from"`
	End              time.Time `json:"date_to"`
	SuccessCount     uint      `json:"num_of_successes"`
	FailureCount     uint      `json:"num_of_failures"`
	AverageDuration  uint      `json:"average_duration"`
}

type MetricValueItem struct {
	OrganizationUuid uuid.UUID     `json:"-"`
	PipelineUuid     uuid.UUID     `json:"pipeline_uuid"`
	MetricValue      float64       `json:"metric_value"`
	Duration         time.Duration `json:"duration"`
	ExecutedAt       time.Time     `json:"executed_at"`
}

type RunStatLogs struct {
	OrganizationUuid uuid.UUID     `json:"-"`
	PipelineRunUuid  uuid.UUID     `json:"pipeline_run_uuid"`
	TaskUuid         uuid.UUID     `json:"task_uuid"`
	TaskType         string        `json:"task_type"`
	IsSuccessful     bool          `json:"is_successful"`
	Duration         time.Duration `json:"duration"`
	Output           []byte        `json:"output"`
	MetricValue      float64       `json:"metric_value"`
	ExecutedAt       time.Time     `json:"executed_at"`
}

type TaskRunLogs struct {
	OrganizationUuid uuid.UUID     `json:"-"`
	PipelineRunUuid  uuid.UUID     `json:"pipeline_run_uuid"`
	TaskUuid         uuid.UUID     `json:"task_uuid"`
	TaskType         string        `json:"task_type"`
	PipelineUuid     uuid.UUID     `json:"pipeline_uuid"`
	IsSuccessful     bool          `json:"is_successful"`
	Duration         time.Duration `json:"duration"`
	Output           []byte        `json:"output"`
	MetricValue      float64       `json:"metric_value"`
	ExecutedAt       time.Time     `json:"executed_at"`
}

func (tr RunStats) MarshalJSON() ([]byte, error) {
	type Alias RunStats
	return json.Marshal(&struct {
		Alias
		ExecutedAt string `json:"executed_at"`
	}{
		ExecutedAt: tr.ExecutedAt.Format(helpers.DateTimeFormat),
		Alias:      Alias(tr),
	})
}

func (ts Stats) MarshalJSON() ([]byte, error) {
	type Alias Stats
	return json.Marshal(&struct {
		Alias
		LastRunExecutedAt string `json:"last_run_executed_at"`
	}{
		LastRunExecutedAt: ts.LastRunExecutedAt.Format(helpers.DateTimeFormat),
		Alias:             Alias(ts),
	})
}

func (ts AggregatedStats) MarshalJSON() ([]byte, error) {
	type Alias AggregatedStats
	return json.Marshal(&struct {
		Start string `json:"date_from"`
		End   string `json:"date_to"`
		Alias
	}{
		Start: ts.Start.Format(helpers.DateTimeFormat),
		End:   ts.End.Format(helpers.DateTimeFormat),
		Alias: Alias(ts),
	})
}

func (mvi MetricValueItem) MarshalJSON() ([]byte, error) {
	type Alias MetricValueItem
	return json.Marshal(&struct {
		Alias
		ExecutedAt string `json:"executed_at"`
	}{
		ExecutedAt: mvi.ExecutedAt.Format(helpers.DateTimeFormat),
		Alias:      Alias(mvi),
	})
}
