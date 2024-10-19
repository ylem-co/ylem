package task

import (
	"time"
	"net/http"
	"encoding/json"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/app/pipeline/common"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Create(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	var reqTask HttpApiNewTask
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqTask)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if !IsTypeSupported(reqTask.Type) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if reqTask.Type == TaskTypeNotification && reqTask.Notification != nil && !IsNotificationTypeSupported(reqTask.Notification.Type) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if reqTask.Type == TaskTypeProcessor && reqTask.Processor != nil && !IsProcessorStrategySupported(reqTask.Processor.Strategy) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if reqTask.Type == TaskTypeApiCall && reqTask.ApiCall != nil && !IsApiCallTypeSupported(reqTask.ApiCall.Type) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if reqTask.Type == TaskTypeTransformer && reqTask.Transformer != nil && !IsTransformerTypeSupported(reqTask.Transformer.Type) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if reqTask.Type == TaskTypeTransformer && reqTask.Transformer != nil && !IsTransformerValid(reqTask.Transformer) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if reqTask.Type == TaskTypeQuery && reqTask.Query != nil && !IsQuerySafe(reqTask.Query.SQLQuery) {
		errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "sql_query"})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(errorJson)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	wf, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if wf == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, wf.OrganizationUuid, services.PermissionActionCreate, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	if reqTask.Type == TaskTypeRunPipeline && reqTask.RunPipeline != nil {
		wfToRun, err := pipeline.GetPipelineByUuid(db, reqTask.RunPipeline.PipelineUuid)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		// forbidden to add run_pipeline tasks to non-metric pipelines
		if wf.Type != common.PipelineTypeMetric {
			helpers.HttpReturnErrorForbidden(w)
			return
		}
		// forbidden to run other metric pipelines from metric pipelines
		if wfToRun.Type == common.PipelineTypeMetric {
			helpers.HttpReturnErrorForbidden(w)
			return
		}

		canPerformOperation := services.ValidatePermissions(authData.Uuid, wfToRun.OrganizationUuid, services.PermissionActionRun, services.PermissionResourceTypePipeline, "")
		if !canPerformOperation {
			helpers.HttpReturnErrorForbidden(w)
			return
		}
	}

	newTask := &Task{
		PipelineId:   wf.Id,
		PipelineUuid: wf.Uuid,
		Name:         reqTask.Name,
		Type:         reqTask.Type,
	}

	_, _, err = CreateTaskWithImplementation(db, newTask, reqTask)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(newTask)
	
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

	task, err := GetTaskByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if task == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	err = DeleteTask(db, task.Id, task.Type, task.ImplementationId)

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

	var reqTask HttpApiUpdatedTask
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqTask)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)
	
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if !IsSeveritySupported(reqTask.Severity) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	pipelineUuid := mux.Vars(r)["pipelineUuid"]
	db := helpers.DbConn()
	defer db.Close()

	wf, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if wf == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, wf.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	task, err := GetTaskByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if task == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	if task.Type == TaskTypeNotification && !IsNotificationTypeSupported(reqTask.Notification.Type) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if task.Type == TaskTypeProcessor && !IsProcessorStrategySupported(reqTask.Processor.Strategy) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if task.Type == TaskTypeApiCall && !IsApiCallTypeSupported(reqTask.ApiCall.Type) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if task.Type == TaskTypeTransformer && !IsTransformerTypeSupported(reqTask.Transformer.Type) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if task.Type == TaskTypeTransformer && !IsUpdatedTransformerValid(reqTask.Transformer) {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if task.Type == TaskTypeQuery && !IsQuerySafe(reqTask.Query.SQLQuery) {
		errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "sql_query"})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(errorJson)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if task.Type == TaskTypeRunPipeline && reqTask.RunPipeline != nil {
		if reqTask.RunPipeline.PipelineUuid == pipelineUuid {
			errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "run_pipeline_uuid"})
			w.WriteHeader(http.StatusBadRequest)
			
			_, error := w.Write(errorJson)
			if error != nil {
				log.Error(error)
			}

			return
		}

		wfToRun, err := pipeline.GetPipelineByUuid(db, reqTask.RunPipeline.PipelineUuid)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		canPerformOperation := services.ValidatePermissions(authData.Uuid, wfToRun.OrganizationUuid, services.PermissionActionRun, services.PermissionResourceTypePipeline, "")
		if !canPerformOperation {
			errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "run_pipeline_uuid_run"})
			w.WriteHeader(http.StatusBadRequest)
			
			_, error := w.Write(errorJson)
			if error != nil {
				log.Error(error)
			}

			return
		}
	}

	task.Name = reqTask.Name
	task.Severity = reqTask.Severity
	task.UpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	err = UpdateTask(db, task, reqTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(task)
	
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

	task, err := GetTaskByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if task == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(task)

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

	tasks, err := GetTasksByPipelineId(db, pipeline.Id)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if tasks == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(tasks)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func SearchInOrganization(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	searchString := mux.Vars(r)["searchString"]

	canPerformOperation := services.ValidatePermissions(authData.Uuid, uuid, services.PermissionActionReadList, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	tasks, err := GetTasksByOrganizationUuidAndSearchString(db, uuid, searchString)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if tasks == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(tasks)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}
