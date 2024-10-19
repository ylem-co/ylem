package api

import (
	"time"
	"encoding/json"
	"net/http"
	"ylem_statistics/domain/readmodel"
	"ylem_statistics/domain/readmodel/dto"
	"ylem_statistics/helpers"
	"ylem_statistics/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type LoadStatsFunc func(uuid uuid.UUID, start time.Time, end time.Time) (dto.Stats, error)

var TaskStatsLoadStatsFunc LoadStatsFunc = func(uuid uuid.UUID, start time.Time, end time.Time) (dto.Stats, error) {
	rm, err := readmodel.NewTaskRunReadModel()
	if err != nil {
		return dto.Stats{}, err
	}

	return rm.GetTaskStats(uuid, start, end)
}

var PipelineStatsLoadStatsFunc LoadStatsFunc = func(uuid uuid.UUID, start time.Time, end time.Time) (dto.Stats, error) {
	rm, err := readmodel.NewPipelineReadModel()
	if err != nil {
		return dto.Stats{}, err
	}

	return rm.GetPipelineStats(uuid, start, end)
}

func NewGetStatsFunc(lsf LoadStatsFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
		if user == nil {
			helpers.HttpReturnErrorUnauthorized(w)

			return
		}

		vars := mux.Vars(r)
		uuidParamStr := vars["uuid"]
		dateFromParamStr := vars["dateFrom"]
		dateToParamStr := vars["dateTo"]

		var err error
		errs := make([]error, 0)

		uuidParam, err := uuid.Parse(uuidParamStr)
		if err != nil {
			errs = append(errs, err)
		}

		dateFromParam, err := time.Parse(helpers.DateTimeFormat, dateFromParamStr)
		if err != nil {
			errs = append(errs, err)
		}

		dateToParam, err := time.Parse(helpers.DateTimeFormat, dateToParamStr)
		if err != nil {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			helpers.HttpReturnErrorBadRequest(w, errs)
			return
		}
		data, err := lsf(uuidParam, dateFromParam, dateToParam)
		if err != nil {
			log.Error(err)
			helpers.HttpReturnErrorInternal(w)
			return
		}

		if data.OrganizationUuid != uuid.Nil {
			canPerformOperation := services.ValidatePermissions(
				user.Uuid,
				data.OrganizationUuid.String(),
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
}
