package api

import (
	"encoding/json"
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ExternalAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	url, err := gothic.GetAuthURL(w, r)
	if err != nil {
		log.Error(err)

		rp, _ := json.Marshal(map[string]string{"error": "Failed to generate a redirect link"})
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(rp)
		if err != nil {
			log.Error(err)
		}

		return
	}

	rp, _ := json.Marshal(map[string]string{"url": url})
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(rp)
	if err != nil {
		log.Error(err)
	}
}
