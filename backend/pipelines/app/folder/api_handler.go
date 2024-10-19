package folder

import (
	"time"
	"encoding/json"
	"net/http"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HttpApiNewFolder struct {
	Name             string `json:"name" valid:"type(string)"`
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4"`
	Type             string `json:"type" valid:"type(string),in(generic|metric)"`
	ParentUuid       string `json:"parent_uuid" valid:"uuidv4, optional"`
}

type HttpApiUpdatedFolder struct {
	Name       string `json:"name" valid:"type(string)"`
	ParentUuid string `json:"parent_uuid" valid:"uuidv4, optional"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	var reqFolder HttpApiNewFolder
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqFolder)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, reqFolder.OrganizationUuid, services.PermissionActionCreate, services.PermissionResourceTypeFolder, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	newFolder := &Folder{
		Name:             reqFolder.Name,
		Type:             reqFolder.Type,
		OrganizationUuid: reqFolder.OrganizationUuid,
	}

	if reqFolder.ParentUuid != "" {
		parentFolder, err := GetFolderByUuid(db, reqFolder.ParentUuid)

		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		if parentFolder == nil {
			helpers.HttpReturnErrorNotFound(w)
			return
		}

		newFolder.ParentId = parentFolder.Id
		newFolder.ParentUuid = parentFolder.Uuid
	}

	existingFolder, err := GetFoldersByOrgNameTypeAndParentId(db, newFolder.Name, newFolder.ParentId, newFolder.Type, newFolder.OrganizationUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
	if existingFolder != nil {
		rp, _ := json.Marshal(map[string]string{"error": "Such folder already exist", "fields": "folder_name_exists"})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	_, err = CreateFolder(db, newFolder)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(newFolder)

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

	deletedFolder, err := GetFolderByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if deletedFolder == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, deletedFolder.OrganizationUuid, services.PermissionActionDelete, services.PermissionResourceTypeFolder, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err = DeleteFolder(db, deletedFolder.Id)

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

	var reqFolder HttpApiUpdatedFolder
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqFolder)
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

	updatedFolder, err := GetFolderByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if updatedFolder == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, updatedFolder.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypeFolder, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	updatedFolder.Name = reqFolder.Name
	updatedFolder.UpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)

	if reqFolder.ParentUuid != "" {
		parentFolder, err := GetFolderByUuid(db, reqFolder.ParentUuid)

		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		if parentFolder == nil {
			helpers.HttpReturnErrorNotFound(w)
			return
		}

		updatedFolder.ParentId = parentFolder.Id
		updatedFolder.ParentUuid = reqFolder.ParentUuid
	} else {
		updatedFolder.ParentId = 0
		updatedFolder.ParentUuid = ""
	}

	existingFolder, err := GetFoldersByOrgNameTypeAndParentId(db, reqFolder.Name, updatedFolder.ParentId, updatedFolder.Type, updatedFolder.OrganizationUuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
	if existingFolder != nil  && existingFolder.Uuid != uuid {
		rp, _ := json.Marshal(map[string]string{"error": "Such folder already exist", "fields": "folder_name_exists"})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	err = UpdateFolder(db, updatedFolder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(updatedFolder)

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

	db := helpers.DbConn()
	defer db.Close()

	folder, err := GetFolderByUuid(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if folder == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, folder.OrganizationUuid, services.PermissionActionRead, services.PermissionResourceTypeFolder, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(folder)

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

	canPerformOperation := services.ValidatePermissions(authData.Uuid, uuid, services.PermissionActionReadList, services.PermissionResourceTypeFolder, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	var folders *Folders
	var err error
	if folderUuid != "" {
		parentFolder, err := GetFolderByUuid(db, folderUuid)

		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		if parentFolder == nil {
			helpers.HttpReturnErrorNotFound(w)
			return
		}

		folders, err = GetFoldersByOrganizationUuidAndParentId(db, uuid, parentFolder.Id)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}
	} else {
		folders, err = GetFoldersByOrganizationUuidAndParentId(db, uuid, 0)
			if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}
	}

	if folders == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(folders)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}
