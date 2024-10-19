package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"ylem_users/helpers"
	"ylem_users/repositories"

	log "github.com/sirupsen/logrus"
)

func GetOrganizationDataKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	organizationUuid := vars["uuid"]

	db := helpers.DbConn()
	defer db.Close()

	org, ok := repositories.GetOrganizationByUuid(db, organizationUuid)
	if !ok {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	_, err := w.Write(org.DataKey)
	if err != nil {
		log.Error(err)
	}
}
