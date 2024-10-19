package trial

import (
	"encoding/json"
	"net/http"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"

	log "github.com/sirupsen/logrus"
)

type HttpApiNewTrialPipelines struct {
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4"`
	SourceUuid       string `json:"source_uuid" valid:"uuidv4"`
	DestinationUuid  string `json:"destination_uuid" valid:"uuidv4"`
}

func CreateTrialOnes(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	var reqPipeline HttpApiNewTrialPipelines
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqPipeline)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, reqPipeline.OrganizationUuid, services.PermissionActionCreate, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	err := CreateTrialPipelines(
		db, 
		reqPipeline.OrganizationUuid,
		reqPipeline.SourceUuid, 
		reqPipeline.DestinationUuid,
		authData.Uuid,
	);

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
