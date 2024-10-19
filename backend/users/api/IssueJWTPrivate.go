package api

import (
	"encoding/json"
	"net/http"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type IssueJWTPrivate struct {
	OrganizationUuid string `json:"organization_uuid"`
	UserUuid         string `json:"user_uuid"`
}

func (auth *AuthMiddleware) IssueJWTPrivate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	reqData := &IssueJWTPrivate{}

	err := helpers.DecodeJSONBody(w, r, &reqData)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	usr, ok := repositories.GetUserByUuid(db, reqData.UserUuid)
	var rp []byte
	if !ok {
		rp, _ = json.Marshal(map[string]string{"error": "Invalid User", "fields": ""})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}
		
		return
	}

	org, ok := repositories.GetOrganizationByUserUuid(db, reqData.UserUuid)
	if !ok || org.Uuid != reqData.OrganizationUuid {
		rp, _ = json.Marshal(map[string]string{"error": "Invalid organization", "fields": ""})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}
		
		return
	}

	tokenErr, tokenString := CreateJWTToken(usr, auth, &org)
	if tokenErr != nil {
		log.Println(tokenErr)
		return
	}

	rp, _ = json.Marshal(map[string]string{"token": tokenString, "uuid": usr.Uuid, "email": usr.Email, "first_name": usr.FirstName, "last_name": usr.LastName, "phone": usr.Phone, "roles": usr.Roles, "is_email_confirmed": strconv.Itoa(usr.IsEmailConfirmed)})
	w.WriteHeader(http.StatusCreated)
	
	_, error := w.Write(rp)
	if error != nil {
		log.Error(error)
	}
		
}
