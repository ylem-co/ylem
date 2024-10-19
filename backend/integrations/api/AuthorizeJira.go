package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"ylem_integrations/config"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	"ylem_integrations/services/aws/kms"
)

func AuthorizeJira(w http.ResponseWriter, r *http.Request) {
	log.Trace("Authorizing Jira")

	code := r.URL.Query()["code"][0]
	state := r.URL.Query()["state"][0]

	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Finding Jira auth by State")
	entity, err := repositories.FindJiraAuthorizationByState(db, state)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	log.Tracef("Exchanging tokens")
	grant, err := services.JiraGrantAuthorization(code)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	log.Debugf("Granted access token: %s", grant.AccessToken)

	log.Tracef("Grabbing available resources")
	resources, err := services.JiraListAvailableResources(grant.AccessToken)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	log.Infof("Found %d resources for the organization %s", len(resources), entity.OrganizationUuid)
	if len(resources) == 0 {
		rerr := fmt.Sprintf("No available resources found for the organization %s", entity.OrganizationUuid)
		log.Error(rerr)
		helpers.HttpReturnErrorBadRequest(w, rerr, nil)

		return
	}

	accessTokenSecretBox, err := encryptJiraAccessToken(w, r, err, entity, grant)
	if err != nil {
		return
	}

	entity.AccessToken = accessTokenSecretBox
	entity.Scopes = &grant.Scopes
	entity.IsActive = false

	if len(resources) == 1 {
		entity.IsActive = true
		entity.Cloudid = &resources[0].Id
	}

	err = repositories.UpdateJiraAuthorization(db, entity)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	url := strings.ReplaceAll(
		config.Cfg().NetworkConfig.JiraAfterAuthorizationRedirectUrl,
		"{uuid}",
		entity.Uuid,
	)

	log.Tracef("Jira claims are gotten. Redirecting...")
	http.Redirect(w, r, url, http.StatusFound)
}

func encryptJiraAccessToken(w http.ResponseWriter, r *http.Request, err error, entity *entities.JiraAuthorization, grant *services.JiraAuthorizedGrant) (*kms.SecretBox, error) {
	encryptionKey, error := services.FetchOrganizationDataKey(entity.OrganizationUuid)
	if error != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, error
	}

	key, error := kms.DecryptDataKey(r.Context(), config.Cfg().Aws.KmsKeyId, encryptionKey)
	if error != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, error
	}

	encryptedAccessToken, error := kms.Encrypt([]byte(grant.AccessToken), key)
	if error != nil {
		helpers.HttpReturnErrorInternal(w)

		return nil, error
	}

	accessTokenSecretBox := kms.NewSealedSecretBox(encryptedAccessToken)
	return &accessTokenSecretBox, nil
}
