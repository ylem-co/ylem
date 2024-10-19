package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	"ylem_integrations/services/incidentio"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func GetIncidentIoSeverities(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Listing IncidentIo severities")
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("authorization"))
	if user == nil {
		log.Debugf("User not authenticated")
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	entity, err := repositories.FindIncidentIoIntegration(db, uuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		entity.Integration.OrganizationUuid,
		services.PermissionActionRead,
		services.PermissionResourceTypeIntegration,
		"",
	)

	if !canPerformOperation {
		log.Debugf(
			"User %s can't perform the operation %s in %s",
			user.Uuid,
			services.PermissionActionCreate,
			services.PermissionResourceTypeIntegration,
		)
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	_, err = decryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.ApiKey)
	if err != nil {
		return
	}

	ioClient := incidentio.Instance()
	severities, err := ioClient.GetSeverities(string(entity.ApiKey.PlainValue))
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	sort.Slice(severities, func(i, j int) bool {
		return severities[i].Rank < severities[j].Rank
	})

	jsonResponse, _ := json.Marshal(map[string][]incidentio.Severity{
		"items": severities,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
