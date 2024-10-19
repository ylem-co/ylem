package services

import (
	"github.com/ylem-co/hubspot-client"
	"ylem_integrations/config"
)

const HubspotTokenSubtractTime = -10

func init() {
	hs := config.Cfg().Hubspot
	hubspotclient.Initiate(hubspotclient.Config{
		ClientID:     hs.OauthClientId,
		ClientSecret: hs.OauthClientSecret,
		RedirectUrl:  hs.OauthRedirectUri,
		Scopes:       []string{"tickets"},
	})
}
