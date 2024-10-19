package entity

import (
	"time"
	"encoding/json"
	"ylem_statistics/config"
	"ylem_statistics/helpers"

	"github.com/google/uuid"
)

type TaskRun struct {
	Uuid             uuid.UUID `json:"uuid" gorm:"primaryKey"`
	ExecutorUuid     uuid.UUID `json:"executor_uuid"`
	OrganizationUuid uuid.UUID `json:"organization_uuid"`
	CreatorUuid      uuid.UUID `json:"creator_uuid"`
	PipelineUuid     uuid.UUID `json:"pipeline_uuid"`
	PipelineRunUuid  uuid.UUID `json:"pipeline_run_uuid"`
	TaskUuid         uuid.UUID `json:"task_uuid"`
	TaskType         string    `json:"task_type"`
	PipelineType     string    `json:"pipeline_type"`
	IsInitialTask    bool      `json:"is_initial_task"`
	IsFinalTask      bool      `json:"is_final_task"`
	IsSuccessful     bool      `json:"is_successful"`
	IsFatalFailure   bool      `json:"is_fatal_failure"`
	ExecutedAt       time.Time `json:"executed_at"`
	Duration         uint32    `json:"duration"`
	IsMetricValueSet bool      `json:"is_metric_value_set"`
	MetricValue      float64   `json:"metric_value"`
	Output           []byte    `json:"output"`
}

func (tr TaskRun) MarshalJSON() ([]byte, error) {
	type Alias TaskRun
	return json.Marshal(&struct {
		ExecutedAt string `json:"executed_at"`
		Alias
	}{
		ExecutedAt: tr.ExecutedAt.Format(helpers.DateTimeFormat),
		Alias:      Alias(tr),
	})
}

func (tr TaskRun) TableName() string {
	return config.Cfg().DB.StatsTable
}
