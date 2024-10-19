package entities

import "ylem_integrations/services/aws/kms"

type Opsgenie struct {
	Id          int64          `json:"-"`
	Integration Integration    `json:"integration"`
	ApiKey      *kms.SecretBox `json:"api_key"`
}

const IntegrationTypeOpsgenie = "opsgenie"
