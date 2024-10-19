package schedule

import (
	"time"
	"ylem_pipelines/app/pipeline/run"

	"github.com/google/uuid"
)

type ScheduledRun struct {
	Id              int64
	PipelineRunUuid uuid.UUID
	PipelineId      int64
	Schedule        string
	Input           []byte
	EnvVars         map[string]interface{}
	Config          run.PipelineRunConfig
	ExecuteAt       *time.Time
}
