package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"ylem_users/services"
)

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := user.UserUuid

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	targetUserUuid := vars["uuid"]

	db := helpers.DbConn()
	defer db.Close()

	org, ok := repositories.GetOrganizationByUserUuid(db, targetUserUuid)
	if !ok {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	permissionCheck := services.HttpPermissionCheck{UserUuid: userUUID, OrganizationUuid: org.Uuid, ResourceUuid: targetUserUuid, ResourceType: entities.RESOURCE_USER, Action: entities.ACTION_CREATE}
	ok = services.IsUserActionAllowed(permissionCheck)
	if !ok {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	ok = repositories.ActivateUser(db, targetUserUuid)

	if ok {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(500)
	}
}
