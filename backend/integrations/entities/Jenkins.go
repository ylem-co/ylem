package entities

import "ylem_integrations/services/aws/kms"

type Jenkins struct {
	Id          int64          `json:"-"`
	Integration Integration    `json:"integration"`
	BaseUrl     string         `json:"base_url"`
	Token       *kms.SecretBox `json:"token"`
}

const IntegrationTypeJenkins = "jenkins"
