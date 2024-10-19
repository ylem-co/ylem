package run

const (
	IdListTypeEnabled  = "enabled"
	IdListTypeDisabled = "disabled"
)

type PipelineRunConfig struct {
	TaskIds        IdList `json:"task_ids"`
	TaskTriggerIds IdList `json:"task_trigger_ids"`
}

type IdList struct {
	Type string   `json:"type"`
	Ids  []string `json:"ids"`
}
