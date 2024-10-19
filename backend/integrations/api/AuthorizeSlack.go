package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"ylem_integrations/config"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

func AuthorizeSlack(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Authorizing Slack")

	code := r.URL.Query()["code"][0]
	state := r.URL.Query()["state"][0]

	db := helpers.DbConn()
	defer db.Close()

	log.Tracef("Finding Slack auth by State")
	entity, err := repositories.FindSlackAuthorizationByState(db, state)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	log.Tracef("Exchanging tokens")
	grant, err := services.SlackGrantAuthorization(code)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	entity.AccessToken = &grant.AccessToken
	entity.Scopes = &grant.Scopes
	entity.BotUserId = &grant.BotUserId
	entity.IsActive = true

	err = repositories.UpdateSlackAuthorization(db, entity)

	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	url := strings.ReplaceAll(
		config.Cfg().NetworkConfig.SlackAfterAuthorizationRedirectUrl,
		"{uuid}",
		entity.Uuid,
	)

	http.Redirect(w, r, url, http.StatusFound)
}
