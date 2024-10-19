package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"ylem_users/services"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Role struct {
	Role string `json:"role"`
}

func AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := user.UserUuid

	var role Role

	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &role)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}
		return
	}

	errorFields := ValidateRole(role, w)
	if len(errorFields) > 0 {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": strings.Join(errorFields, ",")})
		w.WriteHeader(http.StatusBadRequest)
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}
		return
	}

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

	roles, _ := json.Marshal([]string{role.Role})
	ok = repositories.AssignRole(db, targetUserUuid, roles)

	if ok {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(500)
	}
}

func ValidateRole(role Role, w http.ResponseWriter) []string {
	var errorFields []string

	if role.Role != entities.ROLE_ORGANIZATION_ADMIN && role.Role != entities.ROLE_TEAM_MEMBER {
		errorFields = append(errorFields, "role")
	}

	return errorFields
}
