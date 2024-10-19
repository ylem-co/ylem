package api

import (
	"net/http"
	"ylem_users/config"
)

func IsExternalAuthAvailable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if config.Cfg().Google.ClientId != "" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
