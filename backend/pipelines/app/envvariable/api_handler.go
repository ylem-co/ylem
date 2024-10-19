package envvariable

import (
	"encoding/json"
	"net/http"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HttpApiNewEnvVariable struct {
	Name             string `json:"name" valid:"type(string)"`
	OrganizationUuid string `json:"organization_uuid" valid:"uuidv4"`
	Value            string `json:"value" valid:"type(string)"`
}

type HttpApiUpdatedEnvVariable struct {
	Name  string `json:"name" valid:"type(string)"`
	Value string `json:"value" valid:"type(string)"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	var reqEnvVariable HttpApiNewEnvVariable
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqEnvVariable)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if reqEnvVariable.Value != "" && !IsEnvVariableValValid(reqEnvVariable.Value) {
		errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "env_variable_value"})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(errorJson)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if reqEnvVariable.Name != "" && !IsEnvVariableNameValid(reqEnvVariable.Name) {
		errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "env_variable_name"})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(errorJson)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	envVariable, err := GetEnvVariableByNameAndOrganizationUuid(db, reqEnvVariable.Name, reqEnvVariable.OrganizationUuid)
	if err != nil {
		log.Error(err)
		return
	}

	if envVariable != nil {
		rp, _ := json.Marshal(map[string]string{"error": "Such environment variable already exist", "fields": "env_variable_name_exists"})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, reqEnvVariable.OrganizationUuid, services.PermissionActionCreate, services.PermissionResourceTypeEnvVariable, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	newEnvVariablePipeline := &EnvVariable{
		Name:             reqEnvVariable.Name,
		Value:            reqEnvVariable.Value,
		OrganizationUuid: reqEnvVariable.OrganizationUuid,
	}

	_, err = CreateEnvVariable(db, newEnvVariablePipeline)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(newEnvVariablePipeline)

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

	deletedEnvVariable, err := GetEnvVariableByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if deletedEnvVariable == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, deletedEnvVariable.OrganizationUuid, services.PermissionActionDelete, services.PermissionResourceTypeEnvVariable, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err = DeleteEnvVariable(db, deletedEnvVariable.Id)
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

	var reqEnvVariable HttpApiUpdatedEnvVariable
	w.Header().Set("Content-Type", "application/json")

	decodeReqErr := helpers.DecodeJSONBody(w, r, &reqEnvVariable)
	if decodeReqErr != nil {
		rp, _ := json.Marshal(decodeReqErr.Msg)
		w.WriteHeader(decodeReqErr.Status)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if reqEnvVariable.Value != "" && !IsEnvVariableValValid(reqEnvVariable.Value) {
		errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "env_variable_value"})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(errorJson)
		if error != nil {
			log.Error(error)
		}

		return
	}

	if reqEnvVariable.Name != "" && !IsEnvVariableNameValid(reqEnvVariable.Name) {
		errorJson, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": "env_variable_name"})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(errorJson)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	uuid := mux.Vars(r)["uuid"]

	updatedEnvVariable, err := GetEnvVariableByUuid(db, uuid)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if updatedEnvVariable == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, updatedEnvVariable.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypeEnvVariable, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	envVariable, err := GetEnvVariableByNameAndOrganizationUuid(db, reqEnvVariable.Name, updatedEnvVariable.OrganizationUuid)
	if err != nil {
		log.Error(err)
		return
	}
	
	if envVariable != nil  && envVariable.Uuid != uuid {
		rp, _ := json.Marshal(map[string]string{"error": "Such environment variable already exist", "fields": "env_variable_name_exists"})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	updatedEnvVariable.Name = reqEnvVariable.Name
	updatedEnvVariable.Value = reqEnvVariable.Value

	err = UpdateEnvVariable(db, updatedEnvVariable)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(updatedEnvVariable)

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

	envVariable, err := GetEnvVariableByUuid(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if envVariable == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, envVariable.OrganizationUuid, services.PermissionActionRead, services.PermissionResourceTypeEnvVariable, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(envVariable)

	_, error := w.Write(jsonResponse)
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

	canPerformOperation := services.ValidatePermissions(authData.Uuid, uuid, services.PermissionActionReadList, services.PermissionResourceTypeEnvVariable, "")
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	envVariables, err := GetEnvVariablesByOrganizationUuid(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if envVariables == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(envVariables)

	_, error := w.Write(jsonResponse)
	if error != nil {
		log.Error(error)
	}
}
