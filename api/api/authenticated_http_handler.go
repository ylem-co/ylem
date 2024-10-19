package api

import (
	"errors"
	"net/http"
	"ylem_api/model/entity"
	"ylem_api/model/repository"
	"ylem_api/service"
	"ylem_api/service/oauth"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AuthenticatedHttpHandler func(*entity.OauthToken, http.ResponseWriter, *http.Request)

func Authenticate(handler AuthenticatedHttpHandler) http.HandlerFunc {
	return AuthenticateScoped(handler, "")
}

func AuthenticateScoped(handler AuthenticatedHttpHandler, scope string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := getOauthToken(r)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if t == nil {
			log.Debug("Token not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isCallAllowed := service.ValidatePermissions(
			t.OauthClient.UserUuid.String(),
			t.OauthClient.OrganizationUuid.String(),
			service.PermissionActionCreate,
			service.PermissionResourceTypeApiCall,
			"",
		)

		if !isCallAllowed {
			log.Debug("API is disabled for the subscription plan")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if t.IsAccessExpired() {
			log.Debug("Access token expired")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if scope != "" && !oauth.IsScopeGranted(scope, t.GetScope()) {
			log.Debug("Scope not granted")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		handler(t, w, r)
	}
}

func getOauthToken(r *http.Request) (*entity.OauthToken, error) {
	headerVal := r.Header.Get("Authorization")
	parts := strings.Split(headerVal, "Bearer ")

	if len(parts) != 2 {
		return nil, nil
	}

	token := parts[1]

	repo := repository.NewOauthTokenRepository()
	t, err := repo.FindByAccessToken(token)
	if err != nil {
		return nil, err
	}

	return t, nil
}
