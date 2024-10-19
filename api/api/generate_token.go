package api

import (
	"net/http"
	"ylem_api/service/oauth"

	log "github.com/sirupsen/logrus"
)

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	srv, err := oauth.NewServer()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = srv.HandleTokenRequest(w, r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
