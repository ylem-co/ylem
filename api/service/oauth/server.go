package oauth

import (
	"ylem_api/model/repository"
	"strings"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/google/uuid"
)

func NewServer() (*server.Server, error) {
	cs := NewOauthClientStore()
	ts := NewOauthTokenStore()
	m := manage.NewDefaultManager()
	m.MapClientStorage(cs)
	m.MapTokenStorage(ts)
	m.SetClientTokenCfg(
		&manage.Config{
			AccessTokenExp:    time.Hour * 24,
			RefreshTokenExp:   time.Hour * 72,
			IsGenerateRefresh: true,
		},
	)

	m.SetRefreshTokenCfg(&manage.RefreshingConfig{
		AccessTokenExp:     time.Hour * 24,
		RefreshTokenExp:    time.Hour * 72,
		IsGenerateRefresh:  true,
		IsRemoveAccess:     false,
		IsRemoveRefreshing: true,
	})

	srv := server.NewDefaultServer(m)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetClientScopeHandler(clientScopeHandler)

	return srv, nil
}

func clientScopeHandler(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
	repo := repository.NewOauthClientRepository()
	clientUuid, err := uuid.Parse(tgr.ClientID)
	if err != nil {
		return false, err
	}
	cl, err := repo.FindByUuid(clientUuid)
	if err != nil {
		return false, err
	}

	var allowedScopes map[string]bool
	if cl.AllowedScopes == "" {
		allowedScopes = map[string]bool{
			ScopePipelinesRun: true,
			ScopeStatsRead:    true,
		}
	} else {
		allowedScopes = make(map[string]bool)
		for _, v := range NormalizeScopes(cl.AllowedScopes) {
			allowedScopes[v] = true
		}

	}

	if tgr.Scope == "" {
		allowedScopesArr := make([]string, 0)
		for k := range allowedScopes {
			allowedScopesArr = append(allowedScopesArr, k)
		}
		tgr.Scope = strings.Join(allowedScopesArr, ",")

		return true, nil
	}

	requestedScopes := NormalizeScopes(tgr.Scope)

	for _, rs := range requestedScopes {
		if _, ok := allowedScopes[rs]; !ok {
			return false, nil
		}
	}

	return true, nil
}
