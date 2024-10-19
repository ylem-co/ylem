package tasktrigger

import (
	"time"
	"net/http"
	"encoding/json"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/tasktrigger/types"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HttpApiNewTaskTrigger struct {
	TriggerTaskUuid   string `json:"trigger_task_uuid" valid:"uuidv4,optional"`
	TriggeredTaskUuid string `json:"triggered_task_uuid" valid:"uuidv4"`
	TriggerType       string `json:"trigger_type" valid:"type(string)"`
	Schedule          string `json:"schedule" valid:"type(string),optional"`
}

type HttpApiUpdatedTaskTrigger struct {
	TriggerType string `json:"trigger_type" valid:"type(string)"`
	Schedule    string `json:"schedule" valid:"type(string),optional"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	var reqTaskTrigger HttpApiNewTaskTrigger
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqTaskTrigger)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if !IsTriggerTypeSupported(reqTaskTrigger.TriggerType) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if reqTaskTrigger.TriggerType == types.TriggerTypeSchedule && reqTaskTrigger.Schedule != "" && !IsScheduleValid(reqTaskTrigger.Schedule) {
		errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "schedule"})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(errorJson)
		if error != nil {
			log.Error(error)
		}

		return
	}

	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	db := helpers.DbConn()
	defer db.Close()

	pipeline, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, pipeline.OrganizationUuid, services.PermissionActionCreate, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	var triggerTask *task.Task
	if reqTaskTrigger.TriggerTaskUuid != "" {
		triggerTask, err = task.GetTaskByUuid(db, reqTaskTrigger.TriggerTaskUuid)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		if triggerTask == nil {
			helpers.HttpReturnErrorNotFound(w)
			return
		}
	} else {
		if reqTaskTrigger.TriggerType != types.TriggerTypeSchedule {
			helpers.HttpReturnErrorNotFound(w)
			return
		}
	}

	triggeredTask, err := task.GetTaskByUuid(db, reqTaskTrigger.TriggeredTaskUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if triggeredTask == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	var triggerTaskId int64
	if triggerTask != nil {
		triggerTaskId = triggerTask.Id
	}

	newTaskTrigger := &TaskTrigger{
		PipelineId:      pipeline.Id,
		TriggerTaskId:   triggerTaskId,
		TriggeredTaskId: triggeredTask.Id,
		Schedule:        reqTaskTrigger.Schedule,
		TriggerType:     reqTaskTrigger.TriggerType,
	}

	_, _, err = CreateTaskTrigger(db, newTaskTrigger)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(newTaskTrigger)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	db := helpers.DbConn()
	defer db.Close()

	pipeline, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, pipeline.OrganizationUuid, services.PermissionActionDelete, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err = DeleteTaskTrigger(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func Update(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	var reqTaskTrigger HttpApiUpdatedTaskTrigger
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqTaskTrigger)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if !IsTriggerTypeSupported(reqTaskTrigger.TriggerType) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if reqTaskTrigger.TriggerType == types.TriggerTypeSchedule && reqTaskTrigger.Schedule != "" && !IsScheduleValid(reqTaskTrigger.Schedule) {
		errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "schedule"})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(errorJson)
		if error != nil {
			log.Error(error)
		}

		return
	}

	uuid := mux.Vars(r)["uuid"]
	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	db := helpers.DbConn()
	defer db.Close()

	pipeline, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, pipeline.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	taskTrigger, err := GetTaskTriggerByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if taskTrigger == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	taskTrigger.Schedule = reqTaskTrigger.Schedule
	taskTrigger.TriggerType = reqTaskTrigger.TriggerType
	taskTrigger.UpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	err = UpdateTaskTrigger(db, taskTrigger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(taskTrigger)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func Find(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	db := helpers.DbConn()
	defer db.Close()

	pipeline, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
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

	taskTrigger, err := GetTaskTriggerByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if taskTrigger == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(taskTrigger)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func FindAllInPipeline(w http.ResponseWriter, r *http.Request) {
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
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, pipeline.OrganizationUuid, services.PermissionActionReadList, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	taskTriggers, err := GetTaskTriggersByPipelineId(db, pipeline.Id)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if taskTriggers == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(taskTriggers)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}
