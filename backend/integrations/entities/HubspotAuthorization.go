package entities

import (
	"time"
	"ylem_integrations/services/aws/kms"
)

type HubspotAuthorization struct {
	Id                   int64          `json:"-"`
	Uuid                 string         `json:"uuid"`
	CreatorUuid          string         `json:"-"`
	OrganizationUuid     string         `json:"-"`
	Name                 string         `json:"name"`
	State                string         `json:"-"`
	IsActive             bool           `json:"is_active"`
	AccessToken          *kms.SecretBox `json:"-"`
	AccessTokenExpiresAt time.Time      `json:"-"`
	RefreshToken         *kms.SecretBox `json:"-"`
	Scopes               *string        `json:"-"`
}
