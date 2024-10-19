package entities

import "ylem_integrations/services/aws/kms"

type JiraAuthorization struct {
	Id               int64         `json:"-"`
	Uuid             string        `json:"uuid"`
	CreatorUuid      string        `json:"-"`
	OrganizationUuid string        `json:"-"`
	Name             string        `json:"name"`
	State            string        `json:"-"`
	IsActive         bool          `json:"is_active"`
	AccessToken      *kms.SecretBox `json:"-"`
	Cloudid          *string       `json:"resource_id"`
	Scopes           *string       `json:"-"`
}
