package api

import (
	"strconv"
	"encoding/json"
	"net/http"
	"ylem_statistics/domain/readmodel"
	"ylem_statistics/helpers"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func PipelineValueQuantileFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuidParamStr := vars["uuid"]
	levelParamStr := vars["level"]
	periodParamStr := vars["period"]
	perdiodCountParamStr := vars["periodCount"]

	errs := make([]error, 0)

	err := validation.Validate(uuidParamStr, is.UUIDv4)
	if err != nil {
		errs = append(errs, err)
	}

	err = validation.Validate(levelParamStr, is.Float)
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

	level, err := strconv.ParseFloat(levelParamStr, 64)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		helpers.HttpReturnErrorBadRequest(w, errs)
		return
	}

	uuidParam, _ := uuid.Parse(uuidParamStr)
	perdiodCountParam, _ := strconv.ParseUint(perdiodCountParamStr, 10, 64)

	rm, err := readmodel.NewPipelineReadModel()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	value, err := rm.GetPipelineRunResultQuantile(uuidParam, level, readmodel.Period(periodParamStr), perdiodCountParam)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	data := map[string]float64{
		"value": value,
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
