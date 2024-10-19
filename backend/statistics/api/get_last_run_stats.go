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

type LoadLastRunStatsFunc func(uuid uuid.UUID, end time.Time) (dto.RunStats, error)

var TaskLastRunStatsLoadStateFunc LoadLastRunStatsFunc = func(uuid uuid.UUID, end time.Time) (dto.RunStats, error) {
	rm, err := readmodel.NewTaskRunReadModel()
	if err != nil {
		return dto.RunStats{}, err
	}

	return rm.GetLastTaskRun(uuid, end)
}

var PipelineLastRunStatsLoadStateFunc LoadLastRunStatsFunc = func(uuid uuid.UUID, end time.Time) (dto.RunStats, error) {
	rm, err := readmodel.NewPipelineReadModel()
	if err != nil {
		return dto.RunStats{}, err
	}

	return rm.GetLastPipelineRun(uuid, end)
}

func NewGetLastRunStats(lsf LoadLastRunStatsFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
		if user == nil {
			helpers.HttpReturnErrorUnauthorized(w)

			return
		}

		uuidParamStr := mux.Vars(r)["uuid"]
		uuidParam, err := uuid.Parse(uuidParamStr)
		if err != nil {
			helpers.HttpReturnErrorBadRequest(w, []error{err})
			return
		}

		data, err := lsf(uuidParam, time.Now())
		if err != nil {
			log.Error(err)
			helpers.HttpReturnErrorInternal(w)
			return
		}

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
