package result

import (
	"io"
	"time"
	"encoding/json"
	"net/http"
	"ylem_pipelines/app/schedule"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/app/pipeline/run"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func InitiatePipelineRun(w http.ResponseWriter, r *http.Request) {
	input, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	initiatePipelineRun(w, r, input, run.PipelineRunConfig{})
}

func InitiatePipelineRunWithConfig(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	body := &struct {
		Input             []byte                `json:"input"`
		PipelineRunConfig run.PipelineRunConfig `json:"config"`
	}{}

	err = json.Unmarshal(bodyBytes, body)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	initiatePipelineRun(w, r, body.Input, body.PipelineRunConfig)
}

func initiatePipelineRun(w http.ResponseWriter, r *http.Request, input []byte, wrConfig run.PipelineRunConfig) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	db := helpers.DbConn()
	defer db.Close()

	pipeline, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipeline == nil || pipeline.IsTemplate != 0 {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	if pipeline.IsPaused != 0 {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	executeAt := time.Now()
	wrUid := uuid.New()
	sr := schedule.ScheduledRun{
		PipelineRunUuid: wrUid,
		PipelineId:      pipeline.Id,
		Input:           input,
		ExecuteAt:       &executeAt,
		Config:          wrConfig,
	}
	err = schedule.AddScheduledRunsTx(tx, []schedule.ScheduledRun{sr})

	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	tasks, err := task.GetInitialTasks(db, uuid.MustParse(pipeline.Uuid), sr.Config)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	tasksToTrigger := make([]*task.Task, 0)
	tasksToTrigger = append(tasksToTrigger, tasks...)

	if len(tasksToTrigger) == 0 {
		_ = tx.Rollback()
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	err = PurgeTaskRunResultsTx(tx, uuid.MustParse(pipeline.Uuid))
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	trrs := make([]TaskRunResult, 0)
	for _, t := range tasks {
		taskUuid, _ := uuid.Parse(t.Uuid)
		trr, err := CreatePendingTaskRunResultTx(
			tx,
			t.Id,
			taskUuid,
			uuid.Nil,
			wrUid,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Error(err)
			helpers.HttpReturnErrorInternal(w)
			return
		}

		trrs = append(trrs, trr)

		if err != nil {
			_ = tx.Rollback()
			log.Error(err)
			helpers.HttpReturnErrorInternal(w)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	helpers.HttpResponse(w, http.StatusCreated, map[string]interface{}{
		"pipeline_run_uuid": sr.PipelineRunUuid,
		"results":           trrs,
	})
}

func GetPipelineRunResults(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	db := helpers.DbConn()
	defer db.Close()

	pipeline, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, pipeline.OrganizationUuid, services.PermissionActionRead, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	results, err := FindTaskRunResults(db, pipeline.Id)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	wrUuid := uuid.Nil
	if len(results) > 0 {
		wrUuid = results[0].PipelineRunUuid
	}

	helpers.HttpResponse(w, http.StatusOK, map[string]interface{}{
		"pipeline_run_uuid": wrUuid,
		"results":           results,
	})
}
