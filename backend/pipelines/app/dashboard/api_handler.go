package dashboard

import (
	"encoding/json"
	"net/http"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"
	"ylem_pipelines/app/pipeline/common"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Find(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]

	db := helpers.DbConn()
	defer db.Close()

	canPerformOperation := services.ValidatePermissions(authData.Uuid, uuid, services.PermissionActionRead, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	dashboard, err := GetDashboardByOrganizationUuid(db, uuid, authData.Uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if dashboard == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(dashboard)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func FindNewGroupedItems(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	itemType := mux.Vars(r)["type"]
	groupBy := mux.Vars(r)["groupBy"]

	db := helpers.DbConn()
	defer db.Close()

	canPerformOperation := services.ValidatePermissions(authData.Uuid, uuid, services.PermissionActionRead, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	if itemType != common.PipelineTypeGeneric && itemType != common.PipelineTypeMetric {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	groupedItems, err := GetGroupedItemsByOrganizationUuid(db, uuid, itemType, groupBy, authData.Uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if groupedItems == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(groupedItems)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}
