package api

import (
	"encoding/json"
	"net/http"
	"ylem_api/helpers"
	"ylem_api/model/command"
	"ylem_api/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func DeleteOauthClient(w http.ResponseWriter, r *http.Request) {
	authData := service.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	c := command.DeleteOauthClientCommand{}
	decodeReqErr := helpers.DecodeJSONBody(w, r, &c)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)

		_, err := w.Write(rp)
		if err != nil {
		  log.Error(err)
		}
		return
	}

	var err error
	c.UserUuid, err = uuid.Parse(authData.Uuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	c.OrganizationUuid, err = uuid.Parse(authData.OrganizationUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	canPerformOperation := service.ValidatePermissions(authData.Uuid, authData.OrganizationUuid, service.PermissionActionDelete, service.PermissionResourceTypeOauthClient, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	uid := mux.Vars(r)["uuid"]
	ok, err := command.NewDeleteOauthClientHandler().Handle(uid)
	if !ok && err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
