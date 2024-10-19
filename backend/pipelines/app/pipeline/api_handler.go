package pipeline

import (
	"time"
	"io"
	"net/http"
	"encoding/json"
	"ylem_pipelines/app/folder"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"
	"ylem_pipelines/app/pipeline/common"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HttpApiNewPipeline struct {
	Name             string `json:"name" valid:"type(string)"`
	Type             string `json:"type" valid:"type(string),in(generic|metric)"`
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4"`
	FolderUuid       string `json:"folder_uuid" valid:"uuidv4, optional"`
	ElementsLayout   string `json:"elements_layout" valid:"type(string), optional"`
	Schedule         string `json:"schedule" valid:"type(string), optional"`
}

type HttpApiUpdatedPipeline struct {
	Name           string `json:"name" valid:"type(string)"`
	FolderUuid     string `json:"folder_uuid" valid:"uuidv4, optional"`
	ElementsLayout string `json:"elements_layout" valid:"type(string), optional"`
	Schedule       string `json:"schedule" valid:"type(string), optional"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	var reqPipeline HttpApiNewPipeline
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqPipeline)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	currentPipelineCount, err := GetCurrentPipelineCount(db, reqPipeline.OrganizationUuid, reqPipeline.Type)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	resourceType, err := services.GetPipelinePermissionResourceType(reqPipeline.Type)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	canPerformOperation, err := services.ValidateBilledPermissions(
		authData.Uuid,
		reqPipeline.OrganizationUuid,
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

	newPipeline := &Pipeline{
		Name:             reqPipeline.Name,
		Type:             reqPipeline.Type,
		OrganizationUuid: reqPipeline.OrganizationUuid,
		CreatorUuid:      authData.Uuid,
		ElementsLayout:   reqPipeline.ElementsLayout,
		Preview:          make([]byte, 128),
		Schedule:         reqPipeline.Schedule,
	}

	if reqPipeline.FolderUuid != "" {
		folder, err := folder.GetFolderByUuid(db, reqPipeline.FolderUuid)

		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		if folder == nil {
			helpers.HttpReturnErrorNotFound(w)
			return
		}

		newPipeline.FolderId = folder.Id
		newPipeline.FolderUuid = folder.Uuid
	}

	_, err = CreatePipeline(db, newPipeline)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	services.UpdatePipelineConnection(reqPipeline.OrganizationUuid, true)

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(newPipeline)

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
	db := helpers.DbConn()
	defer db.Close()

	deletedPipeline, err := GetPipelineByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if deletedPipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, deletedPipeline.OrganizationUuid, services.PermissionActionDelete, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err = DeletePipeline(db, deletedPipeline.Id)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	pipelines, _ := GetPipelinesByOrganizationUuid(db, deletedPipeline.OrganizationUuid)
	if len(pipelines.Items) == 0 {
		services.UpdatePipelineConnection(deletedPipeline.OrganizationUuid, false)
	}

	w.WriteHeader(http.StatusNoContent)
}

func Toggle(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	db := helpers.DbConn()
	defer db.Close()

	toggledPipeline, err := GetPipelineByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if toggledPipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, toggledPipeline.OrganizationUuid, services.PermissionActionDelete, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err = TogglePipeline(db, toggledPipeline)

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

	var reqPipeline HttpApiUpdatedPipeline
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqPipeline)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	uuid := mux.Vars(r)["uuid"]
	db := helpers.DbConn()
	defer db.Close()

	updatedPipeline, err := GetPipelineByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if updatedPipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, updatedPipeline.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	isScheduleChanged := updatedPipeline.Schedule != reqPipeline.Schedule

	updatedPipeline.Name = reqPipeline.Name
	updatedPipeline.ElementsLayout = reqPipeline.ElementsLayout
	updatedPipeline.UpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	updatedPipeline.Schedule = reqPipeline.Schedule

	if reqPipeline.FolderUuid != "" {
		folder, err := folder.GetFolderByUuid(db, reqPipeline.FolderUuid)

		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		if folder == nil {
			helpers.HttpReturnErrorNotFound(w)
			return
		}

		updatedPipeline.FolderId = folder.Id
		updatedPipeline.FolderUuid = folder.Uuid
	} else {
		updatedPipeline.FolderId = 0
		updatedPipeline.FolderUuid = ""
	}

	err = UpdatePipeline(db, updatedPipeline, isScheduleChanged)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(updatedPipeline)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func UpdatePreview(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	db := helpers.DbConn()
	defer db.Close()

	updatedPipeline, err := GetPipelineByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if updatedPipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, updatedPipeline.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	data, err := io.ReadAll(r.Body)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	updatedPipeline.Preview = data

	err = UpdatePipelinePreview(db, updatedPipeline)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Find(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]

	db := helpers.DbConn()
	defer db.Close()

	pipeline, err := GetPipelineByUuid(db, uuid)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(pipeline)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func FindPreview(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]

	db := helpers.DbConn()
	defer db.Close()

	pipeline, err := GetPipelineByUuid(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipeline == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		helpers.HttpReturnErrorBadRequest(w)
		return
	}

	if r.Form.Get("asTemplate") == "1" {
		if pipeline.IsTemplate != 1 {
			helpers.HttpReturnErrorForbidden(w)
			return
		}
	} else {
		canPerformOperation := services.ValidatePermissions(authData.Uuid, pipeline.OrganizationUuid, services.PermissionActionRead, services.PermissionResourceTypePipeline, "")
		if !canPerformOperation {
			helpers.HttpReturnErrorForbidden(w)
			return
		}
	}

	//w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)

	_, error := w.Write(pipeline.Preview)
	if error != nil {
		log.Error(error)
	}
}

func FindAllInOrganization(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]

	canPerformOperation := services.ValidatePermissions(authData.Uuid, uuid, services.PermissionActionReadList, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	pipelines, err := GetPipelinesByOrganizationUuid(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipelines == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(pipelines)

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

	pipelines, err := GetPipelinesByOrganizationUuidAndSearchString(db, uuid, searchString)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if pipelines == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(pipelines)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func FindAllInOrganizationAndFolder(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	folderUuid := mux.Vars(r)["folderUuid"]

	canPerformOperation := services.ValidatePermissions(authData.Uuid, uuid, services.PermissionActionReadList, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	var pipelines *Pipelines
	var err error
	if folderUuid != "" {
		f, err := folder.GetFolderByUuid(db, folderUuid)

		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		if f == nil {
			helpers.HttpReturnErrorNotFound(w)
			return
		}

		pipelines, err = GetPipelinesByFolderIdAndOrganizationUuid(db, uuid, f.Id)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}
	} else {
		pipelines, err = GetPipelinesByFolderIdAndOrganizationUuid(db, uuid, 0)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}
	}

	if pipelines == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(pipelines)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}

func GetRunsPerOrganizationPerMonth(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))

	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	uuid := mux.Vars(r)["uuid"]
	wfType := mux.Vars(r)["type"]

	if !common.IsTypeSupported(wfType) {
   		helpers.HttpReturnErrorBadRequest(w)
   		return
   	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, uuid, services.PermissionActionReadList, services.PermissionResourceTypePipeline, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	var runs *PipelineRunsPerMonths
	var err error
	runs, err = GetRunsByOrganizationUuid(db, uuid, wfType)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if runs == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(runs)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}
