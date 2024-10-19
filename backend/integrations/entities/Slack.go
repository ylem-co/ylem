package entities

type Slack struct {
	Id                 int64              `json:"-"`
	Integration        Integration        `json:"integration"`
	SlackAuthorization SlackAuthorization `json:"authorization"`
	SlackChannelId     *string            `json:"-"`
}

const IntegrationTypeSlack = "slack"
