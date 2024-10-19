package pipelinetemplate

type SharedPipeline struct {
	Id               int64   `json:"-"`
	PipelineUuid     string  `json:"pipeline_uuid"`
	OrganizationUuid string  `json:"organization_uuid"`
	CreatorUuid      string  `json:"creator_uuid"`
	ShareLink        string  `json:"share_link"`
	IsActive         int8    `json:"-"`
	IsLinkPublished  int8    `json:"is_link_published"`
	CreatedAt        string  `json:"-"`
	UpdatedAt        *string `json:"-"`
}
