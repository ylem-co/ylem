package hubspot

import (
	"context"
	"time"
	"ylem_taskrunner/helpers"

	hubspotclient "github.com/ylem-co/hubspot-client"
	"golang.org/x/oauth2"
)

func init() {
	// we don't really need any of that, yet we need the config to be initiated
	hubspotclient.Initiate(hubspotclient.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectUrl:  "",
		Scopes:       []string{},
	})
}

type Authentication struct {
	EncryptedDataKey     []byte
	EncryptedAccessToken []byte
}

func CreateTicket(ctx context.Context, request hubspotclient.CreateTicketRequest, auth Authentication) error {
	accessToken, err := helpers.DecryptData(ctx, auth.EncryptedDataKey, auth.EncryptedAccessToken)
	if err != nil {
		return err
	}

	token := oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(1 * time.Hour), // it's always fresh here
	}

	client, err := hubspotclient.CreateInstance(ctx, &token)
	if err != nil {
		return err
	}

	return client.CreateTicket(request)
}
