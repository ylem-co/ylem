package entities

import "ylem_integrations/services/aws/kms"

type IncidentIo struct {
	Id          int64          `json:"-"`
	Integration Integration    `json:"integration"`
	ApiKey      *kms.SecretBox `json:"api_key"`
	Visibility  string         `json:"visibility"`
}

const IntegrationTypeIncidentIo = "incidentio"
