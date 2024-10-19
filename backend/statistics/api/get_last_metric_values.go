package api

import (
	"strconv"
	"encoding/json"
	"net/http"
	"ylem_statistics/domain/readmodel"
	"ylem_statistics/helpers"
	"ylem_statistics/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func GetLastMetricValues(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	vars := mux.Vars(r)
	uuidParamStr := vars["uuid"]
	numParamStr := vars["num"]

	var err error
	errs := make([]error, 0)

	uuidParam, err := uuid.Parse(uuidParamStr)
	if err != nil {
		errs = append(errs, err)
	}

	numParam, err := strconv.ParseInt(numParamStr, 10, 64)
	if err != nil {
	    errs = append(errs, err)
	}

	if len(errs) > 0 {
		helpers.HttpReturnErrorBadRequest(w, errs)
		return
	}

	rm, err := readmodel.NewPipelineReadModel()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	data, err := rm.GetLastMetricValues(uuidParam, numParam)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if len(data) > 0 {
		orgUuid := data[0].OrganizationUuid
		canPerformOperation := services.ValidatePermissions(
			user.Uuid,
			orgUuid.String(),
			services.PermissionActionReadList,
			services.PermissionResourceTypeStat,
			"",
		)
		if !canPerformOperation {
			helpers.HttpReturnErrorForbidden(w)

			return
		}
	}

	response, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		log.Error(err)
	}
}
