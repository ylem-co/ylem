package api

import (
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

func PipelineDurationStatsQuantileFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuidParamStr := vars["uuid"]
	errs := make([]error, 0)

	err := validation.Validate(uuidParamStr, is.UUIDv4)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		helpers.HttpReturnErrorBadRequest(w, errs)
		return
	}

	uuidParam, _ := uuid.Parse(uuidParamStr)

	rm, err := readmodel.NewPipelineReadModel()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	value, err := rm.GetPipelineDurationStatsQuantile(uuidParam)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	data := map[string]int{
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
