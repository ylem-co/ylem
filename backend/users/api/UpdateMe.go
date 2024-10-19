package api

import (
	"strings"
	"database/sql"
	"encoding/json"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"ylem_users/services"

	log "github.com/sirupsen/logrus"
)

func UpdateMe(w http.ResponseWriter, r *http.Request) {
	ctxUser := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := ctxUser.UserUuid

	var user services.HttpUser
	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &user)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	errorFields := ValidateUpdatingHttpUser(user, w)
	if len(errorFields) > 0 {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": strings.Join(errorFields, ",")})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	permissionCheck := services.HttpPermissionCheck{UserUuid: userUUID, OrganizationUuid: "", ResourceUuid: userUUID, ResourceType: entities.RESOURCE_USER, Action: entities.ACTION_UPDATE}
	ok := services.IsUserActionAllowed(permissionCheck)
	if !ok {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	if repositories.DoesUserExist(db, user.Email, userUUID) {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid Email", "fields": "email"})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}
		
		return
	}

	ok = SaveUser(db, user, userUUID)

	if ok {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(500)
	}
}

func SaveUser(db *sql.DB, user services.HttpUser, uuid string) bool {
	updateQuery := `UPDATE users 
        SET first_name = ?, last_name = ?, phone = ?
        WHERE uuid = ?
        `

	updateStatement, err := db.Prepare(updateQuery)
	if err != nil {
		return false
	}
	defer updateStatement.Close()

	_, err = updateStatement.Exec(user.FirstName, user.LastName, user.Phone, uuid)
	return err == nil
}

func ValidateUpdatingHttpUser(user services.HttpUser, w http.ResponseWriter) []string {
	var errorFields []string

	if user.FirstName == "" {
		errorFields = append(errorFields, "first_name")
	}

	if user.LastName == "" {
		errorFields = append(errorFields, "last_name")
	}

	if user.Phone != "" {
		if !services.IsPhoneValid(user.Phone) {
			errorFields = append(errorFields, "phone")
		}
	}

	return errorFields
}
