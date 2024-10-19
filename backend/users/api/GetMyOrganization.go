package api

import (
	"encoding/json"
	"net/http"
	"ylem_users/helpers"
	"ylem_users/repositories"

	log "github.com/sirupsen/logrus"
)

func GetMyOrganization(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := user.UserUuid

	w.Header().Set("Content-Type", "application/json")

	db := helpers.DbConn()
	defer db.Close()

	org, ok := repositories.GetOrganizationByUserUuid(db, userUUID)

	if ok {
		rp, _ := json.Marshal(org)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(rp)
		if err != nil {
			log.Error(err)
		}
	} else {
		w.WriteHeader(500)
	}
}
