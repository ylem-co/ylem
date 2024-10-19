package api

import (
	hubspotclient "github.com/ylem-co/hubspot-client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
	"ylem_integrations/config"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	"ylem_integrations/services/aws/kms"
)

func AuthorizeHubspot(w http.ResponseWriter, r *http.Request) {
	log.Trace("Authorizing Hubspot")

	code := r.URL.Query()["code"][0]
	state := r.URL.Query()["state"][0]

	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Finding Hubspot auth by State")
	entity, err := repositories.FindHubspotAuthorizationByState(db, state)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	log.Tracef("Exchanging tokens")
	token, err := hubspotclient.ExchangeCode(r.Context(), code)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	log.Debugf("Granted access token: %s", token.AccessToken)
	log.Debugf("Granted refresh token: %s", token.RefreshToken)

	accessTokenSecretBox, refreshTokenSecretBox, err := encryptHubspotTokens(w, r, entity, token)
	if err != nil {
		return
	}

	ticketsScope := "tickets"
	entity.AccessToken = accessTokenSecretBox
	entity.RefreshToken = refreshTokenSecretBox
	entity.Scopes = &ticketsScope
	entity.IsActive = true
	// Subtract some time to avoid a chance that in the taskrunner the token expires as we refresh a token in ylem_integrations
	entity.AccessTokenExpiresAt = token.Expiry.Add(services.HubspotTokenSubtractTime * time.Minute)

	err = repositories.UpdateHubspotAuthorization(db, entity)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	url := strings.ReplaceAll(
		config.Cfg().NetworkConfig.HubspotAfterAuthorizationRedirectUrl,
		"{uuid}",
		entity.Uuid,
	)

	log.Tracef("Hubspot is authorized. Redirecting...")
	http.Redirect(w, r, url, http.StatusFound)
}

func encryptHubspotTokens(w http.ResponseWriter, r *http.Request, entity *entities.HubspotAuthorization, token *oauth2.Token) (*kms.SecretBox, *kms.SecretBox, error) {
	encryptionKey, err := services.FetchOrganizationDataKey(entity.OrganizationUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, nil, err
	}

	key, err := kms.DecryptDataKey(r.Context(), config.Cfg().Aws.KmsKeyId, encryptionKey)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, nil, err
	}

	encryptedAccessToken, err := kms.Encrypt([]byte(token.AccessToken), key)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, nil, err
	}

	encryptedRefreshToken, err := kms.Encrypt([]byte(token.RefreshToken), key)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, nil, err
	}

	accessTokenSecretBox := kms.NewSealedSecretBox(encryptedAccessToken)
	encryptedRefreshTokenSecretBox := kms.NewSealedSecretBox(encryptedRefreshToken)

	return &accessTokenSecretBox, &encryptedRefreshTokenSecretBox, nil
}
