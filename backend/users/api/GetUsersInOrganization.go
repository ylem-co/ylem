package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"ylem_users/services"

	log "github.com/sirupsen/logrus"
)

func GetUsersInOrganization(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := user.UserUuid
	db := helpers.DbConn()
	defer db.Close()

	vars := mux.Vars(r)
	organizationUuid := vars["uuid"]

	permissionCheck := services.HttpPermissionCheck{UserUuid: userUUID, OrganizationUuid: organizationUuid, ResourceUuid: "", ResourceType: entities.RESOURCE_USER, Action: entities.ACTION_READ_LIST}
	ok := services.IsUserActionAllowed(permissionCheck)
	if !ok {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	users, ok := repositories.GetUsersByOrganizationUuid(db, organizationUuid)
	if ok {
		response, _ := json.Marshal(
			map[string][]entities.UserToExpose{"items": users},
		)

		_, err := w.Write(response)
		if err != nil {
			log.Error(err)
		}
	} else {
		w.WriteHeader(500)
	}
}
