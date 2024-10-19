package api

import (
	"encoding/json"
	"net/http"
	"ylem_api/helpers"
	"ylem_api/model/command"
	"ylem_api/model/repository"
	"ylem_api/service"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type CreateOauthClientResponse struct {
	Uuid   string `json:"uuid"`
	Secret string `json:"secret"`
}

func CreateOauthClient(w http.ResponseWriter, r *http.Request) {
	authData := service.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	c := command.CreateOauthClientCommand{}
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

	canPerformOperation := service.ValidatePermissions(authData.Uuid, authData.OrganizationUuid, service.PermissionActionCreate, service.PermissionResourceTypeOauthClient, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	uid, err := command.NewCreateOauthClientHandler().Handle(c)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	repo := repository.NewOauthClientRepository()
	cl, err := repo.FindByUuid(uid)
	if err != nil {
		log.Error(err)
		return
	}

	clData := &CreateOauthClientResponse{
		Uuid:   uid.String(),
		Secret: cl.Secret,
	}
	resp, err := json.Marshal(clData)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
	  log.Error(err)
	}
}
