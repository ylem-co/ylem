package result

import "github.com/google/uuid"

const (
	StatePending  = "pending"
	StateExecuted = "executed"
)

type TaskRunResult struct {
	Id              int64          `json:"id"`
	State           string         `json:"state"`
	TaskId          int64          `json:"-"`
	TaskUuid        uuid.UUID      `json:"task_uuid"`
	TaskRunUuid     uuid.UUID      `json:"task_run_uuid"`
	PipelineRunUuid uuid.UUID      `json:"pipeline_run_uuid"`
	IsSuccessful    bool           `json:"is_successful"`
	Output          []byte         `json:"output"`
	Errors          []TaskRunError `json:"errors"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}

type TaskRunError struct {
	Id              int64  `json:"id"`
	TaskRunResultId int64  `json:"task_run_result_id"`
	Code            uint   `json:"code"`
	Severity        string `json:"severity"`
	Message         string `json:"message"`
}
