package entities

type Hubspot struct {
	Id                   int64                `json:"-"`
	Integration          Integration          `json:"integration"`
	HubspotAuthorization HubspotAuthorization `json:"authorization"`
	PipelineStageCode    string               `json:"pipeline_stage_code"`
	OwnerCode            string               `json:"owner_code"`
}

const IntegrationTypeHubspot = "hubspot"
