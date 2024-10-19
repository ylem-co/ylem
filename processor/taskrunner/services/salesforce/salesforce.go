package salesforce

import (
	"context"
	"time"
	"ylem_taskrunner/helpers"

	"github.com/ylem-co/salesforce-client"
	"golang.org/x/oauth2"
)

func init() {
	// we don't really need any of that, yet we need the config to be initiated
	salesforceclient.Initiate(salesforceclient.Config{
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

func CreateCase(ctx context.Context, request salesforceclient.CreateCaseRequest, auth Authentication, domain string) error {
	accessToken, err := helpers.DecryptData(ctx, auth.EncryptedDataKey, auth.EncryptedAccessToken)
	if err != nil {
		return err
	}

	token := oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(1 * time.Hour), // it's always fresh here
	}

	client, err := salesforceclient.CreateInstance(ctx, domain, &token)
	if err != nil {
		return err
	}

	return client.CreateCase(request)
}
