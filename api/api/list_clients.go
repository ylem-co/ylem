package api

import (
	"encoding/json"
	"net/http"
	"ylem_api/helpers"
	"ylem_api/model/repository"
	"ylem_api/service"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func ListClientsByUser(w http.ResponseWriter, r *http.Request) {
	authData := service.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil || authData.Uuid == "" {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	repo := repository.NewOauthClientRepository()
	clients, err := repo.FindAllByUserUuid(uuid.MustParse(authData.Uuid))

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if len(clients) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(clients)

		_, err := w.Write(jsonResponse)
		if err != nil {
		  log.Error(err)
		}
		return
	}

	canPerformOperation := service.ValidatePermissions(authData.Uuid, clients[0].OrganizationUuid.String(), service.PermissionActionReadList, service.PermissionResourceTypeOauthClient, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(clients)

	_, err = w.Write(jsonResponse)
	if err != nil {
	  log.Error(err)
	}
}
