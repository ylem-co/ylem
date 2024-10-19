package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"ylem_users/helpers"
	"ylem_users/repositories"
)

func ConfirmEmailInternal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	userUUID := vars["uuid"]

	db := helpers.DbConn()
	defer db.Close()

	user, ok := repositories.GetUserByUuid(db, userUUID)

	if !ok {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if user.IsEmailConfirmed == isEmailConfirmed {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err := repositories.UpdateUserIsEmailConfirmed(db, &user, isEmailConfirmed)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
