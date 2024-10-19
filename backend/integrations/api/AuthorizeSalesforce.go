package api

import (
	"github.com/ylem-co/salesforce-client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"ylem_integrations/config"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	"ylem_integrations/services/aws/kms"
)

func AuthorizeSalesforce(w http.ResponseWriter, r *http.Request) {
	log.Trace("Authorizing Salesforce")

	code := r.URL.Query()["code"][0]
	state := r.URL.Query()["state"][0]

	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Finding Salesforce auth by State")
	entity, err := repositories.FindSalesforceAuthorizationByState(db, state)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	log.Tracef("Exchanging tokens")
	token, err := salesforceclient.ExchangeCode(r.Context(), code)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	log.Debugf("Granted access token: %s", token.AccessToken)

	accessTokenSecretBox, refreshTokenSecretBox, err := encryptSalesforceTokens(w, r, entity, token)
	if err != nil {
		return
	}

	ticketsScope := "api chatter_api refresh_token"
	domain := token.Extra("instance_url").(string)
	entity.AccessToken = accessTokenSecretBox
	entity.RefreshToken = refreshTokenSecretBox
	entity.Scopes = &ticketsScope
	entity.IsActive = true
	entity.Domain = &domain

	err = repositories.UpdateSalesforceAuthorization(db, entity)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	url := strings.ReplaceAll(
		config.Cfg().NetworkConfig.SalesforceAfterAuthorizationRedirectUrl,
		"{uuid}",
		entity.Uuid,
	)

	log.Tracef("Salesforce is authorized. Redirecting...")
	http.Redirect(w, r, url, http.StatusFound)
}

func encryptSalesforceTokens(w http.ResponseWriter, r *http.Request, entity *entities.SalesforceAuthorization, token *oauth2.Token) (*kms.SecretBox, *kms.SecretBox, error) {
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

	accessTokenSecretBox := kms.NewSealedSecretBox(encryptedAccessToken)

	encryptedRefreshToken, err := kms.Encrypt([]byte(token.RefreshToken), key)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, nil, err
	}

	refreshTokenSecretBox := kms.NewSealedSecretBox(encryptedRefreshToken)

	return &accessTokenSecretBox, &refreshTokenSecretBox, nil
}
