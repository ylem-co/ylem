package api

import (
	"strings"
	"database/sql"
	"encoding/json"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"ylem_users/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HttpOrganization struct {
	Name string `json:"name"`
}

func UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := user.UserUuid

	var organization HttpOrganization

	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &organization)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	errorFields := ValidateHttpOrganization(organization, w)
	if len(errorFields) > 0 {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": strings.Join(errorFields, ",")})
		w.WriteHeader(http.StatusBadRequest)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	vars := mux.Vars(r)
	organizationUuid := vars["uuid"]

	permissionCheck := services.HttpPermissionCheck{UserUuid: userUUID, OrganizationUuid: organizationUuid, ResourceUuid: organizationUuid, ResourceType: entities.RESOURCE_ORGANIZATION, Action: entities.ACTION_UPDATE}
	ok := services.IsOrganizationActionAllowed(permissionCheck)
	if !ok {
		helpers.HttpReturnErrorForbidden(w)
		return
	}

	if repositories.DoesOrganizationExist(db, organization.Name, organizationUuid) {
		rp, _ := json.Marshal(map[string]string{"error": "Organization already exist", "fields": "organization_exists"})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}
		
		return
	}

	ok = SaveOrganization(db, organization, organizationUuid)

	if ok {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(500)
	}
}

func SaveOrganization(db *sql.DB, organization HttpOrganization, uuid string) bool {
	updateOrgQuery := `UPDATE organizations 
        SET name = ? 
        WHERE uuid = ?
        `

	updateStatement, err := db.Prepare(updateOrgQuery)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer updateStatement.Close()

	_, err = updateStatement.Exec(organization.Name, uuid)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func ValidateHttpOrganization(organization HttpOrganization, w http.ResponseWriter) []string {
	var errorFields []string

	if organization.Name == "" {
		errorFields = append(errorFields, "name")
	}

	return errorFields
}
