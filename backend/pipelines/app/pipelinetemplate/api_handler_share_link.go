package pipelinetemplate

import (
	"net/http"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func ListMyShareLinks(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	sls, err := FindAllActiveShareLinksForUser(db, authData.Uuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if len(sls) == 0 {
		helpers.HttpResponse(w, http.StatusOK, sls)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, sls[0].OrganizationUuid, services.PermissionActionRead, services.PermissionResourceTypePipeline, sls[0].PipelineUuid)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	helpers.HttpResponse(w, http.StatusOK, map[string]interface{}{
		"items": sls,
	})
}

func GetShareLink(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	shareLink := mux.Vars(r)["shareLink"]
	sl, err := FindActiveShareLink(db, shareLink)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if sl == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, sl.OrganizationUuid, services.PermissionActionRead, services.PermissionResourceTypePipeline, sl.PipelineUuid)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	helpers.HttpResponse(w, http.StatusOK, sl)
}

func GetShareLinkForTemplate(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	templateUuid := mux.Vars(r)["templateUuid"]
	sl, err := FindActiveShareLinkForPipeline(db, templateUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if sl == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, sl.OrganizationUuid, services.PermissionActionRead, services.PermissionResourceTypePipeline, sl.PipelineUuid)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	helpers.HttpResponse(w, http.StatusOK, sl)
}

func ShareTemplate(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	templateUuid := mux.Vars(r)["templateUuid"]

	db := helpers.DbConn()
	defer db.Close()

	tpl, err := pipeline.GetPipelineByUuid(db, templateUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, tpl.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypePipeline, tpl.Uuid)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	if tpl.IsTemplate == 0 {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	isAlreadyShared, err := IsPipelineShared(db, templateUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if isAlreadyShared {
		helpers.HttpReturnErrorConflict(w)
		return
	}

	stpl, err := CreateSharedPipeline(db, tpl, authData.Uuid, authData.OrganizationUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	helpers.HttpResponse(w, http.StatusCreated, stpl)

}

func UnshareTemplate(w http.ResponseWriter, r *http.Request) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	db := helpers.DbConn()
	defer db.Close()

	templateUuid := mux.Vars(r)["templateUuid"]
	sl, err := FindActiveShareLinkForPipeline(db, templateUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	if sl == nil {
		helpers.HttpReturnErrorNotFound(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, sl.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypePipeline, sl.PipelineUuid)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err = DeactivateSharedPipeline(db, templateUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	helpers.HttpResponse(w, http.StatusOK, nil)
}

func PublishShareLink(w http.ResponseWriter, r *http.Request) {
	shareLinkPublishHandler(w, r, true)
}

func UnpublishShareLink(w http.ResponseWriter, r *http.Request) {
	shareLinkPublishHandler(w, r, false)
}

func shareLinkPublishHandler(w http.ResponseWriter, r *http.Request, isPublished bool) {
	authData := services.InitialAuthorization(r.Header.Get("Authorization"))
	if authData == nil {
		helpers.HttpReturnErrorUnauthorized(w)
		return
	}

	templateUuid := mux.Vars(r)["templateUuid"]

	db := helpers.DbConn()
	defer db.Close()

	tpl, err := pipeline.GetPipelineByUuid(db, templateUuid)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	canPerformOperation := services.ValidatePermissions(authData.Uuid, tpl.OrganizationUuid, services.PermissionActionUpdate, services.PermissionResourceTypePipeline, tpl.Uuid)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	if tpl.IsTemplate == 0 {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	err = SetShareLinkPublished(db, tpl.Uuid, isPublished)
	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)
		return
	}

	helpers.HttpResponse(w, http.StatusOK, nil)
}
