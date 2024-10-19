package pipeline

const (
	PipelineIsActive    = 1
	PipelineIsNotActive = 0
)

type Pipeline struct {
	Id               int64  `json:"-"`
	Uuid             string `json:"uuid"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	OrganizationUuid string `json:"organization_uuid"`
	CreatorUuid      string `json:"-"`
	ElementsLayout   string `json:"elements_layout"`
	Preview          []byte `json:"-"`
	FolderUuid       string `json:"folder_uuid"`
	FolderId         int64  `json:"-"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	IsActive         int8   `json:"-"`
	IsPaused         int8   `json:"is_paused"`
	IsTemplate       int8   `json:"is_template"`
	Schedule         string `json:"schedule"`
}

type SearchedPipeline struct {
	Id               int64  `json:"-"`
	Uuid             string `json:"uuid"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	OrganizationUuid string `json:"organization_uuid"`
	CreatorUuid      string `json:"-"`
	FolderUuid       string `json:"folder_uuid"`
	FolderId         int64  `json:"-"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	IsTemplate       int8   `json:"is_template"`
	Schedule         string `json:"schedule"`
}

type PipelineRunsPerMonth struct {
	RunCount         int64  `json:"run_count"`
	YearMonth        string `json:"year_month"`
}
