package dashboard

type Dashboard struct {
	NumActivePipelines          int  `json:"num_active_pipelines"`
	NumActiveMetrics            int  `json:"num_active_metrics"`
	NumNewPipelines             int  `json:"num_new_pipelines"`
	NumNewMetrics               int  `json:"num_new_metrics"`
	NumRecentlyUpdatedPipelines int  `json:"num_recently_updated_pipelines"`
	NumRecentlyUpdatedMetrics   int  `json:"num_recently_updated_metrics"`
	NumMyPipelines              int  `json:"num_my_pipelines"`
	NumMyMetrics                int  `json:"num_my_metrics"`
	NumScheduledPipelines       int  `json:"num_scheduled_pipelines"`
	NumExtTriggeredPipelines    int  `json:"num_externally_triggered_pipelines"`
	NumPipelineTemplates        int  `json:"num_pipeline_templates"`
	NumMetricTemplates          int  `json:"num_metric_templates"`
}

type GroupedItems struct {
	Items []GroupedItem `json:"items"`
}

type GroupedItem struct {
	Count  int    `json:"count"`
	Year   int    `json:"year"`
	Month  string `json:"month"`
	Week   int    `json:"week"`
}

const GroupByMonth = "month"
const GroupByWeek  = "week"
