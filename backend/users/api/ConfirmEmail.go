package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"ylem_users/helpers"
	"ylem_users/repositories"
)

const isEmailConfirmed = 1

func ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	emailToken := vars["key"]

	db := helpers.DbConn()
	defer db.Close()

	user, err := repositories.GetUserByEmailToken(db, emailToken)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if user == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	if user.IsEmailConfirmed == isEmailConfirmed {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err = repositories.UpdateUserIsEmailConfirmed(db, user, isEmailConfirmed)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
