package api

import (
	"strconv"
	"time"
	"encoding/json"
	"net/http"
	"ylem_statistics/domain/readmodel"
	"ylem_statistics/domain/readmodel/dto"
	"ylem_statistics/helpers"
	"ylem_statistics/services"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type LoadAggregatedStatsFunc func(uuid uuid.UUID, start time.Time, period readmodel.Period, periodCount uint) ([]dto.AggregatedStats, error)

var TaskLoadAggregatedStatsFunc LoadAggregatedStatsFunc = func(uuid uuid.UUID, start time.Time, period readmodel.Period, periodCount uint) ([]dto.AggregatedStats, error) {
	rm, err := readmodel.NewTaskRunReadModel()
	if err != nil {
		return []dto.AggregatedStats{}, err
	}

	return rm.GetTaskAggregatedStats(uuid, start, period, periodCount)
}

var PipelineLoadAggregatedStatsFunc LoadAggregatedStatsFunc = func(uuid uuid.UUID, start time.Time, period readmodel.Period, periodCount uint) ([]dto.AggregatedStats, error) {
	rm, err := readmodel.NewPipelineReadModel()
	if err != nil {
		return []dto.AggregatedStats{}, err
	}

	return rm.GetPipelineAggregatedStats(uuid, start, period, periodCount)
}

func NewGetAggregatedStats(lsf LoadAggregatedStatsFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
		if user == nil {
			helpers.HttpReturnErrorUnauthorized(w)

			return
		}

		vars := mux.Vars(r)
		uuidParamStr := vars["uuid"]
		dateFromParamStr := vars["dateFrom"]
		periodParamStr := vars["period"]
		perdiodCountParamStr := vars["periodCount"]

		errs := make([]error, 0)

		err := validation.Validate(uuidParamStr, is.UUIDv4)
		if err != nil {
			errs = append(errs, err)
		}

		err = validation.Validate(periodParamStr, validation.In(readmodel.ValidPeriods...))
		if err != nil {
			errs = append(errs, err)
		}

		err = validation.Validate(perdiodCountParamStr, is.Digit)
		if err != nil {
			errs = append(errs, err)
		}

		dateFromParam, err := time.Parse(helpers.DateTimeFormat, dateFromParamStr)
		if err != nil {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			helpers.HttpReturnErrorBadRequest(w, errs)
			return
		}

		uuidParam, _ := uuid.Parse(uuidParamStr)
		perdiodCountParam, _ := strconv.ParseUint(perdiodCountParamStr, 10, 32)

		data, err := lsf(uuidParam, dateFromParam, readmodel.Period(periodParamStr), uint(perdiodCountParam))
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
}
