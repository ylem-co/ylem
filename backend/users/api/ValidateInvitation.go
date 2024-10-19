package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"ylem_users/helpers"
	"ylem_users/repositories"
)

func ValidateInvitation(w http.ResponseWriter, r *http.Request) {

	db := helpers.DbConn()
	defer db.Close()

	vars := mux.Vars(r)
	key := vars["key"]

	w.Header().Set("Content-Type", "application/json")

	ok := repositories.ValidateInvitationByKey(db, key)
	if ok {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(404)
	}
}
