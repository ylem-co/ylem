package entities

type Organization struct {
	Id                   int    `json:"id"`
	Uuid                 string `json:"uuid"`
	Name                 string `json:"name"`
	IsDataSourceCreated  bool   `json:"is_data_source_created"`
	IsDestinationCreated bool   `json:"is_destination_created"`
	IsPipelineCreated    bool   `json:"is_pipeline_created"`
	DataKey              []byte `json:"-"`
}
