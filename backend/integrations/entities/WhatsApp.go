package entities

type WhatsApp struct {
	Id             int64              `json:"-"`
	Integration    Integration        `json:"integration"`
	ContentSid     string             `json:"content_sid"`
}

const IntegrationTypeWhatsApp = "whatsapp"
