package pipeline

type Pipelines struct {
	Items []Pipeline `json:"items"`
}

type SearchedPipelines struct {
	Items []SearchedPipeline `json:"items"`
}

type PipelineRunsPerMonths struct {
	Items []PipelineRunsPerMonth `json:"items"`
}
