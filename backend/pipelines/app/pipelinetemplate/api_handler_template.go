package pipelinetemplate

import (
	"net/http"
	"database/sql"
	"encoding/json"
	"ylem_pipelines/app/folder"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/tasktrigger"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/config"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func ListTemplates(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	orgUuid := ""
	if r.Form.Get("onlySystem") == "1" {
		if len(config.Cfg().SystemOrganizationUuid) == 0 {
			helpers.HttpReturnErrorNotFound(w)
			return
		}

		orgUuid = config.Cfg().SystemOrganizationUuid
	}

	templates, err := pipeline.FindAllAvailableTemplates(db, orgUuid, "", true)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if len(templates.Items) == 0 {
		helpers.HttpResponse(w, http.StatusOK, templates)
		return
	}

	helpers.HttpResponse(w, http.StatusOK, templates)
}

func ListOrganizationTemplates(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	creatorUuid := ""
	if r.Form.Get("onlyMy") == "1" {
		creatorUuid = authData.Uuid
	}

	onlyShared := r.Form.Get("onlyShared") == "1"

	db := helpers.DbConn()
	defer db.Close()

	templates, err := pipeline.FindAllAvailableTemplates(db, authData.OrganizationUuid, creatorUuid, onlyShared)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if len(templates.Items) == 0 {
		helpers.HttpResponse(w, http.StatusOK, templates)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, templates.Items[0].OrganizationUuid, services.PermissionActionReadList, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	helpers.HttpResponse(w, http.StatusOK, templates)
}

func SaveAsTemplate(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	decodedBody := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&decodedBody)
	if err != nil {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	uuid, ok := decodedBody["pipeline_uuid"]
	if !ok {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	wf, tasks, tts, err := loadPipelineData(db, uuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if wf == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	if wf.IsTemplate != 0 {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, wf.OrganizationUuid, services.PermissionActionCreate, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	var f *folder.Folder
	if wf.FolderUuid != "" {
		f, err = folder.GetFolderByUuid(db, wf.FolderUuid)
		if err != nil {
			log.Error(err)
			helpers.HttpReturnErrorInternal(w)
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	newWfUuid, err := ClonePipelineTx(tx, authData.OrganizationUuid, authData.Uuid, wf, tasks, tts, 1, f)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		_ = tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	response := map[string]string{
		"uuid": newWfUuid,
	}

	helpers.HttpResponse(w, http.StatusCreated, response)
}

func CreateFromTemplate(w http.ResponseWriter, r *http.Request) {
	var err error
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	templateUuid := mux.Vars(r)["templateUuid"]

	db := helpers.DbConn()
	defer db.Close()

	tplWf, tasks, tts, err := loadPipelineData(db, templateUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if tplWf == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	if tplWf.IsTemplate == 0 {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	currentPipelineCount, err := pipeline.GetCurrentPipelineCount(db, authData.OrganizationUuid, tplWf.Type)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	resourceType, err := services.GetPipelinePermissionResourceType(tplWf.Type)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	canPerformOperation, err := services.ValidateBilledPermissions(
		authData.Uuid,
		authData.OrganizationUuid,
		services.PermissionActionCreate,
		resourceType,
		"",
		currentPipelineCount,
	)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if !canPerformOperation {
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)
		return
	}

	decodedBody := make(map[string]string)
	err = json.NewDecoder(r.Body).Decode(&decodedBody)
	if err != nil {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	folderUuid := decodedBody["folder_uuid"]
	var f *folder.Folder
	if folderUuid != "" {
		f, err = folder.GetFolderByUuid(db, folderUuid)
		if err != nil {
			log.Error(err)
			helpers.HttpReturnErrorInternal(w)
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	newWfUuid, err := ClonePipelineTx(tx, authData.OrganizationUuid, authData.Uuid, tplWf, tasks, tts, 0, f)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		_ = tx.Rollback()
		return
	}

	newWf, err := pipeline.GetPipelineByUuidTx(tx, newWfUuid)
	if err != nil {
		_ = tx.Rollback()
		helpers.HttpReturnErrorInternal(w)
		return
	}

	newWf.IsTemplate = 0
	err = pipeline.UpdatePipelineTx(tx, newWf, false)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	response := map[string]string{
		"uuid": newWfUuid,
	}

	helpers.HttpResponse(w, http.StatusCreated, response)
}

func ClonePipeline(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	db := helpers.DbConn()
	defer db.Close()

	wf, tasks, tts, err := loadPipelineData(db, uuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if wf == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	currentPipelineCount, err := pipeline.GetCurrentPipelineCount(db, wf.OrganizationUuid, wf.Type)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	resourceType, err := services.GetPipelinePermissionResourceType(wf.Type)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	canPerformOperation, err := services.ValidateBilledPermissions(
		authData.Uuid,
		wf.OrganizationUuid,
		services.PermissionActionCreate,
		resourceType,
		"",
		currentPipelineCount,
	)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if !canPerformOperation {
		helpers.HttpReturnErrorForbiddenQuotaExceeded(w)
		return
	}

	var f *folder.Folder
	if wf.FolderUuid != "" {
		f, err = folder.GetFolderByUuid(db, wf.FolderUuid)
		if err != nil {
			log.Error(err)
			helpers.HttpReturnErrorInternal(w)
			return
		}
	}

	wf.Name = "Copy of " + wf.Name

	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	newWfUuid, err := ClonePipelineTx(tx, authData.OrganizationUuid, authData.Uuid, wf, tasks, tts, 0, f)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		_ = tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	response := map[string]string{
		"uuid": newWfUuid,
	}

	helpers.HttpResponse(w, http.StatusCreated, response)
}

func loadPipelineData(db *sql.DB, pipelineUuid string) (*pipeline.Pipeline, *task.Tasks, *tasktrigger.TaskTriggers, error) {
	wf, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		return nil, nil, nil, err
	}

	if wf == nil {
		return nil, nil, nil, nil
	}

	tasks, err := task.GetTasksByPipelineId(db, wf.Id)
	if err != nil {
		return nil, nil, nil, err
	}

	tts, err := tasktrigger.GetTaskTriggersByPipelineId(db, wf.Id)
	if err != nil {
		return nil, nil, nil, err
	}

	return wf, tasks, tts, nil
}
