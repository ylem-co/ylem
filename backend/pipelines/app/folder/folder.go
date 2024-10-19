package folder

const FolderIsActive = 1
const FolderIsNotActive = 0

type Folder struct {
	Id               int64  `json:"-"`
	Uuid             string `json:"uuid"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	OrganizationUuid string `json:"organization_uuid"`
	ParentUuid       string `json:"parent_uuid"`
	ParentId         int64  `json:"-"`
	IsActive         int8   `json:"-"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
