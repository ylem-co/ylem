package entities

type SlackAuthorization struct {
	Id               int64   `json:"-"`
	Name             string  `json:"name"`
	Uuid             string  `json:"uuid"`
	CreatorUuid      string  `json:"-"`
	OrganizationUuid string  `json:"-"`
	State            string  `json:"-"`
	AccessToken      *string `json:"-"`
	Scopes           *string `json:"-"`
	BotUserId        *string `json:"-"`
	IsActive         bool    `json:"is_active"`
}
