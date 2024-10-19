package api

import (
	"encoding/json"
	"net/http"
	"ylem_users/helpers"
	"ylem_users/repositories"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HttpOrganizationConnections struct {
	IsDataSourceCreated  *bool `json:"is_data_source_created" valid:"type(*bool),optional"`
	IsDestinationCreated *bool `json:"is_destination_created" valid:"type(*bool),optional"`
	IsPipelineCreated    *bool `json:"is_pipeline_created" valid:"type(*bool),optional"`
}

func UpdateOrganizationConnections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	organizationUuid := vars["uuid"]

	db := helpers.DbConn()
	defer db.Close()

	org, ok := repositories.GetOrganizationByUuid(db, organizationUuid)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var connections HttpOrganizationConnections

	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &connections)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if connections.IsDataSourceCreated == nil && connections.IsDestinationCreated == nil && connections.IsPipelineCreated == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if connections.IsDataSourceCreated != nil {
		if *connections.IsDataSourceCreated {
			org.IsDataSourceCreated = true
		} else {
			org.IsDataSourceCreated = false
		}
	}

	if connections.IsDestinationCreated != nil {
		if *connections.IsDestinationCreated {
			org.IsDestinationCreated = true
		} else {
			org.IsDestinationCreated = false
		}
	}

	if connections.IsPipelineCreated != nil {
		if *connections.IsPipelineCreated {
			org.IsPipelineCreated = true
		} else {
			org.IsPipelineCreated = false
		}
	}

	ok = repositories.UpdateConnections(db, org)

	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
