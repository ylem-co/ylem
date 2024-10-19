package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
)

func ChangeIntegrationStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var entity *entities.Integration
	db := helpers.DbConn()
	defer db.Close()

	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid uuidv4", "fields": "uuid"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		
		_, err := w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		return
	}

	var err error
	entity, err = repositories.FindIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	status := mux.Vars(r)["status"]
	if !entities.IsIntegrationStatusSupported(status) {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid status"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		
		_, err = w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		return
	}

	err = repositories.ChangeIntegrationStatus(db, entity, status)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}
