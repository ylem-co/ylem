package services

import (
	"github.com/ylem-co/salesforce-client"
	"ylem_integrations/config"
)

func init() {
	sf := config.Cfg().Salesforce
	salesforceclient.Initiate(salesforceclient.Config{
		ClientID:     sf.OauthClientId,
		ClientSecret: sf.OauthClientSecret,
		RedirectUrl:  sf.OauthRedirectUri,
		Scopes:       []string{"api", "chatter_api", "refresh_token"},
	})
}
