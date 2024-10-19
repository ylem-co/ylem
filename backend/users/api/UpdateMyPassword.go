package api

import (
	"strings"
	"database/sql"
	"encoding/json"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/services"

	log "github.com/sirupsen/logrus"
)

type HttpPassword struct {
	Password       string `json:"password"`
	ConirmPassword string `json:"confirm_password"`
}

func UpdateMyPassword(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := user.UserUuid

	var pwd HttpPassword

	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &pwd)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	errorFields := ValidateUpdatingHttpPassword(pwd, w)
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

	ok = UpdateUsersPassword(db, pwd, userUUID)

	if ok {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(500)
	}
}

func UpdateUsersPassword(db *sql.DB, pwd HttpPassword, uuid string) bool {
	updateQuery := `UPDATE users 
        SET password = ?
        WHERE uuid = ?
        `

	updateStatement, err := db.Prepare(updateQuery)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer updateStatement.Close()

	hash, _ := HashPassword(pwd.Password)

	_, err = updateStatement.Exec(hash, uuid)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func ValidateUpdatingHttpPassword(pwd HttpPassword, w http.ResponseWriter) []string {
	var errorFields []string

	if pwd.Password == "" || !services.IsPasswordValid(pwd.Password) {
		errorFields = append(errorFields, "password")
	}

	if pwd.Password != pwd.ConirmPassword {
		errorFields = append(errorFields, "password")
		errorFields = append(errorFields, "confirm_password")
	}

	return errorFields
}
