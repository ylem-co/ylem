package entities

import "ylem_integrations/services/aws/kms"

const IntegrationTypeTableau = "tableau"

const (
	TableauModeOverwrite = "overwrite"
	TableauModeAppend    = "append"
)

type Tableau struct {
	Id             int64          `json:"-"`
	Integration    Integration    `json:"integration"`
	Username       *kms.SecretBox `json:"username"`
	Password       *kms.SecretBox `json:"password"`
	Sitename       string         `json:"site_name"`
	ProjectName    string         `json:"project_name"`
	DatasourceName string         `json:"datasource_name"`
	Mode           string         `json:"mode"`
}
