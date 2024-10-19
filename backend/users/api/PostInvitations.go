package api

import (
	"strings"
	"net/http"
	_ "database/sql"
	"encoding/json"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"ylem_users/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HttpEmailsForInvitations struct {
	Emails string `json:"emails"`
}

func PostInvitations(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := user.UserUuid
	userId := user.UserId

	var emails HttpEmailsForInvitations

	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &emails)
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

	vars := mux.Vars(r)
	organizationUuid := vars["uuid"]

	permissionCheck := services.HttpPermissionCheck{UserUuid: userUUID, OrganizationUuid: organizationUuid, ResourceUuid: "", ResourceType: entities.RESOURCE_INVITATION, Action: entities.ACTION_CREATE}
	ok := services.IsInvitationActionAllowed(permissionCheck)
	if !ok {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	emailsToSave := strings.Split(emails.Emails, ",")
	org, ok := repositories.GetOrganizationByUuid(db, organizationUuid)

	if !ok {
		log.Errorf("Could not get organization by uuid")
		helpers.HttpReturnErrorInternal(w)

		return
	}

	var invitationUuid string
	var invitationCode string
	for _, e := range emailsToSave {
		invitationUuid = uuid.NewString()
		invitationCode = helpers.RandSeq(40)
		_ = repositories.SaveInvitation(db, invitationUuid, e, invitationCode, org.Id, userId)
	}

	w.WriteHeader(201)
}
